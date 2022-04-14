package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db := mustOpenDB("sqlite", ":memory:")
	defer db.Close()
	mustCreateTables(db)
	mustInsertNames(db, "foo", "bar", "larch", "spam", "egg")
	mustPrintAllNames(db)

	fmt.Println("Hello World")
}

func mustOpenDB(driver, source string) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("unable to open sqlite connection: %v", err)
	}
	return db
}

const (
	createTableStmt = `CREATE TABLE test (
		id 		INTEGER PRIMARY KEY AUTOINCREMENT,
		name 	TEXT
	)`
	insertStmt = `INSERT INTO test (name) VALUES (?)`
	selectStmt = `SELECT id, name FROM test`
)

func mustCreateTables(db *sql.DB) {
	if _, err := db.Exec(createTableStmt); err != nil {
		log.Fatalf("unable to create table: %v", err)
	}
}

func mustInsertNames(db *sql.DB, names ...string) {
	px, err := db.Prepare(insertStmt)
	if err != nil {
		log.Fatalf("unable to prepare insert: %v", err)
	}
	defer px.Close()

	for _, name := range names {
		if _, err := px.Exec(name); err != nil {
			log.Fatalf("unable to exec prepared insert: %v", err)
		}
	}
}

func mustPrintAllNames(db *sql.DB) {
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
		fmt.Printf("Scanned row with id=%v name=%q\n", id, name)
	}
}
