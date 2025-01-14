package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FilePermissions represents the permissions of a file or directory in a human-readable format.
type FilePermissions struct {
	Path        string `json:"path"`        // Absolute path of the file or directory
	Type        string `json:"type"`        // "file" or "directory"
	User        string `json:"user"`        // User permissions (rwx)
	Group       string `json:"group"`       // Group permissions (rwx)
	Other       string `json:"other"`       // Other permissions (rwx)
	IsDirectory bool   `json:"isDirectory"` // Is it a directory?
}

// GetDirectoryPermissions returns a JSON object tracking files and directories along with their permissions.
func GetDirectoryPermissions(absDirPath string) (string, error) {
	// Slice to hold the permissions of all files and directories
	var permissions []FilePermissions

	// Walk the directory tree
	err := filepath.Walk(absDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path '%s': %v", path, err)
		}

		// Get file mode (permissions)
		mode := info.Mode()

		// Convert permissions to rwx format
		user := rwx(uint32(mode.Perm()>>6) & 7)  // User permissions (bits 6-8)
		group := rwx(uint32(mode.Perm()>>3) & 7) // Group permissions (bits 3-5)
		other := rwx(uint32(mode.Perm()) & 7)    // Other permissions (bits 0-2)

		// Append the file/directory information
		permissions = append(permissions, FilePermissions{
			Path:        path,
			Type:        fileType(info),
			User:        user,
			Group:       group,
			Other:       other,
			IsDirectory: info.IsDir(),
		})
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error walking the directory '%s': %v", absDirPath, err)
	}

	// Convert the permissions slice to JSON
	permissionsJSON, err := json.MarshalIndent(permissions, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling permissions to JSON: %v", err)
	}

	return string(permissionsJSON), nil
}

// rwx converts a numeric permission to rwx format.
func rwx(perm uint32) string {
	chars := []rune{'-', '-', '-'}
	if perm&4 != 0 { // Read bit
		chars[0] = 'r'
	}
	if perm&2 != 0 { // Write bit
		chars[1] = 'w'
	}
	if perm&1 != 0 { // Execute bit
		chars[2] = 'x'
	}
	return string(chars)
}

// fileType determines if the file info represents a file or directory.
func fileType(info os.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}

func LogMessage(message string) {
	fmt.Println("[LOG]:", message)
}

// ResolveAbsolutePath resolves the absolute path for a given relative or absolute path.
// If checkExistence is true (default), it checks if the file or directory exists.
// If checkExistence is false, it skips the existence check.
func ResolveAbsolutePath(path string, checkExistence ...bool) (string, error) {
	// Default value for checkExistence is true
	check := true
	if len(checkExistence) > 0 {
		check = checkExistence[0]
	}

	// Expand `~` to the home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("unable to determine home directory")
		}

		// Ensure proper joining without a double slash
		if len(path) > 1 && path[1] == '/' {
			path = filepath.Join(homeDir, path[2:])
		} else {
			path = homeDir
		}
	}

	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err // Error resolving the absolute path
	}

	// Check if the file or directory exists, only if checkExistence is true
	if check {
		_, err = os.Stat(absPath)
		if os.IsNotExist(err) {
			return "", errors.New("file or directory does not exist: " + absPath)
		} else if err != nil {
			return "", err // Other errors, like permissions
		}
	}

	return absPath, nil
}

// IsRegularFile checks if the given path is a regular file.
func IsRegularFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file '%s' does not exist", path)
		}
		return fmt.Errorf("error accessing file '%s': %w", path, err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("'%s' is not a regular file", path)
	}

	return nil
}

// IsDirectory checks if the given path is a valid directory.
func IsDirectory(path string) error {
	dirInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' does not exist", path)
		}
		return fmt.Errorf("error accessing directory '%s': %w", path, err)
	}

	if !dirInfo.IsDir() {
		return fmt.Errorf("'%s' is not a directory", path)
	}

	return nil
}

// ValidateFileOrDirectory ensures the given path exists as a file or directory.
func ValidateFileOrDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path '%s' does not exist", path)
	} else if err != nil {
		return fmt.Errorf("error accessing path '%s': %w", path, err)
	}

	return nil
}

