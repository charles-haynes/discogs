package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

func Human(i int64) string {
	v := i
	if v < 0 {
		v = -v
	}
	if v < 1000 {
		return fmt.Sprintf("%d", i)
	}
	if v < 1000000 {
		return fmt.Sprintf("%7.2fk", float64(i)/1000.0)
	}
	if v < 1000000000 {
		return fmt.Sprintf("%7.2fm", float64(i)/1000000.0)
	}
	return fmt.Sprintf("%7.2fg", float64(i)/1000000000.0)
}

func FatalN(err error, n int) {
	rpc := make([]uintptr, 1)
	runtime.Callers(n, rpc)
	frame, _ := runtime.CallersFrames(rpc).Next()
	fmt.Printf("FATAL %s:%d %s: %s\n",
		frame.File, frame.Line, frame.Function, err)
	os.Exit(-1)
}

func Fatal(err error) {
	FatalN(err, 3)
}

func DieIfError(err error) {
	if err != nil {
		FatalN(err, 3)
	}
}

type Time struct{ time.Time }

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var timeString string
	if err := d.DecodeElement(&timeString, &start); err != nil {
		return err
	}
	ts, err := time.Parse("2006-01-02T15:04:05.000Z", timeString)
	if err != nil {
		return err
	}
	*t = Time{ts}
	return nil
}

type Contents struct {
	Key          string
	LastModified Time
	ETag         string
	Size         int64
	StorageCass  string
}

// ByModified implements sort.Interface for []Contents based on
// the LastModified field
type ByModified []Contents

func (c ByModified) Len() int      { return len(c) }
func (c ByModified) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByModified) Less(i, j int) bool {
	return c[i].LastModified.Before(c[j].LastModified.Time)
}

type Result struct {
	Name        string
	Prefix      string
	MaxKeys     int
	Delimiter   string
	IsTruncated bool
	Contents    []Contents
}

const urlPrefix = "https://discogs-data.s3-us-west-2.amazonaws.com/"

func main() {
	url := fmt.Sprintf("%s?delimiter=/&prefix=data/%s/",
		urlPrefix, time.Now().Format("2006"))
	resp, err := http.Get(url)
	DieIfError(err)
	defer resp.Body.Close()
	var r Result
	DieIfError(xml.NewDecoder(resp.Body).Decode(&r))
	sort.Sort(ByModified(r.Contents))
	for i := len(r.Contents) - 1; i >= len(r.Contents)-5; i-- {
		name := r.Contents[i].Key
		fmt.Printf("fetching %9s %s\n",
			Human(r.Contents[i].Size), name)
		start := time.Now()
		resp, err := http.Get(urlPrefix + name)
		DieIfError(err)
		f, err := os.Create(filepath.Base(name))
		DieIfError(err)
		_, err = io.Copy(f, resp.Body)
		DieIfError(err)
		resp.Body.Close()
		f.Close()
		end := time.Since(start)
		fmt.Printf("fetch took %s %s/s\n",
			end,
			Human(r.Contents[i].Size/int64(math.Round(end.Seconds()))))
	}
}
