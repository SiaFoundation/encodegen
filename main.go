package main

import (
	"flag"
	"go.sia.tech/encodegen/internal/codegen"
	"log"
)

var pkg = flag.String("pkg", "", "the package name of the generated file")
var dst = flag.String("o", "", "destination file to output generated code")
var src = flag.String("s", "", "Source dir or file (absolute or relative path), omit for stdout")
var types = flag.String("t", "", `Types to generate, comma separated.  To enable memory reuse, put "true" after a type, e.g. Message,true,SubMessage,SubMessage2.  Memory reuse defaults to false if not specified.`)

func main() {
	flag.Parse()
	options := codegen.NewOptionsWithFlagSet(flag.CommandLine)
	generator, err := codegen.NewGenerator(options)
	if err != nil {
		log.Fatal(err)
	}
	err = generator.Generate()
	if err != nil {
		log.Fatal(err)
	}
}
