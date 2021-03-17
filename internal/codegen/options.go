package codegen

import (
	"flag"
	"github.com/go-errors/errors"
	"go.sia.tech/encodegen/internal/toolbox/url"
	"strings"
)

type Type struct {
	Name        string
	ReuseMemory bool
}

type Options struct {
	Source  string
	Dest    string
	Types   []Type
	TagName string
	Pkg     string
}

func (o *Options) Validate() error {
	if o.Source == "" {
		return errors.New("Source was empty")
	}
	if len(o.Types) == 0 {
		return errors.New("Types was empty")
	}
	return nil
}

const (
	optionKeySource = "s"
	optionKeyDest   = "o"
	optionKeyTypes  = "t"
	optionKeyPkg    = "pkg"
)

//NewOptionsWithFlagSet creates a new options for the supplide flagset
func NewOptionsWithFlagSet(set *flag.FlagSet) *Options {
	var result = &Options{}
	var reuseMemory bool

	result.Dest = set.Lookup(optionKeyDest).Value.String()
	result.Source = set.Lookup(optionKeySource).Value.String()

	types := strings.Split(set.Lookup(optionKeyTypes).Value.String(), ",")
	if len(types) == 1 {
		result.Types = append(result.Types, Type{
			Name:        types[0],
			ReuseMemory: false,
		})
	} else if len(types) != 0 {
		for i := 0; i < len(types)-1; i++ {
			if types[i] == "true" || types[i] == "false" {
				continue
			}

			reuseMemory = false
			nextArgument := types[i+1]
			if nextArgument == "true" {
				reuseMemory = true
			}
			result.Types = append(result.Types, Type{
				Name: types[i],
				ReuseMemory: reuseMemory,
			})

		}
	}
	result.Pkg = set.Lookup(optionKeyPkg).Value.String()
	if result.Source == "" {
		result.Source = url.NewResource(".").ParsedURL.Path
	}
	return result
}
