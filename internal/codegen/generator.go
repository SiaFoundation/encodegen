package codegen

import (
	"fmt"
	"go.sia.tech/encodegen/internal/toolbox"
	"go/format"
	"go/types"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const encodingPackage = "go.sia.tech/encodegen/pkg/encodegen"
const nebulousEncodingPackage = "gitlab.com/NebulousLabs/encoding"
const sizeofPrefix = "EncodegenSizeofEmpty"

type Import struct {
	Path    string
	Types   []string
	Enabled bool
}

// Generator holds the content to generate the gojay code
type Generator struct {
	fileInfo    *toolbox.FileSetInfo
	types       map[string]string
	structTypes map[string]string
	imports     map[string]Import
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
func (g *Generator) addImport(shortName string, path string, enabled bool) {
	g.imports[shortName] = Import{
		Path:    path,
		Enabled: enabled,
	}
}

// we initiate the variables containing the code to be generated
func (g *Generator) init() {
	g.imports = make(map[string]Import)
	g.types = make(map[string]string)
	g.structTypes = make(map[string]string)
}

// NewGenerator creates a new generator with the given options
func NewGenerator(options *Options) (*Generator, error) {
	var g = &Generator{}
	// first we validate the flags
	err := options.Validate()
	if err != nil {
		return nil, err
	}
	g.options = options
	// we initiate the values on the generator
	g.init()
	return g, nil
}

// Generate generates the gojay implementation code
func (g *Generator) Generate() error {
	// first we read the code from which we should find the types
	err := g.readPackageCode(g.options.Source)
	if err != nil {
		return err
	}

	// add whitespace trim character to the front of all templates - see https://golang.org/pkg/text/template/
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

	g.addImport("encodegen", encodingPackage, true)
	g.addImport("encoding", nebulousEncodingPackage, true)
	for _, shortName := range toolbox.MapKeysToStringSlice(g.imports) {
		currentImport := g.imports[shortName]
		if currentImport.Enabled {
			pkgs, err := packages.Load(&packages.Config{Mode: packages.LoadTypes}, g.imports[shortName].Path)
			// invalid package
			if err != nil {
				return err
			}

			// ensure imported types have marshal/unmarshalbuffer
			for _, pkg := range pkgs {
				for _, importedType := range currentImport.Types {
					if pkg.Types != nil && pkg.Types.Scope() != nil {
						lookedUpType := pkg.Types.Scope().Lookup(importedType)
						if lookedUpType == nil || lookedUpType.Type() == nil {
							return fmt.Errorf("We could not find type %+s in scope %+s", importedType, currentImport.Path)
						}

						marshalBufferObject, _, _ := types.LookupFieldOrMethod(lookedUpType.Type(), true, pkg.Types, "MarshalBuffer")
						if marshalBufferObject == nil {
							return fmt.Errorf("This type (%s.%s) does not implement MarshalBuffer", currentImport.Path, importedType)
						}
						unmarshalBufferObject, _, _ := types.LookupFieldOrMethod(lookedUpType.Type(), true, pkg.Types, "UnmarshalBuffer")
						if unmarshalBufferObject == nil {
							return fmt.Errorf("This type (%s.%s) does not implement UnmarshalBuffer", currentImport.Path, importedType)
						}
					} else {
						return fmt.Errorf("This type does not exist or its scope cannot be assessed.")
					}
				}
			}
			g.Imports += fmt.Sprintf(`%s "%s"%s`, shortName, currentImport.Path, "\n")
		}
	}

	return g.writeCode()
}

func (g *Generator) writeCode() error {
	var generatedCode = []string{}
	var structInit string

	generatedCode = append(generatedCode, "var (")
	for _, key := range sortedKeys(g.structTypes) {
		keyType := g.Type(key)
		primitiveType, primitiveDerived := isPrimitiveDerived(g.fileInfo, keyType)
		if primitiveDerived {
			structInit = "(" + primitiveType.ResetString + ")"
		} else {
			structInit = "{}"
		}
		generatedCode = append(generatedCode, fmt.Sprintf("%s%s = len(encoding.Marshal(%s%s))", sizeofPrefix, key, key, structInit))
	}
	generatedCode = append(generatedCode, ")")

	for _, key := range sortedKeys(g.structTypes) {
		code := g.structTypes[key]
		generatedCode = append(generatedCode, code)
	}

	g.Code = strings.Join(generatedCode, "\n")

	expandedCode, err := expandBlockTemplate(fileCode, g)
	if err != nil {
		return err
	}

	// log.Printf("UNFORMATTED\n:%s", string(expandedCode))

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

	aStruct := NewStruct(typeInfo, g, structType.ReuseMemory)
	code, err := aStruct.Generate()

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

	// add all imports to the map.  if they are referenced by a field they will be enabled
	for _, fileInfo := range g.fileInfo.FilesInfo() {
		for shortName, importData := range fileInfo.Imports {
			if shortName != "" {
				g.addImport(shortName, importData, false)
			}
		}
	}

	if g.options.Pkg != "" {
		g.Pkg = g.options.Pkg
	}
	return err

}
