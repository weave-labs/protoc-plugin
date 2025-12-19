package plugin

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func Run(r io.Reader, w io.Writer, generator Generator) error {
	if len(os.Args) > 1 {
		return fmt.Errorf("unknown args supplied: %v", os.Args[1:])
	}

	in, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(in, req); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	if len(req.GetFileToGenerate()) == 0 {
		return nil
	}

	var setFn func(name string, value string) error
	if generator.FlagSet() != nil {
		setFn = generator.FlagSet().Set
	}

	plugin, err := protogen.Options{
		ParamFunc: setFn,
	}.New(req)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	plugin.SupportedFeatures = generator.Features()
	plugin.SupportedEditionsMinimum, plugin.SupportedEditionsMaximum = generator.SupportedEditionsRange()

	if err := generator.Generate(plugin); err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	out, err := proto.Marshal(plugin.Response())
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
