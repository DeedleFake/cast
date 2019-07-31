package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/DeedleFake/cast/db"
)

func main() {
	dbpath := flag.String("db", "", "database location (defaults to a directory in the current user's home directory)")
	flag.Parse()

	if *dbpath == "" {
		u, err := user.Current()
		if err != nil {
			log.Fatalf("Failed to get current user: %v", err)
		}

		*dbpath = filepath.Join(u.HomeDir, ".cast", "cast.db")
	}

	err := os.MkdirAll(filepath.Dir(*dbpath), 0644)
	if err != nil {
		log.Fatalf("Failed to create directory for DB: %v", err)
	}

	db, err := db.Open(*dbpath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
}
