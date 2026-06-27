package client

import "testing"

func TestDoDNS(t *testing.T) {
	type Case struct {
		Input string
		Host  string
		portN uint16
	}
	cases := []Case{
		{"localhost", "localhost", 25565},
		{"play.hypixel.net", "play.hypixel.net", 25565},
		{"1f2d.net", "tcpshield.1f2d.net", 25565},
	}
	for _, c := range cases {
		t.Logf("Testing %s", c.Input)
		_, host, portN, err := DoDns(c.Input)
		if err != nil {
			t.Error(err)
			continue
		}
		if host != c.Host || portN != c.portN {
			t.Errorf("expected %s:%d, got %s:%d", c.Host, c.portN, host, portN)
		}
	}
}
