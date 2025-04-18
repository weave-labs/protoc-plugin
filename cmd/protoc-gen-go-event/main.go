package main

import (
	"log"
	"os"

	"github.com/weave-labs/protoc-plugin/pkg/plugin"
)

func main() {
	generator := new(Generator)
	if err := plugin.Run(os.Stdin, os.Stdout, generator); err != nil {
		log.Fatalf("%s: %v", generator.Name(), err)
	}
}