// ValidateOutputPath checks if all parent directories in the path exist
// and ensures the target file doesn't exist.
func ValidateOutputPath(outputPath string) error {
	// Expand `~` to the home directory
	if strings.HasPrefix(outputPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to determine home directory: %w", err)
		}
		outputPath = filepath.Join(homeDir, outputPath[1:])
	}

	// Split the path into directory part and file name
	dir := filepath.Dir(outputPath)

	// Check if all parent directories exist
	if err := IsDirectory(dir); err != nil {
		return fmt.Errorf("invalid parent directory path: %w", err)
	}

	// Check if a file already exists at the given output path
	_, err := os.Stat(outputPath)
	if err == nil {
		return fmt.Errorf("file '%s' already exists", outputPath)
	} else if !os.IsNotExist(err) {
		// If an error occurred other than the file not existing (e.g., permissions), return it
		return fmt.Errorf("error checking if file '%s' exists: %v", outputPath, err)
	}

	return nil
}

// ReadJSONFile takes a file path, resolves it to an absolute path, reads its content, and returns the data.
func ReadJSONFile(path string) ([]byte, error) {
	// Resolve the absolute path of the file
	absPath, err := ResolveAbsolutePath(path, true)
	if err != nil {
		return nil, fmt.Errorf("error resolving path: %v", err)
	}

	// Open the JSON file
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %v", absPath, err)
	}
	defer file.Close()

	// Read the contents of the file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", absPath, err)
	}

	// Optionally: Decode the JSON content (if you need the content as a Go struct)
	var jsonContent interface{}
	if err := json.Unmarshal(data, &jsonContent); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON content from file '%s': %v", absPath, err)
	}

	// Return the raw data (or jsonContent if you need it parsed)
	return data, nil
}

// RestorePermissionsFromJSON restores permissions based on the provided JSON data.
func RestorePermissionsFromJSON(data []byte) error {
	var permissions []FilePermissions
	if err := json.Unmarshal(data, &permissions); err != nil {
		return fmt.Errorf("error unmarshaling JSON data: %v", err)
	}

	// Loop through each file permission and apply the changes
	for _, perm := range permissions {
		// Validate the permission strings
		if err := ValidatePermissionString(perm.User); err != nil {
			return fmt.Errorf("error converting permissions for file '%s': %v", perm.Path, err)
		}
		if err := ValidatePermissionString(perm.Group); err != nil {
			return fmt.Errorf("error converting permissions for file '%s': %v", perm.Path, err)
		}
		if err := ValidatePermissionString(perm.Other); err != nil {
			return fmt.Errorf("error converting permissions for file '%s': %v", perm.Path, err)
		}

		// Convert rwx to octal for file permission
		permValue, err := convertPermissionsToOctal(perm.User, perm.Group, perm.Other)
		if err != nil {
			return fmt.Errorf("error converting permissions for file '%s': %v", perm.Path, err)
		}

		// Check if the file or directory exists
		if _, err := os.Stat(perm.Path); os.IsNotExist(err) {
			return fmt.Errorf("file or directory '%s' does not exist", perm.Path)
		}

		// Apply the permissions for the file or directory
		if perm.IsDirectory {
			// Ensure it's a directory, and then change permissions
			if err := os.Chmod(perm.Path, os.FileMode(permValue)); err != nil {
				return fmt.Errorf("error restoring directory permissions for '%s': %v", perm.Path, err)
			}
		} else {
			// If it's a file, change the permissions
			if err := os.Chmod(perm.Path, os.FileMode(permValue)); err != nil {
				return fmt.Errorf("error restoring file permissions for '%s': %v", perm.Path, err)
			}
		}

	}

	return nil
}

// convertPermissionsToOctal converts a user, group, and other permissions to an octal value.
func convertPermissionsToOctal(user, group, other string) (int, error) {
	permMap := map[rune]int{
		'r': 4,
		'w': 2,
		'x': 1,
	}

	// Convert each permission string to an integer
	var userPerm, groupPerm, otherPerm int

	for _, c := range user {
		userPerm += permMap[c]
	}
	for _, c := range group {
		groupPerm += permMap[c]
	}
	for _, c := range other {
		otherPerm += permMap[c]
	}

	// Combine the permissions into a single octal value
	return (userPerm << 6) | (groupPerm << 3) | otherPerm, nil
}

// ValidatePermissionString validates a Unix-style permission string (rwx format)
func ValidatePermissionString(permission string) error {
	// Permission string must be exactly 3 characters long and contain only r, w, x, or -
	if len(permission) != 3 {
		return fmt.Errorf("permission string must be exactly 3 characters long")
	}

	// Regex to match valid permission patterns: "r", "w", "x", or "-" for each position
	validPerm := regexp.MustCompile(`^[rwx-]{3}$`)
	if !validPerm.MatchString(permission) {
		return fmt.Errorf("invalid permission string: %s", permission)
	}

	return nil
}
