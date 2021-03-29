package codegen

import (
	"github.com/stretchr/testify/assert"
	"go.sia.tech/encodegen/internal/toolbox"
	"path"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {

	parent := path.Join(toolbox.CallerDirectory(3), "test")

	var useCases = []struct {
		description string
		options     *Options
		hasError    bool
	}{
		{
			description: "basic type to test imports",
			options: &Options{
				Source: path.Join(parent, "importedtype"),
				Types:  []Type{{Name: "Imported", ReuseMemory: true}, {Name: "Hash", ReuseMemory: true}},
				Dest:   path.Join(parent, "importedtype", "encoding.go"),
			},
		},
		{
			description: "basic struct code generation",
			options: &Options{
				Source: path.Join(parent, "basic_struct"),
				Types:  []Type{{Name: "Message", ReuseMemory: true}},
				Dest:   path.Join(parent, "basic_struct", "encoding.go"),
			},
		},
		{
			description: "struct composed of a bunch of different aliased types",
			options: &Options{
				Source: path.Join(parent, "alias_struct"),
				Types:  []Type{{Name: "Message", ReuseMemory: true}},
				Dest:   path.Join(parent, "alias_struct", "encoding.go"),
			},
		},
		{
			description: "struct with anonymous struct in it",
			options: &Options{
				Source: path.Join(parent, "embedded_struct"),
				Types:  []Type{{Name: "Message", ReuseMemory: true}},
				Dest:   path.Join(parent, "embedded_struct", "encoding.go"),
			},
		},
	}

	for _, useCase := range useCases {
		generator, err := NewGenerator(useCase.options)
		if err != nil {
			t.Fatal(err)
		}
		err = generator.Generate()
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		if !assert.Nil(t, err, useCase.description) {
			t.Fatal(err)
		}
	}

}
