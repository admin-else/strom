package proto_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var nanIgnorer32 = cmp.Comparer(func(x, y float32) bool {
	return (math.IsNaN(float64(x)) && math.IsNaN(float64(y))) || x == y
})

var nanIgnorer64 = cmp.Comparer(func(x, y float64) bool {
	return (math.IsNaN(x) && math.IsNaN(y)) || x == y
})

func NewStateFuzzer(state proto_base.State, actor proto_base.Actor) func(t *testing.T, inbytes []byte) {
	return func(t *testing.T, inbytes []byte) {
		b1 := bytes.NewBuffer(inbytes)
		b2 := bytes.NewBuffer(nil)
		conn := proto.NewConn()
		conn.SetState(state)
		conn.Actor = actor
		conn.R = b1
		conn.W = b2
		_ = conn.SetVersion("1.21.8") // cant fail
		p, err := conn.Receive()
		if err != nil {
			t.Skip(err)
			return
		}
		if b1.Len() != 0 {
			t.Skip("not all bytes read")
			return
		}
		err = conn.Send(p)
		if err != nil {
			t.Fatal(err)
			return
		}
		b1.Reset()
		conn.R = b2
		conn.W = b1
		p2, err := conn.Receive()
		if err != nil {
			t.Fatal(err)
			return
		}
		if !cmp.Equal(p, p2, nanIgnorer32, nanIgnorer64, cmpopts.IgnoreFields(proto.UnCodablePacket{}, "Err")) {
			t.Error("decoded and encoded decoded data are not equal",
				cmp.Diff(p, p2, nanIgnorer32, nanIgnorer64, cmpopts.IgnoreFields(proto.UnCodablePacket{}, "Err")),
			)
		}
	}
}

func FuzzConnHandshakingToServer(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Handshaking, proto_base.Servee))
}

func FuzzConnHandshakingToClient(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Handshaking, proto_base.Client))
}

func FuzzConnStatusToServer(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Status, proto_base.Servee))
}

func FuzzConnStatusToClient(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Status, proto_base.Client))
}

func FuzzConnLoginToServer(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Login, proto_base.Servee))
}

func FuzzConnLoginToClient(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Login, proto_base.Client))
}

func FuzzConnConfigToServer(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Configuration, proto_base.Servee))
}

func FuzzConnConfigToClient(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Configuration, proto_base.Client))
}

func FuzzConnPlayToServer(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Play, proto_base.Servee))
}

func FuzzConnPlayToClient(f *testing.F) {
	f.Fuzz(NewStateFuzzer(proto_base.Play, proto_base.Client))
}

func TestEntityMetadata(t *testing.T) {
	packetData := []byte{97, 198, 9, 9, 3, 65, 160, 0, 0, 16, 0, 127, 255}
	packet, err := proto.SimpleBytesToPacket(packetData, data.MustLookupProtocolVersion("1.21.11"), proto_base.ToClient, proto_base.Play)
	if err != nil {
		t.Fatal(err)
	}
	packetData2, err := proto.SimplePacketToBytes(packet)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(packetData, packetData2) {
		t.Error("packet data is not equal", cmp.Diff(packetData, packetData2))
	}
}
