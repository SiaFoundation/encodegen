package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("encodegen: ")
	flag.Usage = func() {
		os.Stderr.WriteString("Usage: encodgen [flags] -t T [directory]\n\nFlags:\n")
		flag.PrintDefaults()
	}

	dst := flag.String("o", "", "output file name; default srcdir/encoding.go")
	typs := flag.String("t", "", "comma-separated list of type names; to (optionally) override default unmarshaler allocation limit, add a colon and then the expression replacing it after; required")
	flag.Parse()
	args := flag.Args()

	if len(args) > 1 || *typs == "" {
		flag.Usage()
		os.Exit(2)
	}
	dir := "."
	if len(args) > 0 {
		dir = args[0]
		if !isDir(dir) {
			log.Fatalln(dir, "is not a directory")
		}
	}

	code, err := Generate(dir, strings.Split(*typs, ",")...)
	if err != nil {
		log.Fatal(err)
	}
	if *dst == "" {
		*dst = filepath.Join(dir, "encoding.go")
	}
	if err := ioutil.WriteFile(*dst, code, 0644); err != nil {
		log.Fatal(err)
	}
}
