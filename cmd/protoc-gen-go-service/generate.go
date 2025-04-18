package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/weave-labs/protoc-plugin/pkg/plugin"
)

const (
	PluginName    = "protoc-gen-go-service"
	PluginVersion = "v0.1.0"
	GenFileSuffix = "_service.pb.go"
)

type Generator struct{}

func (g *Generator) FlagSet() *flag.FlagSet {
	return nil
}

func (g *Generator) Generate(plug *protogen.Plugin) error {
	for _, f := range plug.Files {
		if !f.Generate || len(f.Services) == 0 {
			continue
		}

		genFile := plug.NewGeneratedFile(f.GeneratedFilenamePrefix+GenFileSuffix, f.GoImportPath)

		servicesNames := make([]string, len(f.Services))
		for i, svc := range f.Services {
			servicesNames[i] = svc.GoName
		}

		if err := ServiceTemplate.Execute(genFile, ServiceTemplateInput{
			PluginVersion:      g.Version(),
			ProtocVersion:      plugin.ProtocVersion(plug),
			FileDescriptorPath: f.Desc.Path(),
			PackageName:        string(f.GoPackageName),
			ServiceNames:       servicesNames,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) Name() string {
	return PluginName
}

func (g *Generator) Version() string {
	return PluginVersion
}
