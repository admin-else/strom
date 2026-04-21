package print_nbt

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"

	"github.com/admin-else/strom/nbt"
)

var (
	cmd            = flag.NewFlagSet("print-nbt", flag.ContinueOnError)
	FileFlag       = cmd.String("file", "", "The NBT file to print")
	HashFormatFlag = cmd.Bool("snbt", false, "will print with go %#v")
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
	n, err := nbt.ReadUnstructuredFile(f)
	if err != nil {
		return
	}
	var b []byte
	if *HashFormatFlag {
		bb := bytes.NewBuffer(nil)
		err = nbt.PrintSNBTAny(n.Value, bb)
		if err != nil {
			return
		}
		b = bb.Bytes()
		//b = []byte(strings.ReplaceAll(string(b), ",", ",\n"))
	} else {
		b, err = json.MarshalIndent(n, "", "  ")
	}
	if err != nil {
		return
	}
	_, err = os.Stdout.Write(b)
	return
}
