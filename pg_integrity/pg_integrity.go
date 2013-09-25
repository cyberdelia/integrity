package main

import (
	"archive/tar"
	"flag"
	"github.com/cyberdelia/integrity"
	"io"
	"log"
	"os"
	"regexp"
)

var (
	file    string
	pattern *regexp.Regexp
)

func init() {
	flag.StringVar(&file, "f", "", "File name of tar backup")
	pattern = regexp.MustCompile(`^(base|global)(/\d+)?/(\d+)$`)
}

func archiveReader(name string) (archive io.ReadCloser, err error) {
	if file != "" {
		return os.Open(file)
	}
	return os.Stdin, nil
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	archive, err := archiveReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	tr := tar.NewReader(archive)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if h.Typeflag != tar.TypeReg {
			// Ignore directories
			continue
		}
		if !pattern.MatchString(h.Name) {
			// Ignore non-pages files
			continue
		}
		err = integrity.Verify(tr)
		if err != nil {
			log.Printf("%s: %s\n", h.Name, err)
		}
	}
}
