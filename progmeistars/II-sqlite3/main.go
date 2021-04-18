package main

import (
	_ "embed"
	"log"
	"os"
	"text/template"
)

//go:embed music.dat
var albumsData string

const outputTemplate = `{{$g := .Grouping}}
Total albums: {{len .Albums}}

+----------------+--------+
| Genre          | Albums |
+----------------+--------+
{{- range .Grouping.Genres}}
| {{printf "%-14.14s " . -}}
| {{$g.AlbumsForGenre . | len | printf "% 6d " -}}
|{{end}}
+----------------+--------+

+------+--------+
| Year | Albums |
+------+--------+
{{- range .Grouping.Years}}
| {{printf "%4d " . -}}
| {{$g.AlbumsForYear . | len | printf "% 6d " -}}
|{{end}}
+------+--------+
`

func main() {
	albums := NewAlbumsString(albumsData)
	grouping, err := NewGrouping(albums)
	if err != nil {
		log.Fatalf("new grouping: %v", err)
	}

	tmpl := template.Must(template.New("albums").Parse(outputTemplate))
	tmpl.Execute(os.Stdout, struct {
		Albums   []Album
		Grouping Grouping
	}{
		Albums:   albums,
		Grouping: grouping,
	})

}
