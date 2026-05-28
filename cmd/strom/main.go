package main

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"runtime/debug"
	"slices"

	"github.com/admin-else/strom/cmd/strom/api"
	"github.com/admin-else/strom/cmd/strom/data"
	"github.com/admin-else/strom/cmd/strom/extract_mca"
	"github.com/admin-else/strom/cmd/strom/offline_uuid"
	"github.com/admin-else/strom/cmd/strom/packet_info"
	"github.com/admin-else/strom/cmd/strom/packet_inspector"
	"github.com/admin-else/strom/cmd/strom/print_nbt"
	"github.com/admin-else/strom/cmd/strom/print_replay"
	"github.com/admin-else/strom/cmd/strom/status"
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
	"offline-uuid": offline_uuid.Run,
	"packet-info":  packet_info.Run,
	"print-nbt":    print_nbt.Run,
	"status":       status.Run,
	"packet-spy":   packet_inspector.Run,
	"api":          api.Run,
	"version":      Version,
	"extract-mca":  extract_mca.Run,
	"data":         data.Run,
	"print-replay": print_replay.Run,
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
