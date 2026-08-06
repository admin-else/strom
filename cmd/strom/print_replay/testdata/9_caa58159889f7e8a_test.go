package failed_packets_test

import (
	"bytes"
	"testing"

	generated "git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v26_2"
)

// Error: unexpected EOF
func TestFailedPacketCAA58159889F7E8A(t *testing.T) {
	p := &generated.PlayToClientPacketUpdateTime{}
	b := bytes.NewReader([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1c, 0x2e, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatalf("expected decode error but got nil")
	}
}
