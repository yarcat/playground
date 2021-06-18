package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatalln("open:", err)
	}
	defer db.Close()

	if err = createTimeseriesTable(db, "test_tt"); err != nil {
		log.Fatal("create:", err)
	}

	tt, err := newTimeseriesTable(db, "test_tt")
	if err != nil {
		log.Fatalln("new:", err)
	}

	for _, cmd := range os.Args[1:] {
		switch cmd {
		case "gen":
			now := time.Now()
			transaction(db, func() { generate(tt) })
			log.Println("gen:", time.Since(now))
		case "del":
			now := time.Now()
			r, err := tt.Cleanup(now.Add(-12 * time.Hour))
			log.Println("del:", time.Since(now))
			log.Printf("del: %v, %v", r, err)
		}
	}
}

func transaction(db *sql.DB, f func()) {
	if _, err := db.Exec("BEGIN"); err != nil {
		log.Fatalln("begin:", err)
	}
	f()
	if _, err := db.Exec("COMMIT"); err != nil {
		log.Fatalln("commit:", err)
	}
}

func generate(tt timeseriesTable) {
	now := time.Now()
	for t := now.Add(-24 * time.Hour); now.After(t); t = t.Add(time.Second) {
		if err := tt.Add(t, rand.Float32()); err != nil {
			log.Fatal(err)
		}
	}
}

type timeseriesTable struct {
	name     string
	db       *sql.DB
	add, del *sql.Stmt
}

func newTimeseriesTable(db *sql.DB, name string) (timeseriesTable, error) {
	const addStmt = `INSERT INTO %s (timestamp, value) VALUES (?, ?)`
	add, err := db.Prepare(fmt.Sprintf(addStmt, name))
	if err != nil {
		return timeseriesTable{}, err
	}
	const delStmt = `DELETE FROM %s WHERE timestamp < ?`
	del, err := db.Prepare(fmt.Sprintf(delStmt, name))
	if err != nil {
		return timeseriesTable{}, err
	}
	return timeseriesTable{
		name: name,
		db:   db,
		add:  add,
		del:  del,
	}, nil
}

func createTimeseriesTable(db *sql.DB, name string) error {
	const createStmt = `CREATE TABLE IF NOT EXISTS %v (timestamp DATETIME, value REAL)`
	if _, err := db.Exec(fmt.Sprintf(createStmt, name), name); err != nil {
		return err
	}
	const indexStmt = `CREATE INDEX IF NOT EXISTS %[1]s_timestamp ON %[1]s(timestamp)`
	if _, err := db.Exec(fmt.Sprintf(indexStmt, name), name); err != nil {
		return err
	}
	return nil
}

func (tt timeseriesTable) Add(t time.Time, v float32) error {
	_, err := tt.add.Exec(t.UTC().Unix(), v)
	return err
}

func (tt timeseriesTable) Values(round time.Duration) ([]time.Time, []float32, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (tt timeseriesTable) Cleanup(before time.Time) (int, error) {
	r, _ := tt.del.Exec(before.UTC().Unix())
	i, err := r.RowsAffected()
	return int(i), err
}
