package codegen

import (
	"fmt"
	"github.com/viant/toolbox"
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
	sliceTypes  map[string]string
	imports     map[string]bool
	filedInit   []string
	Pkg         string
	Code        string
	Init        string
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
	g.filedInit = []string{}
	g.imports = map[string]bool{}
	g.structTypes = map[string]string{}
	g.sliceTypes = map[string]string{}
	g.addImport(encodingPackage)
}

// NewGenerator creates a new generator with the given options
func NewGenerator(options *Options) *Generator {
	var g = &Generator{}
	// first we validate the flags
	if err := options.Validate(); err != nil {
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
	if err := g.readPackageCode(g.options.Source); err != nil {
		return err
	}

	// then we generate code for the types given
	for _, rootType := range g.options.Types {
		if err := g.generateStructCode(rootType); err != nil {
			return err
		}
	}

	g.Imports = strings.Join(toolbox.MapKeysToStringSlice(g.imports), "\n")
	return g.writeCode()
}

func (g *Generator) writeCode() error {
	var generatedCode = []string{}

	generatedCode = append(generatedCode, "")
	for _, key := range sortedKeys(g.sliceTypes) {
		code := g.sliceTypes[key]
		generatedCode = append(generatedCode, code)
	}
	generatedCode = append(generatedCode, "")
	for _, key := range sortedKeys(g.structTypes) {
		code := g.structTypes[key]
		generatedCode = append(generatedCode, code)
	}

	g.Code = strings.Join(generatedCode, "\n")

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

func (g *Generator) generatePrimitiveArray(field *Field) error {
	key := field.ComponentType + toolbox.AsString(field.IsPointerComponent)
	if _, ok := g.sliceTypes[key]; ok {
		return nil
	}
	code, err := expandBlockTemplate(baseTypeSlice, field)
	g.sliceTypes[key] = code
	return err
}

func (g *Generator) generateObjectArray(field *Field) error {
	if _, ok := g.sliceTypes[field.RawComponentType]; ok {
		return nil
	}

	return g.generateStructCode(field.ComponentType)
	// code, err := expandBlockTemplate(structTypeSlice, field)
	// if err != nil {
	// 	return err
	// }
	// g.sliceTypes[field.RawComponentType] = code
	// return err
}

func (g *Generator) generateTypedArray(field *Field) error {
	if _, ok := g.sliceTypes[field.RawComponentType]; ok {
		return nil
	}

	code, err := expandBlockTemplate(typeSlice, field)
	if err != nil {
		return err
	}
	g.sliceTypes[field.RawComponentType] = code
	return err
}

func (g *Generator) generateStructCode(structType string) error {
	typeInfo := g.Type(structType)
	if typeInfo == nil {
		return nil
	}

	_, hasCode := g.structTypes[structType]
	if hasCode {
		return nil
	}

	aStruct := NewStruct(typeInfo, g)
	code, err := aStruct.Generate()

	if err != nil {
		return err
	}

	g.structTypes[structType] = code
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
