package main

import (
	"debug/buildinfo"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"

	"github.com/admin-else/strom/cmd/strom/api"
	"github.com/admin-else/strom/cmd/strom/offline_uuid"
	"github.com/admin-else/strom/cmd/strom/packet_inspector"
	"github.com/admin-else/strom/cmd/strom/print_nbt"
	"github.com/admin-else/strom/cmd/strom/status"
)

var ExpectedASubcommandErr = errors.New("expected a subcommand")
var UnknownSubcommandErr = errors.New("unknown subcommand try 'strom help' for a list of subcommands")

func Version(args []string) (err error) {
	selfPath, err := os.Executable()
	if err != nil {
		return
	}
	info, err := buildinfo.ReadFile(selfPath)
	if err != nil {
		return
	}
	fmt.Println("strom cli ", info.Main.Version, " (", info.Main.Sum, ")")
	return
}

var subcommands = map[string]func(args []string) error{
	"offline-uuid": offline_uuid.Run,
	"print-nbt":    print_nbt.Run,
	"status":       status.Run,
	"packet-spy":   packet_inspector.Run,
	"api":          api.Run,
	"version":      Version,
}

func Help(args []string) (err error) {
	fmt.Println("Available subcommands", slices.Collect(maps.Keys(subcommands)))
	return
}

func mainE() (err error) {
	subcommands["help"] = Help // this is here because if it isn't tools/cli/main.go:17:5: initialization cycle for subcommands
	//	tools/cli/main.go:17:5: subcommands refers to Help
	//	tools/cli/main.go:22:6: Help refers to subcommands
	args := os.Args

	// Debug stuff
	slog.SetLogLoggerLevel(slog.LevelDebug)
	//args = []string{"TEST", "version"}

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
