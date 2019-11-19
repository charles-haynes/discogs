module github.com/charles-haynes/cross-seed

go 1.13

replace github.com/charles-haynes/discogs => ../discogs
replace github.com/charles-haynes/gazelle => ../gazelle
replace github.com/charles-haynes/whatapi => ../whatapi

require (
	github.com/anacrolix/tagflag v0.0.0-20180803105420-3a8ff5428f76
	github.com/anacrolix/torrent v1.5.2
	github.com/antchfx/htmlquery v1.0.0
	github.com/antchfx/xpath v1.0.0 // indirect
	github.com/bradfitz/iter v0.0.0-20190303215204-33e6a9893b0c
	github.com/charles-haynes/discogs latest
	github.com/charles-haynes/gazelle v0.1.0
	github.com/charles-haynes/html2bbcode v0.0.0-20191010140350-63d1fbb9d7d7
	github.com/charles-haynes/munkres v0.0.0-20190922052047-4487b45c0a2c
	github.com/charles-haynes/strsim v0.0.0-20191011181331-ef9ead4980ee
	github.com/charles-haynes/whatapi v0.0.14
	github.com/cosiner/argv v0.0.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0
	github.com/go-delve/delve v1.3.2 // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/irlndts/go-discogs v0.0.0-20181211134731-618b88263431
	github.com/jmoiron/sqlx v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mattn/go-runewidth v0.0.5 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/peterh/liner v1.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xrash/smetrics v0.0.0-20170218160415-a3153f7040e9
	go.starlark.net v0.0.0-20191021185836-28350e608555 // indirect
	golang.org/x/arch v0.0.0-20190927153633-4e8777c89be4 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/net v0.0.0-20191109021931-daa7c04131f5
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/yaml.v2 v2.2.4
)
