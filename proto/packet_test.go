package proto_test

import (
	"bytes"
	"math"

	"testing"

	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/google/uuid"

	"github.com/google/go-cmp/cmp"
)

var nanIgnorer32 = cmp.Comparer(func(x, y float32) bool {
	return (math.IsNaN(float64(x)) && math.IsNaN(float64(y))) || x == y
})

var nanIgnorer64 = cmp.Comparer(func(x, y float64) bool {
	return (math.IsNaN(x) && math.IsNaN(y)) || x == y
})

func TestChat(t *testing.T) {
	//p := v1_21_8.PlayToServerPacketChatMessage{Message: "dfs", Timestamp: 1764153707364, Salt: -5634100397920314628, Signature: (*[256]uint8)(nil), Offset: 0, Acknowledged: [3]uint8{0x0, 0x0, 0x0}, Checksum: 0x1}
	p := v1_21_8.PlayToClientPacketPlayerChat{GlobalIndex: 0, SenderUuid: uuid.UUID{0x1d, 0xc5, 0x57, 0x82, 0x3c, 0x0, 0x3e, 0xd4, 0xbf, 0x8b, 0x76, 0xd4, 0xc6, 0xd9, 0xcc, 0x47}, Index: 0, Signature: (*[256]uint8)(nil), PlainMessage: "csdf", Timestamp: 1764153823132, Salt: 0, PreviousMessages: v1_21_8.PreviousMessages{Val: []struct {
		Id        int32
		Signature interface{}
	}{}}, UnsignedChatContent: (*nbt.Anon)(nil), FilterType: 0, FilterTypeMask: proto_base.Void{}, Type: v1_21_8.PlayToClientChatTypesHolder{Val: 1}, NetworkName: nbt.Anon{Value: map[string]interface{}{"click_event": map[string]interface{}{"action": "suggest_command", "command": "/tell Player659 "}, "hover_event": map[string]interface{}{"action": "show_entity", "id": "minecraft:player", "name": "Player659", "uuid": []int32{499472258, 1006649044, -1081379116, -958804921}}, "insertion": "Player659", "text": "Player659"}}, NetworkTargetName: (*nbt.Anon)(nil)}
	b := bytes.NewBuffer(nil)
	err := p.Encode(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%q", b.Bytes())
}

func FuzzPlayPackets(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		t.Logf("Fuzz data %q\n", b)
		p := v1_21_8.PlayToClientPacket{}
		b1 := bytes.NewBuffer(b)
		p, err := p.Decode(b1)
		if err != nil {
			t.Skipf("failed to decode %q because %v (partial %v)", b, err, p)
			return
		}
		if b1.Len() != 0 {
			t.Skipf("there are %d bytes left it should be read entirely", b1.Len())
			return
		}
		t.Logf("decoded %#v\n", p)
		b2 := bytes.NewBuffer(nil)
		err = p.Encode(b2)
		if err != nil {
			t.Errorf("failed to encode %q because %v", b, err)
			return
		}
		t.Logf("encoded %q\n", b2.Bytes())
		p2 := v1_21_8.PlayToClientPacket{}
		p2, err = p2.Decode(b2)
		if err != nil {
			t.Errorf("failed to decode encoded data: %v (partial %#v)", err, p2)
			return
		}
		if !cmp.Equal(p, p2, nanIgnorer32, nanIgnorer64) {
			t.Errorf("decoded and encoded decoded data are not equal: %#v", p2)
		}
	})
}
