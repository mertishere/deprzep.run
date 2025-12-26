package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func GetDB() (*sql.DB, error) {
	// if os.Getenv("NODE_ENV") == "development" {
	// 	dbName := os.Getenv("DB_URL")
	// 	db, err := sql.Open("libsql", dbName)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
	// 		os.Exit(1)
	// 	}
	// 	return db, nil
	// }

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", os.Getenv("DB_URL"), os.Getenv("DB_TOKEN")))
	if err != nil {
		return nil, err
	}

	return db, nil
}
