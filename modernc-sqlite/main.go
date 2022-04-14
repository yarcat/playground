package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

var (
	createTableStmt = `
CREATE TABLE test (
	id 		INTEGER PRIMARY KEY AUTOINCREMENT,
	name 	TEXT
)
`
	insertStmt = `
INSERT INTO test (name) VALUES (?)
`
	selectStmt = `
SELECT id, name FROM test
`
)

func main() {
	db := sqlx.MustOpen("sqlite", ":memory:")
	defer db.Close()
	db.MustExec(createTableStmt)
	for _, name := range []string{"foo", "bar", "larch", "spam", "egg"} {
		db.MustExec(insertStmt, name)
	}
	var entries []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	db.Select(&entries, selectStmt)
	fmt.Println(entries)
}
