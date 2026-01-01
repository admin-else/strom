package proto_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/server"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var nanIgnorer32 = cmp.Comparer(func(x, y float32) bool {
	return (math.IsNaN(float64(x)) && math.IsNaN(float64(y))) || x == y
})

var nanIgnorer64 = cmp.Comparer(func(x, y float64) bool {
	return (math.IsNaN(x) && math.IsNaN(y)) || x == y
})

func FuzzConn(f *testing.F) {
	f.Fuzz(func(t *testing.T, inbytes []byte) {
		b1 := bytes.NewBuffer(inbytes)
		b2 := bytes.NewBuffer(nil)
		conn := server.Servee(nil)
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
		return
	})
}
