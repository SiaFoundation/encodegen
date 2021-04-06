package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	pkg := flag.String("pkg", "", "name of target package")
	dst := flag.String("o", "", "destination of generated code (optional; omit for stdout)")
	typs := flag.String("t", "", "types to generate, comma separated")
	flag.Parse()

	code, err := Generate(*pkg, strings.Split(*typs, ",")...)
	if err != nil {
		log.Fatal(err)
	}
	if *dst == "" {
		fmt.Println(code)
		return
	}
	if err := ioutil.WriteFile(*dst, []byte(code), 0644); err != nil {
		log.Fatal(err)
	}
}
