package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Master struct {
	ID      int      `xml:"id,attr"`
	Title   string   `xml:"title"`
	Artists []string `xml:"artists>artist>name"`
}

// FatalN dies after printing info about it's Nth parent caller
// it's the low level function used to print information about where
// a fatal error happened
func FatalN(err error, n int) {
	rpc := make([]uintptr, 1)
	runtime.Callers(n, rpc)
	frame, _ := runtime.CallersFrames(rpc).Next()
	fmt.Printf("FATAL %s:%d %s: %s\n",
		frame.File, frame.Line, frame.Function, err)
	os.Exit(-1)
}

// Fatal dies printing the location from where it was invoked
func Fatal(err error) {
	FatalN(err, 3)
}

// DieIfError dies if the error param is non-nil, printing the location
// in it's caller if there's an error. Used for checking the error status
// of calls, and dying with the location and error if there was an error
func DieIfError(err error) {
	if err != nil {
		FatalN(err, 3)
	}
}

func DecodeMaster(tx *sqlx.Tx, d *xml.Decoder, se xml.StartElement) {
	m := Master{}
	d.DecodeElement(&m, &se)
	_ = tx.MustExec(`INSERT INTO album_fts (rowid, artist, album) VALUES (?,?,?);`,
		m.ID, strings.Join(m.Artists, ","), m.Title)
}

func main() {
	var f io.ReadCloser = os.Stdin
	var err error
	if len(os.Args) > 1 {
		f, err = os.Open(os.Args[1])
		DieIfError(err)
		defer f.Close()
	}
	f, err = gzip.NewReader(f)
	DieIfError(err)
	fmt.Printf("Connecting")
	start := time.Now()
	db := sqlx.MustConnect("sqlite3", "/home/haynes/projects/cross-seed/discogs/discogs.db")
	fmt.Printf(" took %s\n", time.Since(start))
	fmt.Printf("Dropping/Creating")
	start = time.Now()
	_ = db.MustExec(`
DROP TABLE IF EXISTS album_fts;
CREATE VIRTUAL TABLE album_fts USING fts5(artist, album);
`)
	fmt.Printf(" took %s\n", time.Since(start))
	tx := db.MustBegin()
	defer tx.Rollback()
	d := xml.NewDecoder(f)
	tick := time.Tick(1 * time.Second)
	start = time.Now()
	var c int64 = 0
	for {
		select {
		case <-tick:
			fmt.Printf("\r%d %d/s",
				c,
				c/int64(time.Since(start)/time.Second))
		default:
		}
		t, err := d.Token()
		if err != nil || t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "master" {
				c++
				DecodeMaster(tx, d, se)
			}
		}
	}
	fmt.Printf("\n%d records in %s, %d/s\n",
		c,
		time.Since(start),
		c/int64(time.Since(start)/time.Second))
	DieIfError(tx.Commit())
	fmt.Printf("Optimizing/Analyzing")
	start = time.Now()
	_ = db.MustExec(`
INSERT INTO album_fts(album_fts)
  VALUES('optimize');
ANALYZE;
`)
	fmt.Printf(" took %s\n", time.Since(start))
}
