package extract_mca

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/admin-else/strom/mc/level/anvil"
	"github.com/admin-else/strom/mc/nbt"
)

var (
	cmd        = flag.NewFlagSet("print-nbt", flag.ContinueOnError)
	FileFlag   = cmd.String("file", "", "The MCA file to extract")
	IndexXFlag = cmd.Int("x", 0, "The x chunk coordinate")
	IndexZFlag = cmd.Int("z", 0, "The z chunk coordinate")
	ToFlag     = cmd.String("to", "mca-chunk.nbt", "The output file")
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
	mca, err := anvil.ReadChunkHeader(f)
	if err != nil {
		return
	}
	t, err := mca.TimeAt(f, int32(*IndexXFlag), int32(*IndexZFlag))
	if err != nil {
		return
	}
	slog.Info("Chunk last written", "time", t)
	n, err := mca.ChunkAt(f, int32(*IndexXFlag), int32(*IndexZFlag))
	if err != nil {
		return
	}
	err = nbt.WriteUnstructuredFilePath(*ToFlag, n)
	if err != nil {
		return
	}
	return
}
