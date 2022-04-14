package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("unable to open sqlite connection: %v", err)
	}

	const createTableStmt = `CREATE TABLE test (
		id 		INTEGER PRIMARY KEY AUTOINCREMENT,
		name 	TEXT
	)`
	if _, err := db.Exec(createTableStmt); err != nil {
		log.Fatalf("unable to create table: %v", err)
	}

	const insertStmt = `INSERT INTO test (name) VALUES ("Foo"), ("Bar")`
	if _, err := db.Exec(insertStmt); err != nil {
		log.Fatalf("unable to insert: %v", err)
	}

	const selectStmt = `SELECT id, name FROM test`
	rows, err := db.Query(selectStmt)
	if err != nil {
		log.Fatalf("unable to execute select: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("unable to scan: %v", err)
		}
		fmt.Printf("Scanned rows with id=%v name=%q\n", id, name)
	}

	fmt.Println("Hello World")
}
