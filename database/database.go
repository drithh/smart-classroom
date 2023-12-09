// database.go
package database

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func ConnectDB() {
	connString := "user=postgres password=postgres host=localhost port=5432 dbname=classroom sslmode=disable"

	var err error
	db, err = sqlx.Connect("pgx", connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Check if the 'device_settings' table exists, if not, run migration.sql
	if !tableExists("device_settings") {
		fmt.Println("Devices table does not exist. Running migration...")
		if err := runMigration(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running migration: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Connected to the database")
}

func tableExists(tableName string) bool {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", tableName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking if table exists: %v\n", err)
		os.Exit(1)
	}
	return exists
}

func runMigration() error {
	migrationContent, err := os.ReadFile("./database/migration.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(migrationContent))
	return err

	// Placeholder for actual migration logic
	// return nil
}

func CloseDB() {
	db.Close()
	fmt.Println("Connection to the database closed")
}

func GetDB() *sqlx.DB {
	return db
}
