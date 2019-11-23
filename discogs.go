package discogs

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "/home/haynes/projects/discogs/discogs.db"

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

// DB holds the internal state of a discogs database
type DB struct {
	*sqlx.DB
	mean, sd float64
}

// ArtistAlbum is just an artist and album for use in queries and other types
type ArtistAlbum struct {
	Artist string
	Album  string
}

func unicode61(r rune) bool {
	rt := []*unicode.RangeTable{
		unicode.L,
		unicode.N,
		unicode.Co,
	}
	return !unicode.IsOneOf(rt, r)
}

// Terms takes a string and returns all the search terms from
// it for searching artists and albums combined
func (d *DB) Terms(query string) []string {
	words := strings.FieldsFunc(query, unicode61)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	if len(words) == 1 {
		return words
	}
	var r []struct {
		Term  string
		Delta float64
	}
	q, a, err := sqlx.In(`
SELECT term, abs(?-cnt) as delta FROM fts_v WHERE term IN (?);`,
		d.mean, words)
	DieIfError(err)
	DieIfError(d.Select(&r, q, a...))
	for i := range r {
		if r[i].Delta < d.sd {
			r[i].Delta = 0.1
		}
	}
	switch len(r) {
	case 0:
		return []string{}
	case 1:
		return []string{r[0].Term}
	}
	const gamma = 0.8
	weight := gamma
	r[1].Delta /= weight
	if r[1].Delta < r[0].Delta {
		r[0], r[1] = r[1], r[0]
	}
	for _, dx := range r[2:] {
		weight *= gamma
		r[1].Delta /= weight
		if r[1].Delta > dx.Delta {
			if r[0].Delta > dx.Delta {
				r[0], r[1] = dx, r[0]
			} else {
				r[1] = dx
			}
		}
	}
	return []string{r[0].Term, r[1].Term}
}

// NewDB returns a new discogs db instance
func NewDB(db *sqlx.DB) DB {
	var d DB
	d.DB = db
	DieIfError(d.Get(&d.mean, `SELECT mean FROM stats`))
	DieIfError(d.Get(&d.sd, `SELECT sd FROM stats`))
	return d
}

// New returns a new discogs db instance using a fixed sqlite db
func New() DB {
	db, err := sqlx.Connect("sqlite3", dbPath)
	DieIfError(err)
	return NewDB(db)
}
