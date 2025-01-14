package main

import (
	"flag"
	"fmt"
	"fs-perm-backup/packages/backup"
	"fs-perm-backup/packages/restore"
	"log"
)

func main() {
	// Define flags for mode, target directory, and backup file
	mode := flag.String("mode", "", "Specify the operation mode: 'backup' or 'restore'")                           // or -m
	target := flag.String("target", "", "Specify the target directory for permissions backup or restore")          // or -t
	backupFile := flag.String("backup-file", "", "Specify the backup file (default: './permissions_backup.json')") // or -b

	// Short flags for mode, target, and backup file
	flag.StringVar(mode, "m", "", "Short form for 'mode' flag")
	flag.StringVar(target, "t", "", "Short form for 'target' flag")
	flag.StringVar(backupFile, "b", "", "Short form for 'backup-file' flag")

	flag.Parse()

	// Check if flags were used, otherwise, infer from positional arguments
	if *mode == "" && len(flag.Args()) >= 3 {
		*mode = flag.Args()[0]
		*target = flag.Args()[1]
		*backupFile = flag.Args()[2]
	}

	// Validate inputs
	if *mode == "" {
		log.Fatal("Error: 'mode' is required. Specify it as a flag (--mode or -m) or as the first argument.")
	}
	if *target == "" {
		log.Fatal("Error: 'target' is required. Specify it as a flag (--target or -t) or as the second argument.")
	}
	if *backupFile == "" {
		log.Fatal("Error: 'backup-file' is required. Specify it as a flag (--backup-file or -b) or as the third argument.")
	}

	switch *mode {
	case "backup":
		err := backup.BackupPermissions(*target, *backupFile)
		if err != nil {
			log.Fatalf("Error during backup: %v\n", err)
		}
		fmt.Println("Permissions backup successful")
	case "restore":
		err := restore.RestorePermissions(*target, *backupFile)
		if err != nil {
			log.Fatalf("Error during restore: %v\n", err)
		}
		fmt.Println("Permissions restored successfully")
	default:
		log.Fatal("Error: Invalid mode. Use 'backup' or 'restore' as the mode value.")
	}
}
