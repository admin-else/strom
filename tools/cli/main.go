package main

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"

	"github.com/admin-else/strom/tools/cli/offline_uuid"
	"github.com/admin-else/strom/tools/cli/print_nbt"
	"github.com/admin-else/strom/tools/cli/status"
)

var ExpectedASubcommandErr = errors.New("expected a subcommand")
var UnknownSubcommandErr = errors.New("unknown subcommand")

var subcommands = map[string]func(args []string) error{
	"offline-uuid": offline_uuid.Run,
	"print-nbt":    print_nbt.Run,
	"status":       status.Run,
}

func Help(args []string) (err error) {
	fmt.Println("Available subcommands", slices.Collect(maps.Keys(subcommands)))
	return
}

func mainE() (err error) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	subcommands["help"] = Help // this is here because if it isn't tools/cli/main.go:17:5: initialization cycle for subcommands
	//	tools/cli/main.go:17:5: subcommands refers to Help
	//	tools/cli/main.go:22:6: Help refers to subcommands
	args := os.Args
	args = []string{"TEST", "status", "-addr", "tcpshield.1f2d.net:25565"}

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
