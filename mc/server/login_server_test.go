package server_test

import (
	"reflect"
	"testing"

	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
)

func assertSendConvertible(t *testing.T, srcType reflect.Type, name string, state proto_base.State, direction proto_base.Direction, targetVersion int32) {
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

func assertReceiveConvertible(t *testing.T, dstType reflect.Type, name string, state proto_base.State, direction proto_base.Direction, targetVersion int32) {
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

func TestLoginServerSendAssumptions(t *testing.T) {
	v26_1Version := data.MustLookupProtocolVersion("26.1")

	// Packets sent by LoginServer as v1_21_8 must be convertible to every version
	// they are actually sent to. login_success is excluded because it is not
	// convertible across the 26.2 boundary; it has dedicated version-gated send code.
	assertSendConvertible(t, reflect.TypeFor[*v1_21_8.StatusToClientPacketServerInfo](), "server_info", proto_base.Status, proto_base.ToClient, v26_1Version)
	assertSendConvertible(t, reflect.TypeFor[*v1_21_8.StatusToClientPacketPing](), "ping", proto_base.Status, proto_base.ToClient, v26_1Version)
	assertSendConvertible(t, reflect.TypeFor[*v1_21_8.LoginToClientPacketCompress](), "compress", proto_base.Login, proto_base.ToClient, v26_1Version)
	assertSendConvertible(t, reflect.TypeFor[*v1_21_8.LoginToClientPacketEncryptionBegin](), "encryption_begin", proto_base.Login, proto_base.ToClient, v26_1Version)
}

func TestLoginServerReceiveAssumptions(t *testing.T) {
	v26_1Version := data.MustLookupProtocolVersion("26.1")

	// Packets received by LoginServer via RegisterCriticalUntilLatest must be
	// convertible from each later version back to v1_21_8.
	assertReceiveConvertible(t, reflect.TypeFor[*v1_21_8.HandshakingToServerPacketSetProtocol](), "set_protocol", proto_base.Handshaking, proto_base.ToServer, v26_1Version)
	assertReceiveConvertible(t, reflect.TypeFor[*v1_21_8.StatusToServerPacketPingStart](), "ping_start", proto_base.Status, proto_base.ToServer, v26_1Version)
	assertReceiveConvertible(t, reflect.TypeFor[*v1_21_8.StatusToServerPacketPing](), "ping", proto_base.Status, proto_base.ToServer, v26_1Version)
	assertReceiveConvertible(t, reflect.TypeFor[*v1_21_8.LoginToServerPacketLoginStart](), "login_start", proto_base.Login, proto_base.ToServer, v26_1Version)
	assertReceiveConvertible(t, reflect.TypeFor[*v1_21_8.LoginToServerPacketEncryptionBegin](), "encryption_begin", proto_base.Login, proto_base.ToServer, v26_1Version)
	assertReceiveConvertible(t, reflect.TypeFor[*v1_21_8.LoginToServerPacketLoginAcknowledged](), "login_acknowledged", proto_base.Login, proto_base.ToServer, v26_1Version)
}
