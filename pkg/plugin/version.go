package plugin

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func ProtocVersion(plugin *protogen.Plugin) string {
	ver := plugin.Request.GetCompilerVersion()
	if ver == nil {
		return "(unknown)"
	}

	var suffix string
	if s := ver.GetSuffix(); s != "" {
		suffix = "-" + s
	}

	return fmt.Sprintf("v%d.%d.%d%s", ver.GetMajor(), ver.GetMinor(), ver.GetPatch(), suffix)
}
