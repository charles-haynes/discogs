package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Artist struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Score int    `db:"score"`
}

type DiscogsAlbum struct {
	Artist   string  `db:"artist"`
	Album    string  `db:"album"`
	Year     int64   `db:"year"`
	Rating   float64 `db:"rating"`
	ArtistID int64   `db:"artistid"`
	AlbumID  int64   `db:"albumid"`
}

func Fatal(err error) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("FATAL %s:%d %s: %s\n", file, line, f.Name(), err)
	os.Exit(-1)
}

func main() {
	db, err := sqlx.Connect("sqlite3", "cross-seed.db")
	if err != nil {
		Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`ATTACH "discogs/discogs.db" AS discogs`)
	if err != nil {
		Fatal(err)
	}
	var artists []Artist
	err = db.Select(&artists,
		"SELECT id, name, score FROM artists ORDER BY score DESC")
	if err != nil {
		Fatal(err)
	}
	for _, artist := range artists {
		var albumsHave []string
		err = db.Select(&albumsHave,
			`SELECT DISTINCT(album) FROM have WHERE artist = ?`,
			artist)
		if err != nil {
			Fatal(err)
		}
		fmt.Printf("%s: have %d\n", artist, len(albumsHave))
		have := make(map[string]interface{}, len(albumsHave))
		for _, album := range albumsHave {
			fmt.Printf("-- %s\n", album)
			have[strings.ToLower(album)] = nil
		}
		var albums []DiscogsAlbum
		err = db.Select(&albums, `
SELECT da.*
FROM discogs.albums AS da
WHERE da.artist = ?
ORDER BY da.rating DESC
LIMIT ?`,
			artist, artist.Score)
		if err != nil {
			Fatal(err)
		}
		for _, album := range albums {
			if _, ok := have[strings.ToLower(album.Album)]; ok {
				fmt.Printf(
					"-- %s - %s (%04d)\n",
					artist, album.Album, album.Year)
				continue
			}
			fmt.Printf(
				"%s - %s (%04d)\n",
				artist, album.Album, album.Year)
		}
	}
}
