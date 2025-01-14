#!/bin/bash

# Clean up any previous data
rm -r ./tmp
mkdir ./tmp

# Create nested directories and files
mkdir -p ./tmp/dir1/dir2/dir3
touch ./tmp/dir1/file1.txt
touch ./tmp/dir1/dir2/file2.txt
touch ./tmp/dir1/dir2/dir3/file3.txt

# List the structure with permissions before backup
echo "Directory structure before backup (original permissions):"
ls -lR ./tmp

# Remove any existing backup file
rm -f ./permissions.json

# Build the Go project
go build ./cmd/fs-perm-backup.go

# Backup permissions recursively
./fs-perm-backup backup "./tmp" "permissions.json"

# Change permissions to 000 (no access) for testing restore
chmod -R 000 ./tmp

# List the structure after changing permissions to 000
echo "Directory structure after changing permissions to 000:"
ls -lR ./tmp

# Restore the permissions from the backup
./fs-perm-backup restore "./tmp" "permissions.json"

# List the structure after restoring permissions
echo "Directory structure after restore (restored permissions):"
ls -lR ./tmp

# Clean up after test
rm -f ./permissions.json
rm -r ./tmp
