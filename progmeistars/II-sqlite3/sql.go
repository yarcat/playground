package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Grouping struct {
	db *sql.DB
}

func NewGrouping(albums []Album) (g Grouping, err error) {
	g.db, err = sql.Open("sqlite3", "file::memory:?cache=shared")
	if err == nil {
		err = createAlbumsTable(g.db)
	}
	if err == nil {
		err = insertAlbums(g.db, albums)
	}
	return
}

func (g Grouping) Years() []int {
	years, err := distinctYears(g.db)
	if err != nil {
		log.Printf("years: %v", err)
	}
	return years
}

func (g Grouping) AlbumsForYear(year int) []Album {
	albums, err := albumsForYear(g.db, year)
	if err != nil {
		log.Printf("year albums: %v", err)
	}
	return albums
}

func (g Grouping) Genres() []string {
	genres, err := distinctGenres(g.db)
	if err != nil {
		log.Printf("genres: %v", err)
	}
	return genres
}

func (g Grouping) AlbumsForGenre(genre string) []Album {
	albums, err := albumsForGenre(g.db, genre)
	if err != nil {
		log.Println("genre albums: %v", err)
	}
	return albums
}

func createAlbumsTable(db *sql.DB) error {
	_, err := db.Exec(createAlbumsTableStmt)
	return err
}

func insertAlbums(db *sql.DB, albums []Album) error {
	stmt, err := db.Prepare(insertAlbumsStmt)
	if err != nil {
		return err
	}
	for _, a := range albums {
		_, err = stmt.Exec(a.Artist, a.Genre, a.Title, a.Year, a.Tracks)
		if err != nil {
			return err
		}
	}
	return nil
}

func distinctYears(db *sql.DB) (years []int, err error) {
	var rows *sql.Rows
	if rows, err = db.Query(distinctYearsStmt); err != nil {
		return
	}
	for rows.Next() {
		var y int
		if err = rows.Scan(&y); err == nil {
			years = append(years, y)
		}
	}
	return
}

func albumsForYear(db *sql.DB, year int) (albums []Album, err error) {
	var rows *sql.Rows
	rows, err = db.Query(albumsForYearStmt, year)
	if err == nil {
		albums, err = albumsFromRows(rows)
	}
	return
}

func distinctGenres(db *sql.DB) (genres []string, err error) {
	var rows *sql.Rows
	if rows, err = db.Query(distinctGenresStmt); err != nil {
		return
	}
	for rows.Next() {
		var g string
		if err = rows.Scan(&g); err != nil {
			return
		}
		genres = append(genres, g)
	}
	return
}

func albumsForGenre(db *sql.DB, genre string) (albums []Album, err error) {
	var rows *sql.Rows
	rows, err = db.Query(albumsForGenreStmt, genre)
	if err == nil {
		albums, err = albumsFromRows(rows)
	}
	return
}

func albumsFromRows(rows *sql.Rows) (albums []Album, err error) {
	for rows.Next() {
		var a Album
		err = rows.Scan(&a.Artist, &a.Genre, &a.Title, &a.Year, &a.Tracks)
		if err != nil {
			return
		}
		albums = append(albums, a)
	}
	return
}

const (
	createAlbumsTableStmt = `
CREATE TABLE albums (
	artist TEXT,
	genre  TEXT,
	title  TEXT,
	year   INT,
	tracks INT,
	PRIMARY KEY(artist, title, year)
)`

	insertAlbumsStmt = `
INSERT INTO albums (artist, genre, title, year, tracks)
VALUES (?, ?, ?, ?, ?)`

	distinctYearsStmt  = `SELECT DISTINCT year FROM albums ORDER BY year DESC`
	distinctGenresStmt = `SELECT DISTINCT genre FROM albums ORDER by genre`
	albumsForYearStmt  = `SELECT * FROM albums WHERE year IN (?)`
	albumsForGenreStmt = `SELECT * FROM albums WHERE genre IN (?)`
)
