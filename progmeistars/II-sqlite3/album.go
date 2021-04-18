package main

import (
	"bufio"
	"io"
	"strings"
)

type Album struct {
	Artist, Genre, Title string
	Year, Tracks         int
}

func ParseAlbum(s string) (a Album) {
	p := Parser{Input: s}
	for p.Next() {
		switch p.Tag() {
		case "artist":
			a.Artist = p.Text()
		case "title":
			a.Title = p.Text()
		case "genre":
			a.Genre = p.Text()
		case "year":
			a.Year = p.Int()
		case "tracks":
			a.Tracks = p.Int()
		}
	}
	return
}

func NewAlbums(r io.Reader) (albums []Album) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		albums = append(albums, ParseAlbum(s.Text()))
	}
	return
}

func NewAlbumsString(data string) []Album {
	return NewAlbums(strings.NewReader(data))
}
