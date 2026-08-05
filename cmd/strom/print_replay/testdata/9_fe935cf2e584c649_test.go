package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v26_2"
)

// Error: unexpected EOF
func TestFailedPacketFE935CF2E584C649(t *testing.T) {
	p := &generated.PlayToClientPacketUpdateTime{}
	b := bytes.NewReader([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1c, 0x1a, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatalf("expected decode error but got nil")
	}
}
