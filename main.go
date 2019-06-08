package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type processFn func(path, format string) error

var counters = make(map[string]int)

func main() {
	srcDir := flag.String("in", ".", "Source directory")
	format := flag.String("out", "dist/%Y/%m/%d-%H%M%S",
		"Output format (use strftime symbols %Y, %y, %m, %d, %H, %M, %S)",
	)
	dryRun := flag.Bool("dryrun", false, "Dry run")
	flag.Parse()

	formatNorm := toTimeFormat(*format)
	var fn processFn
	n := -1
	if *dryRun {
		fmt.Println("Dry run. Examples of moved files:")
		fn = previewFile
		n = 10
	} else {
		fn = processFile
	}
	if err := processDir(*srcDir, formatNorm, n, fn); err != nil {
		log.Println(err)
	}
}

func processDir(dir, format string, max int, fn processFn) error {
	var i int
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i < 0 || i > max {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			i++
			return fn(path, format)
		}
		return nil
	})
}

func processFile(path, format string) error {
	if ext := strings.ToLower(filepath.Ext(path)); !isSupportedExt(ext) {
		return nil
	}

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	ts, err := getDate(src)
	if err != nil {
		return err
	}

	newPath := ts.Format(format) + "-" + strconv.Itoa(nextIndex(ts)) + ".jpg"
	dir := filepath.Dir(newPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	dest, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return err
	}

	return nil
}

func previewFile(path, format string) error {
	if ext := strings.ToLower(filepath.Ext(path)); !isSupportedExt(ext) {
		return nil
	}

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	ts, err := getDate(src)
	if err != nil {
		return err
	}

	newPath := ts.Format(format) + "-" + strconv.Itoa(nextIndex(ts)) + ".jpg"
	fmt.Printf("%s -> %s\n", path, newPath)

	return nil
}

func isSupportedExt(ext string) bool {
	switch ext {
	case ".jpg":
		return true
	default:
		return false
	}
}

func getDate(fd io.ReadSeeker) (time.Time, error) {
	ef, err := exif.Decode(fd)
	if err != nil {
		return time.Time{}, err
	}

	_, err = fd.Seek(0, 0)
	if err != nil {
		return time.Time{}, err
	}

	return ef.DateTime()
}

func nextIndex(ts time.Time) int {
	key := ts.Format("20060102150405")
	counters[key]++
	return counters[key]
}

func toTimeFormat(f string) string {
	f = strings.Replace(f, "%y", "2006", -1)
	f = strings.Replace(f, "%Y", "06", -1)
	f = strings.Replace(f, "%m", "01", -1)
	f = strings.Replace(f, "%d", "02", -1)

	f = strings.Replace(f, "%H", "15", -1)
	f = strings.Replace(f, "%M", "04", -1)
	f = strings.Replace(f, "%S", "05", -1)

	return f
}
