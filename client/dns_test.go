package client

import "testing"

func TestDoDNS(t *testing.T) {
	type Case struct {
		Input  string
		Output string
	}
	for _, c := range []Case{{"localhost", "127.0.0.1:25565"}, {"1f2d.net", "50.114.4.126:25565"}} {
		t.Logf("Testing %s", c.Input)
		out, err := DoDns(c.Input)
		if err != nil {
			t.Error(err)
			continue
		}
		if out != c.Output {
			t.Errorf("Expected %s, got %s", c.Output, out)
		}
	}
}
