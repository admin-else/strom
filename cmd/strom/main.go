package main

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"runtime/debug"
	"slices"

	"git.anygate.cloud/anygatecloud/strom/cmd/strom/api"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/data"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/extract_mca"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/mca_to_tmcpr"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/messenger"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/offline_uuid"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/packet_info"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/packet_inspector"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/print_nbt"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/print_replay"
	// "git.anygate.cloud/anygatecloud/strom/cmd/strom/raw_capture" // removed: use tcpflow instead for raw TCP capture
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/raw_to_tmcpr"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/serve_tmcpr"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/serve_world"
	"git.anygate.cloud/anygatecloud/strom/cmd/strom/status"
)

var ExpectedASubcommandErr = errors.New("expected a subcommand")
var UnknownSubcommandErr = errors.New("unknown subcommand try 'strom help' for a list of subcommands")

func ToolVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "devel"
	}
	return info.Main.Version
}

func ToolSum() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Sum
}

func Version(args []string) error {
	fmt.Println("strom cli ", ToolVersion(), " (", ToolSum(), ")")
	return nil
}

var subcommands = map[string]func(args []string) error{
	"offline-uuid":         offline_uuid.Run,
	"packet-info":          packet_info.Run,
	"print-nbt":            print_nbt.Run,
	"status":               status.Run,
	"packet-spy":           packet_inspector.Run,
	"api":                  api.Run,
	"version":              Version,
	"extract-mca":          extract_mca.Run,
	"data":                 data.Run,
	"print-replay":         print_replay.Run,
	// "raw-capture":          raw_capture.Run, // removed: use tcpflow instead for raw TCP capture, then raw-to-tmcpr
	"raw-to-tmcpr":         raw_to_tmcpr.Run,
	"serve-world":          serve_world.Run,
	"mca-to-tmcpr":         mca_to_tmcpr.Run,
	"serve-tmcpr":          serve_tmcpr.Run,
	"messenger":            messenger.Run,
}

func Help(args []string) (err error) {
	subcommandNames := slices.Collect(maps.Keys(subcommands))
	slices.Sort(subcommandNames)
	fmt.Println("Available subcommands", subcommandNames)
	return
}

func mainE() (err error) {
	subcommands["help"] = Help
	args := os.Args
	//args = []string{"TEST", "messenger", "-addr", "localhost"}

	if len(args) < 2 {
		err = ExpectedASubcommandErr
		return
	}

	f, ok := subcommands[args[1]]
	if !ok {
		err = UnknownSubcommandErr
		return
	}
	err = f(args[2:])
	return
}

func main() {
	err := mainE()
	if err != nil {
		slog.Error("exiting", "err", err)
	}
}
