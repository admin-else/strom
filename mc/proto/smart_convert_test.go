package proto_test

import (
	"reflect"
	"testing"

	proto2 "github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_9"
	"github.com/admin-else/strom/mc/proto_generated/v26_2"
)

func TestCanSmartConvert(t *testing.T) {
	packet := &v1_21_9.PlayToClientPacketPosition{Flags: v1_21_9.PlayToClientPositionUpdateRelatives{Val: 2}, X: 5}
	can := proto2.SmartConvertibleTo(reflect.TypeOf(packet), reflect.TypeFor[*v1_21_8.PlayToClientPacketPosition]())
	if !can {
		t.Fatal("can't smart convert")
	}
	ret := proto2.SmartConvert(reflect.ValueOf(packet), reflect.TypeFor[*v1_21_8.PlayToClientPacketPosition]())
	t.Logf("%#v", ret)
}

func TestCanSmartConvertTags(t *testing.T) {
	packet := &v1_21_9.ConfigurationToClientPacketTags{}
	targetType := reflect.TypeFor[*v1_21_8.ConfigurationToClientPacketTags]()
	can := proto2.SmartConvertibleTo(reflect.TypeOf(packet), targetType)
	if !can {
		t.Fatal("can't smart convert")
	}
	ret := proto2.SmartConvert(reflect.ValueOf(packet), targetType)
	t.Logf("%#v", ret)
}

func TestRegisterUntilLatest(t *testing.T) {
	conn := proto2.NewConn()
	ok := false
	conn.RegisterCriticalUntilLatest(func(p *v1_21_8.ConfigurationToServerPacketFinishConfiguration) error {
		ok = true
		return nil
	})
	err := conn.Fire(&v1_21_11.ConfigurationToServerPacketFinishConfiguration{})
	if err != nil {
		t.Errorf("fire failed err: %v", err)
	}
	if !ok {
		t.Errorf("handler not called")
	}
}

var TestData = []struct{ V, T, E any }{
	{V: &v1_21_9.PlayToClientPacketPosition{Flags: v1_21_9.PlayToClientPositionUpdateRelatives{Val: 2}, X: 5}, T: &v1_21_8.PlayToClientPacketPosition{}},
	{V: &v1_21_9.ConfigurationToClientPacketTags{}, T: &v1_21_8.ConfigurationToClientPacketTags{}},
	{V: &v1_21_11.PlayToServerPacketChatMessage{Signature: &[256]byte{}}, T: &v1_21_8.PlayToServerPacketChatMessage{}},
}

func TestSmartConvert(t *testing.T) {
	for _, d := range TestData {
		can := proto2.SmartConvertibleTo(reflect.TypeOf(d.V), reflect.TypeOf(d.T))
		if !can {
			t.Errorf("can't smart convert %#v to %#v", d.V, d.T)
			continue
		}
		got := proto2.SmartConvert(reflect.ValueOf(d.V), reflect.TypeOf(d.T))
		if d.E != nil && !reflect.DeepEqual(got.Interface(), d.E) {
			t.Errorf("got %#v, want %#v", got.Interface(), d.E)
		}
	}
}

func TestLoginSuccessNotConvertibleAcross26_2(t *testing.T) {
	// login_success gained an extra SessionId field in 26.2, so it is not
	// smart-convertible in either direction. This is why the login server has
	// version-gated send code and the login client has a dedicated 26.2 handler.
	if proto2.SmartConvertibleTo(reflect.TypeFor[*v1_21_8.LoginToClientPacketSuccess](), reflect.TypeFor[*v26_2.LoginToClientPacketSuccess]()) {
		t.Fatal("expected v1_21_8 login_success NOT to be convertible to v26_2")
	}
	if proto2.SmartConvertibleTo(reflect.TypeFor[*v26_2.LoginToClientPacketSuccess](), reflect.TypeFor[*v1_21_8.LoginToClientPacketSuccess]()) {
		t.Fatal("expected v26_2 login_success NOT to be convertible to v1_21_8")
	}
}
