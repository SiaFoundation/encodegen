package codegen

import (
	"regexp"
	"fmt"
	"go.sia.tech/encodegen/internal/toolbox"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const encodingPackage = "go.sia.tech/encodegen/pkg/encodegen"

// Generator holds the content to generate the gojay code
type Generator struct {
	fileInfo    *toolbox.FileSetInfo
	types       map[string]string
	structTypes map[string]string
	imports     map[string]bool
	Pkg         string
	Code        string
	Imports     string
	options     *Options
}

// Returns the type from the the fileInfo
func (g *Generator) Type(typeName string) *toolbox.TypeInfo {
	return g.fileInfo.Type(typeName)
}

// addImport adds an import package to be printed on the generated code
func (g *Generator) addImport(pkg string) {
	g.imports[`"`+pkg+`"`] = true
}

// we initiate the variables containing the code to be generated
func (g *Generator) init() {
	g.imports = map[string]bool{}
	g.structTypes = map[string]string{}
	g.addImport(encodingPackage)
}

// NewGenerator creates a new generator with the given options
func NewGenerator(options *Options) *Generator {
	var g = &Generator{}
	// first we validate the flags
	err := options.Validate()
	if err != nil {
		panic(err)
	}
	g.options = options
	// we initiate the values on the generator
	g.init()
	return g
}

// Generate generates the gojay implementation code
func (g *Generator) Generate() error {
	// first we read the code from which we should find the types
	err := g.readPackageCode(g.options.Source)
	if err != nil {
		return err
	}

	// add whitespace trim character to the front of all templates https://golang.org/pkg/text/template/
	var re = regexp.MustCompile(`{{(if|else|end)(.*)}}`)
	var substitution = "{{- $1 $2}}"

	for key, value := range fieldTemplate {
		fieldTemplate[key] = re.ReplaceAllString(value, substitution)
	}
	for key, value := range blockTemplate {
		blockTemplate[key] = re.ReplaceAllString(value, substitution)
	}

	// then we generate code for the types given
	for _, rootType := range g.options.Types {
		err := g.generateStructCode(rootType)
		if err != nil {
			return err
		}
	}

	g.Imports = strings.Join(toolbox.MapKeysToStringSlice(g.imports), "\n")
	return g.writeCode()
}

func (g *Generator) writeCode() error {
	var generatedCode = []string{}

	generatedCode = append(generatedCode, "")
	generatedCode = append(generatedCode, "")
	for _, key := range sortedKeys(g.structTypes) {
		code := g.structTypes[key]
		generatedCode = append(generatedCode, code)
	}

	g.Code = strings.Join(generatedCode, "\n")

	// g.Code = strings.Replace(g.Code, "\n\n", "\n", -1)
	expandedCode, err := expandBlockTemplate(fileCode, g)
	if err != nil {
		return err
	}

	// fmt.Printf("UNFORMATTED\n:%s", string(expandedCode))

	code, err := format.Source([]byte(expandedCode))
	if err != nil {
		return err
	}

	// code destination is empty, we just print to stdout
	if g.options.Dest == "" {
		fmt.Print(string(code))
		return nil
	}

	return ioutil.WriteFile(g.options.Dest, code, 0644)
}

func (g *Generator) generateStructCode(structType Type) error {
	typeInfo := g.Type(structType.Name)
	if typeInfo == nil {
		return nil
	}

	_, hasCode := g.structTypes[structType.Name]
	if hasCode {
		return nil
	}

	aStruct := NewStruct(typeInfo, g)
	code, err := aStruct.Generate(structType.ReuseMemory)

	if err != nil {
		return err
	}

	g.structTypes[structType.Name] = code
	return nil
}

func (g *Generator) readPackageCode(pkgPath string) error {
	p, err := filepath.Abs(pkgPath)
	if err != nil {
		return err
	}

	f, err := os.Stat(p)
	if err != nil {
		// path/to/whatever does not exist
		return err
	}

	if !f.IsDir() {
		g.Pkg = filepath.Dir(p)
		dir, _ := filepath.Split(p)
		g.fileInfo, err = toolbox.NewFileSetInfo(dir)

	} else {
		g.Pkg = filepath.Base(p)
		g.fileInfo, err = toolbox.NewFileSetInfo(p)
	}

	// if Pkg flag is set use it
	if g.options.Pkg != "" {
		g.Pkg = g.options.Pkg
	}
	return err
}
