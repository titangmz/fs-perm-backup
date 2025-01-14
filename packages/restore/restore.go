package restore

import (
	"fmt"
	"fs-perm-backup/packages/utils"
)

// RestorePermissions restores the permissions of files from a backup file to the specified directory.
func RestorePermissions(directory string, inputFile string) error {
	// Resolve and validate input file
	inputPath, err := utils.ResolveAbsolutePath(inputFile, true)
	if err != nil {
		return fmt.Errorf("invalid input file: %v", err)
	}

	// Resolve and validate directory
	dirPath, err := utils.ResolveAbsolutePath(directory)
	if err != nil {
		return fmt.Errorf("invalid directory: %v", err)
	}

	// Read the JSON content from the input file
	data, err := utils.ReadJSONFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading input file '%s': %v", inputPath, err)
	}

	if err := utils.IsDirectory(dirPath); err != nil {
		return fmt.Errorf("directory validation failed: %v", err)
	}

	// Log restoration process
	fmt.Printf("Restoring permissions...\nInput file: %s\nDirectory: %s\n", inputPath, dirPath)

	// Restore permissions from the JSON data
	if err := utils.RestorePermissionsFromJSON(data); err != nil {
		return fmt.Errorf("error restoring permissions from JSON: %v", err)
	}

	return nil
}
