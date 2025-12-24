package proto_test

import (
	"bytes"
	"math"

	"testing"

	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/google/go-cmp/cmp"
)

var nanIgnorer32 = cmp.Comparer(func(x, y float32) bool {
	return (math.IsNaN(float64(x)) && math.IsNaN(float64(y))) || x == y
})

var nanIgnorer64 = cmp.Comparer(func(x, y float64) bool {
	return (math.IsNaN(x) && math.IsNaN(y)) || x == y
})

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
