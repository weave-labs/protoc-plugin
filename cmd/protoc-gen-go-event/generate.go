package main

import (
	"flag"
	"strings"

	"github.com/chanced/caps"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/weave-labs/protoc-plugin/pkg/plugin"
)

const (
	PluginName     = "protoc-gen-go-event"
	PluginVersion  = "v0.1.0"
	GenFileSuffix  = "_event.pb.go"
	topicSuffix    = "Topic"
	sqlTopicSuffix = "SQLTopic"
)

type Generator struct{}

type Event struct {
	TopicKey      string
	TopicValue    string
	SQLTopicKey   string
	SQLTopicValue string
}

func (g *Generator) FlagSet() *flag.FlagSet {
	return nil
}

func (g *Generator) Generate(plug *protogen.Plugin) error {
	for _, f := range plug.Files {
		if !f.Generate {
			continue
		}

		events := make([]Event, 0, len(f.Messages))
		for _, msg := range f.Messages {
			msgNameStr := string(msg.Desc.Name())
			if !strings.HasSuffix(strings.ToLower(msgNameStr), "event") {
				continue
			}

			var (
				pkgStr       = string(f.Desc.Package())
				msgNameSnake = caps.ToSnake(msgNameStr)
				msgNameCamel = caps.ToCamel(msgNameStr)
			)
			events = append(events, Event{
				TopicKey:      msgNameCamel + topicSuffix,
				TopicValue:    pkgStr + "." + msgNameSnake,
				SQLTopicKey:   msgNameCamel + sqlTopicSuffix,
				SQLTopicValue: makeSQLTopic(pkgStr, msgNameSnake),
			})
		}

		if len(events) == 0 {
			continue
		}

		genFile := plug.NewGeneratedFile(f.GeneratedFilenamePrefix+GenFileSuffix, f.GoImportPath)
		if err := ServiceTemplate.Execute(genFile, ServiceTemplateInput{
			PluginVersion:      g.Version(),
			ProtocVersion:      plugin.ProtocVersion(plug),
			FileDescriptorPath: f.Desc.Path(),
			PackageName:        string(f.GoPackageName),
			Events:             events,
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

const sqlPartsNeed = 3

func makeSQLTopic(pkg string, name string) string {
	pkgParts := strings.Split(pkg, ".")
	if len(pkgParts) < sqlPartsNeed {
		return ""
	}

	return strings.Join(pkgParts[2:], "_") + "." + name
}
