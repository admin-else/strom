package print_nbt

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	nbt2 "github.com/admin-else/strom/mc/nbt"
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
	if *FileFlag == "" {
		return fmt.Errorf("no file specified")
	}
	f, err := os.Open(*FileFlag)
	if err != nil {
		return
	}
	defer f.Close()
	n, err := nbt2.ReadUnstructuredFile(f)
	if err != nil {
		return
	}
	var b []byte
	if *HashFormatFlag {
		bb := bytes.NewBuffer(nil)
		err = nbt2.PrintSNBTAny(n.Value, bb)
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
