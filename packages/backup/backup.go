package backup

import (
	"fmt"
	"fs-perm-backup/packages/utils"
	"os"
)

// BackupPermissions saves the permissions of files in the specified directory to a backup file.
func BackupPermissions(directory string, outputFile string) error {
	// Resolve and validate the output file path
	// Validate the output file path
	err := utils.ValidateOutputPath(outputFile)
	if err != nil {
		return fmt.Errorf("invalid output file: %v", err)
	}

	// Resolve and validate the directory
	dirPath, err := utils.ResolveAbsolutePath(directory)
	if err != nil {
		return fmt.Errorf("invalid directory: %v", err)
	}

	outputPath, err := utils.ResolveAbsolutePath(outputFile, false)
	if err != nil {
		return fmt.Errorf("invalid output path: %v", err)
	}

	if err := utils.IsDirectory(dirPath); err != nil {
		return fmt.Errorf("directory validation failed: %v", err)
	}

	// Log backup process
	fmt.Printf("Backing up permissions from directory '%s' to file '%s'\n", dirPath, outputPath)

	// Get directory permissions in JSON format
	permissions, err := utils.GetDirectoryPermissions(dirPath)
	if err != nil {
		return fmt.Errorf("error getting directory permissions: %v", err)
	}

	// Write the permissions to the output file
	if err := os.WriteFile(outputPath, []byte(permissions), 0644); err != nil {
		return fmt.Errorf("error writing permissions to file: %v", err)
	}

	// Placeholder logic for backup
	return nil
}
