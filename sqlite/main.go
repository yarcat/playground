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

	const sqlFile = "file:foo.db?cache=shared&timeout=5000"
	db, err := sql.Open("sqlite3", sqlFile)
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
			go func() {
				db, err := sql.Open("sqlite3", sqlFile+"&mode=ro")
				if err != nil {
					log.Fatalf("CNT: OPEN ERR = %v", err)
				}
				for {
					now := time.Now()
					if cnt, err := count(db, "test_tt"); err != nil {
						log.Printf("CNT: %s, ERR = %v", time.Since(now), err)
					} else {
						log.Printf("CNT: %s, %v", time.Since(now), cnt)
					}
					time.Sleep(10 * time.Second)
				}
			}()
			for {
				now := time.Now()
				if err := generate(tt); err != nil {
					log.Printf("GEN: %s, ERR = %v", time.Since(now), err)
				} else {
					log.Printf("GEN: %s", time.Since(now))
				}
				time.Sleep(100 * time.Millisecond)
			}
		case "del":
			now := time.Now()
			r, err := tt.Cleanup(now.Add(-12 * time.Hour))
			log.Println("del:", time.Since(now))
			log.Printf("del: %v, %v", r, err)
		case "scan":
			now := time.Now()
			ts, vs, err := tt.Values(0)
			log.Println("scan:", time.Since(now))
			log.Println("scan:", len(ts), len(vs), err)
		}
	}
}

const countStmt stmt = "SELECT COUNT(*) FROM %s"

func count(db *sql.DB, table string) (int64, error) {
	r := db.QueryRow(countStmt.For(table))
	var s int64
	err := r.Err()
	if err == nil {
		err = r.Scan(&s)
	}
	return s, err
}

func generate(tt timeseriesTable) error {
	now := time.Now()
	var data []TimeValue
	for t := 1; t <= 30; t++ {
		data = append(data, TimeValue{t, now, rand.Float32()})
	}
	return tt.AddBatch(data)
}

// func generate(tt timeseriesTable) {
// 	now := time.Now()
// 	var data []TimeValue
// 	for t := now.Add(-24 * time.Hour); now.After(t); t = t.Add(time.Second) {
// 		data = append(data, TimeValue{t, rand.Float32()})
// 	}
// 	if err := tt.AddBatch(data); err != nil {
// 		log.Fatal(err)
// 	}
// }

type TimeValue struct {
	Typ int
	T   time.Time
	V   float32
}

type timeseriesTable struct {
	name string
	db   *sql.DB
	del  *sql.Stmt
}

func newTimeseriesTable(db *sql.DB, name string) (timeseriesTable, error) {
	const delStmt = `DELETE FROM %s WHERE timestamp < ?`
	del, err := db.Prepare(fmt.Sprintf(delStmt, name))
	if err != nil {
		return timeseriesTable{}, err
	}
	return timeseriesTable{
		name: name,
		db:   db,
		del:  del,
	}, nil
}

func createTimeseriesTable(db *sql.DB, name string) error {
	const createStmt = `CREATE TABLE IF NOT EXISTS %v (row_type INTEGER, timestamp DATETIME, value REAL)`
	if _, err := db.Exec(fmt.Sprintf(createStmt, name), name); err != nil {
		return err
	}
	const indexStmt = `CREATE INDEX IF NOT EXISTS %[1]s_timestamp ON %[1]s(row_type, timestamp)`
	if _, err := db.Exec(fmt.Sprintf(indexStmt, name), name); err != nil {
		return err
	}
	return nil
}

type stmt string

func (s stmt) For(tableName string) string { return fmt.Sprintf(string(s), tableName) }

const insertStmt stmt = "INSERT INTO `%s` (row_type, timestamp, value) VALUES (?, ?, ?)"

func (tt timeseriesTable) AddBatch(tv []TimeValue) error {
	tx, err := tt.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(insertStmt.For(tt.name))
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := range tv {
		if _, err := stmt.Exec(tv[i].Typ, tv[i].T, tv[i].V); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (tt timeseriesTable) Values(round time.Duration) ([]time.Time, []float32, error) {
	const stmt = `
select
	cast(30*round(timestamp/30) as integer) as timestamp_bucket,
	avg(value) as avg_value
from test_tt
group by timestamp_bucket
order by timestamp_bucket`
	rows, err := tt.db.Query(stmt)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var (
		ts   []time.Time
		vals []float32
	)
	for rows.Next() {
		var (
			t int64
			v float32
		)
		if err := rows.Scan(&t, &v); err != nil {
			return nil, nil, err
		}
		ts = append(ts, time.Unix(t, 0))
		vals = append(vals, v)
	}
	return ts, vals, nil
}

func (tt timeseriesTable) Cleanup(before time.Time) (int, error) {
	r, _ := tt.del.Exec(before.UTC().Unix())
	i, err := r.RowsAffected()
	return int(i), err
}
