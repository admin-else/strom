package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/admin-else/strom/nbt"
)

var (
	FileFlag = flag.String("file", "", "The NBT file to print")
)

func mainE() (err error) {
	flag.Parse()
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

func main() {
	err := mainE()
	if err != nil {
		panic(err)
	}
}
