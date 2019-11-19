package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/irlndts/go-discogs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Artist struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Score int    `db:"score"`
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
	db, err := sqlx.Connect("sqlite3", "discogs.db")
	if err != nil {
		Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`ATTACH "../cross-seed.db" as cs;`)
	if err != nil {
		Fatal(err)
	}
	d, err := discogs.NewClient(&discogs.Options{
		UserAgent: "bz.ceh.updaterating/0.0.1",
		Token:     "syxidsTOJubHYggGpAWWerYMciDKHtcpPadiIYIr",
	})
	if err != nil {
		Fatal(err)
	}
	var artists []Artist
	if err := db.Select(
		&artists,
		`SELECT * FROM cs.artists WITH score > 0 ORDER BY score DESC`); err != nil {
		Fatal(err)
	}
	fmt.Println("PRAGMA foreign_keys=OFF;")
	fmt.Println("BEGIN TRANSACTION;")
	for _, a := range artists {
		pages := 1
		for page := 1; page <= pages; page++ {
			fmt.Printf("-- %s page %d/%d\n", a.Name, page, pages)
			time.Sleep(1 * time.Second)
			res, _ := d.Search.Search(discogs.SearchRequest{
				Artist:  a.Name,
				Type:    "master",
				Page:    page,
				PerPage: 100,
			})
			pages = res.Pagination.Pages

			for _, r := range res.Results {
				fmt.Printf(
					"UPDATE masters SET rating = %d WHERE id=%d;\n",
					r.Community.Have, r.ID)
			}
		}
	}
	fmt.Println("COMMIT;")
}
