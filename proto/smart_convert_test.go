package proto_test

import (
	"reflect"
	"testing"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_generated/v1_21_11"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/admin-else/strom/proto_generated/v1_21_9"
)

func TestCanSmartConvert(t *testing.T) {
	packet := &v1_21_9.PlayToClientPacketPosition{Flags: v1_21_9.PlayToClientPositionUpdateRelatives{Val: 2}, X: 5}
	can := proto.SmartConvertibleTo(reflect.TypeOf(packet), reflect.TypeFor[*v1_21_8.PlayToClientPacketPosition]())
	if !can {
		t.Fatal("can't smart convert")
	}
	ret := proto.SmartConvert(reflect.ValueOf(packet), reflect.TypeFor[*v1_21_8.PlayToClientPacketPosition]())
	t.Logf("%#v", ret)
}

func TestCanSmartConvertTags(t *testing.T) {
	packet := &v1_21_9.ConfigurationToClientPacketTags{}
	targetType := reflect.TypeFor[*v1_21_8.ConfigurationToClientPacketTags]()
	can := proto.SmartConvertibleTo(reflect.TypeOf(packet), targetType)
	if !can {
		t.Fatal("can't smart convert")
	}
	ret := proto.SmartConvert(reflect.ValueOf(packet), targetType)
	t.Logf("%#v", ret)
}

func TestRegisterUntilLatest(t *testing.T) {
	conn := proto.NewConn()
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
