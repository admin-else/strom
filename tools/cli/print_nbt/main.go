package print_nbt

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/admin-else/strom/nbt"
)

var (
	cmd      = flag.NewFlagSet("print-nbt", flag.ContinueOnError)
	FileFlag = cmd.String("file", "", "The NBT file to print")
)

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}
	f, err := os.Open(*FileFlag)
	if err != nil {
		return
	}
	defer f.Close()
	n, err := nbt.ReadFile(f)
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return
	}
	_, err = os.Stdout.Write(b)
	return
}
