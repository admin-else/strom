package client_test

import (
	"reflect"
	"testing"

	"git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v1_21_8"
)

func assertClientSendConvertible(t *testing.T, srcType reflect.Type, name string, state proto_base.State, direction proto_base.Direction, targetVersion int32) {
	t.Helper()
	targetInfo, ok := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(name, targetVersion, state, direction)
	if !ok {
		t.Fatalf("packet %s not found for version %d", name, targetVersion)
	}
	targetType := reflect.TypeOf(targetInfo.Type)
	if !proto.SmartConvertibleTo(srcType, targetType) {
		t.Fatalf("%s is not smart-convertible from v1_21_8 to version %d", name, targetVersion)
	}
}

func assertClientReceiveConvertible(t *testing.T, dstType reflect.Type, name string, state proto_base.State, direction proto_base.Direction, targetVersion int32) {
	t.Helper()
	targetInfo, ok := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(name, targetVersion, state, direction)
	if !ok {
		t.Fatalf("packet %s not found for version %d", name, targetVersion)
	}
	targetType := reflect.TypeOf(targetInfo.Type)
	if !proto.SmartConvertibleTo(targetType, dstType) {
		t.Fatalf("%s is not smart-convertible from version %d to v1_21_8", name, targetVersion)
	}
}

func TestLoginClientSendAssumptions(t *testing.T) {
	v26_1Version := data.MustLookupProtocolVersion("26.1")

	// Packets sent by LoginClient as v1_21_8.
	assertClientSendConvertible(t, reflect.TypeFor[*v1_21_8.HandshakingToServerPacketSetProtocol](), "set_protocol", proto_base.Handshaking, proto_base.ToServer, v26_1Version)
	assertClientSendConvertible(t, reflect.TypeFor[*v1_21_8.LoginToServerPacketLoginStart](), "login_start", proto_base.Login, proto_base.ToServer, v26_1Version)
	assertClientSendConvertible(t, reflect.TypeFor[*v1_21_8.LoginToServerPacketLoginAcknowledged](), "login_acknowledged", proto_base.Login, proto_base.ToServer, v26_1Version)
	assertClientSendConvertible(t, reflect.TypeFor[*v1_21_8.LoginToServerPacketEncryptionBegin](), "encryption_begin", proto_base.Login, proto_base.ToServer, v26_1Version)
}

func TestLoginClientReceiveAssumptions(t *testing.T) {
	v26_1Version := data.MustLookupProtocolVersion("26.1")

	// Packets received by LoginClient via RegisterCriticalUntilLatest.
	// login_success is excluded because 26.2 is not convertible.
	assertClientReceiveConvertible(t, reflect.TypeFor[*v1_21_8.LoginToClientPacketCompress](), "compress", proto_base.Login, proto_base.ToClient, v26_1Version)
	assertClientReceiveConvertible(t, reflect.TypeFor[*v1_21_8.LoginToClientPacketEncryptionBegin](), "encryption_begin", proto_base.Login, proto_base.ToClient, v26_1Version)
	assertClientReceiveConvertible(t, reflect.TypeFor[*v1_21_8.LoginToClientPacketDisconnect](), "disconnect", proto_base.Login, proto_base.ToClient, v26_1Version)
}
