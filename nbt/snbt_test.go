package nbt

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/top_0_20.nbt
var topSNBT []byte

func TestPrintSNBTAny(t *testing.T) {
	n, err := ReadUnstructuredFile(bytes.NewReader(topSNBT))
	if err != nil {
		panic(err)
	}
	bb := bytes.NewBuffer(nil)
	err = PrintSNBTAny(n.Value, bb)
	if err != nil {
		panic(err)
	}
	t.Log(bb.String())
}
