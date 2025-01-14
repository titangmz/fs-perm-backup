# fs-perm-backup

A lightweight tool to backup and restore file and directory permissions in Linux.

## Features
- **Backup Permissions**: Save the permissions of all files and directories in a specified path to a portable file.
- **Restore Permissions**: Restore permissions from a backup file to ensure consistency across systems.

## Why Use fs-perm-backup?
When transferring files or setting up backups, permissions can be altered or lost. `fs-perm-backup` helps:
- Safeguard permissions during data migrations.
- Restore permissions after accidental changes.

## Installation
To install `fs-perm-backup`, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/titangmz/fs-perm-backup.git
    cd fs-perm-backup
    ```

                
2. Build the project:

    ```bash
    go build ./cmd/fs-perm-backup.go
    ```

3. Add the binary to your PATH:
    ```bash
    sudo mv fs-perm-backup /usr/local/bin
    ```
    ## Usage

    The fs-perm-backup tool provides two main operations: backup and restore file permissions.

    ### Basic Commands

    ```bash
    # Using positional arguments
    fs-perm-backup backup <target-path> <backup-file>

    # Using short flags
    fs-perm-backup -m backup -t <target-path> -b <backup-file>

    # Using long flags
    fs-perm-backup --mode backup --target <target-path> --backup-file <backup-file>
    ```

    ### Examples

    ```bash
    # Using positional arguments
    fs-perm-backup backup ./myproject perms.json

    # Using short flags
    fs-perm-backup -m restore -t ./myproject -b perms.json

    # Using long flags
    fs-perm-backup --mode backup --target /var/www/html --backup-file site-perms.json
    ```

    ### Notes
    - Backup files are stored in JSON format
    - Requires appropriate permissions to read/write target paths
    - Supports both flag-style and positional argument syntax


## Contributing
Contributions are welcome! Feel free to submit issues, feature requests, or pull requests. Follow these steps to contribute:

1. Fork the repository.
2. Create a new branch:
    ```bash
    git checkout -b feature-name
    ```
3. Commit changes:
    ```bash
    git commit -m 'Add new feature'
    ```
4. Push to your branch:
    ```bash
    git push origin feature-name
    ```
5. Open a pull request.

## License
This project is licensed under the MIT License.