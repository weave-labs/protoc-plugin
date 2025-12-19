package plugin

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Generator interface {
	FlagSet() *flag.FlagSet
	Generate(gen *protogen.Plugin) error
	Name() string
	Version() string
	Features() uint64
	SupportedEditionsRange() (descriptorpb.Edition, descriptorpb.Edition)
}
