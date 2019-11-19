package discogs

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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

var db *sqlx.DB

// Terms takes a string and returns all the search terms from
// it for searching artists and albums combined
func Terms(s string) []string {
	return []string{}
}

// ArtistTerms takes a string and returns the all the search terms from
// it for searching artists
func ArtistTerms(s string) []string {
	return []string{}
}

// AlbumTerms takes a string and returns all the search terms from
// it for searching albums
func AlbumTerms(s string) []string {
	return []string{}
}

func init() {
	var err error
	db, err = sqlx.Connect(
		"sqlite3",
		"/home/haynes/projects/discogs/discogs.db")
	DieIfError(err)
}
