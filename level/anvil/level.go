package anvil

import (
	"io"

	"github.com/admin-else/strom/nbt"
)

type Level struct {
	DataPacks struct {
		Enabled  []string
		Disabled []string
	}
}

func ReadLevelData(r io.Reader) (l Level, err error) {
	var lDtatWrapped struct{ Data Level }
	err = nbt.ReadFile(r, &lDtatWrapped)
	if err != nil {
		return
	}
	l = lDtatWrapped.Data
	return
}
