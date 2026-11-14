package checkrun

import "runtime/debug"

func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(devel)"
}
