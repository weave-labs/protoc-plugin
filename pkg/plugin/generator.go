package plugin

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
)

type Generator interface {
	FlagSet() *flag.FlagSet
	Generate(gen *protogen.Plugin) error
	Name() string
	Version() string
}
