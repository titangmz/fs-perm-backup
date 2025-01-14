package utils

import (
	"os"
	"testing"
)

func TestLogMessage(t *testing.T) {
	// Since LogMessage only prints to stdout, we just ensure it doesn't panic
	LogMessage("Test log message")
}

func TestResolveAbsolutePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "Current directory",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "Non-existent path",
			path:    "/path/that/does/not/exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAbsolutePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveAbsolutePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("ResolveAbsolutePath() returned empty path")
			}
		})
	}
}

func TestIsRegularFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test-regular-file-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Regular file",
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			path:    "/non/existent/file",
			wantErr: true,
		},
		{
			name:    "Directory instead of file",
			path:    os.TempDir(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsRegularFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsRegularFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid directory",
			path:    tmpDir,
			wantErr: false,
		},
		{
			name:    "Non-existent directory",
			path:    "/non/existent/dir",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileOrDirectory(t *testing.T) {
	// Create temporary file and directory for testing
	tmpFile, err := os.CreateTemp("", "test-file-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tmpDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid file",
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name:    "Valid directory",
			path:    tmpDir,
			wantErr: false,
		},
		{
			name:    "Non-existent path",
			path:    "/non/existent/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileOrDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileOrDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
