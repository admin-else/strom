package litematic_test

import (
	"embed"
	"testing"

	data2 "git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/level/litematic"
)

//go:embed testdata
var testData embed.FS

func TestLitematic(t *testing.T) {
	f, err := testData.Open("testdata/place.litematic")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer f.Close()
	s, err := litematic.LoadFromFile(f)
	if err != nil {
		t.Fatal(err)
		return
	}
	stateId, err := s.GetBlockAt("Cactus", 4, 0, 4)
	if err != nil {
		t.Fatal(err)
		return
	}
	version, err := data2.LookUpVersionByDataVersion(s.MinecraftDataVersion)
	if err != nil {
		t.Fatal(err)
		return
	}
	b, state, err := data2.FromBlockState(version.MinecraftVersion, stateId)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Block: %v, state: %v", b.Name, state)
}
