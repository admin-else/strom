package failed_packets_test

import (
	"bytes"
	"testing"

	generated "github.com/admin-else/strom/mc/proto_generated/v26_2"
)

// Error: to do
func TestFailedPacketB3E03D743AA3FD86(t *testing.T) {
	p := &generated.PlayToClientPacketEntityEquipment{}
	b := bytes.NewReader([]byte{0xee, 0x1, 0x0, 0x1, 0x9a, 0x7, 0x0, 0x0})
	err := p.Decode(b)
	if err == nil {
		t.Fatalf("expected decode error but got nil")
	}
}
