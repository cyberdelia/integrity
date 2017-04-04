package main

import (
	"archive/tar"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/cyberdelia/integrity"
)

var (
	file      string
	stdout    bool
	porcelain bool
	pattern   *regexp.Regexp
)

func init() {
	flag.StringVar(&file, "f", "", "File name of tar backup.")
	flag.BoolVar(&stdout, "c", false, "Copy input on standard output.")
	flag.BoolVar(&porcelain, "p", false, "Output in an easy-to-parse format for scripts.")
	pattern = regexp.MustCompile(`^(base|global)(/\d+)?/(\d+)$`)
}

func archiveReader(name string) (archive io.Reader, err error) {
	if file != "" {
		// Read from file
		return os.Open(file)
	}
	if stdout {
		// Read stdin and copy to stdout
		return io.TeeReader(os.Stdin, os.Stdout), nil
	}
	// Read form stdin
	return os.Stdin, nil
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	archive, err := archiveReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Copy(ioutil.Discard, archive)

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
		if err = integrity.Verify(tr); err != nil {
			log.Printf("%s: %s\n", h.Name, err)
		}
	}
}
