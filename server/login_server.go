package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/google/uuid"
)

var StatusServedErr = errors.New("status served")

var UnexpectedStatusRequest = UnexpectedNextStateError{proto_base.Status}

type NameAndUUID struct {
	Name string
	UUID uuid.UUID
}

type UnexpectedNextStateError struct {
	NextState proto_base.State
}

func (u UnexpectedNextStateError) Error() string {
	return fmt.Sprintf("unexpected next state: %v", u.NextState)
}

type LoginServer struct {
	*proto.Conn
	//OnlineMode           bool FIXME: implement this
	ServerHost           string
	ServerPort           uint16
	Requested            NameAndUUID
	Given                *NameAndUUID
	CompressionThreshold int32 //FIXME: implement this
	Status               []byte
	CompatibleVersions   []int32
}

func (l *LoginServer) OnHandshake(packet *v1_21_8.HandshakingToServerPacketSetProtocol) (err error) {
	l.SetState(proto_base.State(packet.NextState))
	l.ServerHost = packet.ServerHost
	l.ServerPort = packet.ServerPort

	if packet.NextState == int32(proto_base.Status) && l.Status != nil {
		return
	}

	if packet.NextState != int32(proto_base.Login) {
		err = UnexpectedNextStateError{proto_base.State(packet.NextState)}
		return
	}

	if len(l.CompatibleVersions) != 0 && !slices.Contains(l.CompatibleVersions, packet.ProtocolVersion) {
		err = fmt.Errorf("incompatible protocol version please use one of %v you can look up the real corresponding version at https://minecraft.wiki/w/Minecraft_Wiki:Projects/wiki.vg_merge/Protocol_version_numbers", l.CompatibleVersions)
		return
	}

	return
}

func (l *LoginServer) OnStatusRequest(_ *v1_21_8.StatusToServerPacketPingStart) (err error) {
	err = l.Send(&v1_21_8.StatusToClientPacketServerInfo{Response: string(l.Status)})
	return
}

func (l *LoginServer) OnStatusPing(packet *v1_21_8.StatusToServerPacketPing) (err error) {
	err = l.Send(&v1_21_8.StatusToClientPacketPing{Time: packet.Time})
	if err != nil {
		return
	}
	return event.HandlerDoneErr{Return: StatusServedErr}
}

func (l *LoginServer) SetCompressionThreshold(threshold int32) (err error) {
	err = l.Send(&v1_21_8.LoginToClientPacketCompress{Threshold: threshold})
	if err != nil {
		return
	}
	l.CompressionThreshold = threshold
	return
}

func (l *LoginServer) OnLoginStart(packet *v1_21_8.LoginToServerPacketLoginStart) (err error) {
	l.Requested = NameAndUUID{packet.Username, packet.PlayerUUID}
	if l.Given == nil {
		l.Given = &l.Requested
	}
	err = l.Send(&v1_21_8.LoginToClientPacketSuccess{
		Uuid:       l.Given.UUID,
		Username:   l.Given.Name,
		Properties: nil,
	})
	return
}

func (l *LoginServer) OnLoginAcknowledged(_ *v1_21_8.LoginToServerPacketLoginAcknowledged) (err error) {
	l.SetState(proto_base.Configuration)
	err = event.HandlerDoneErr{}
	return
}

func (l *LoginServer) OnDefault(event event.Unhandled) (err error) {
	err = fmt.Errorf("unexpected event during login: %#v", event)
	return
}

func (l *LoginServer) OnStart() (err error) {
	return
}

func (l *LoginServer) OnCycle(_ event.Tick) (err error) {
	return
}

func (l *LoginServer) OnClose(_ event.Close) (err error) {
	return
}

type LoginServerSetting func(*LoginServer)

func WithOtherAccount(a *api.Account) LoginServerSetting {
	return func(s *LoginServer) {
		if s.Given == nil {
			return
		}
		s.Given = &NameAndUUID{
			Name: a.Name,
			UUID: a.Uuid,
		}
	}
}

func WithRawStatus(status []byte) LoginServerSetting {
	return func(s *LoginServer) {
		s.Status = status
	}
}

func WithStatus(response StatusResponse) LoginServerSetting {
	var status []byte
	status, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	slog.Debug("With status", "status", status)
	return WithRawStatus(status)
}

func WithCompatibleVersionsRange(start, end int32) LoginServerSetting {
	var versions []int32
	for i := start; i <= end; i++ {
		versions = append(versions, i)
	}
	return WithCompatibleVersions(versions...)
}

func WithCompatibleVersions(versions ...int32) LoginServerSetting {
	return func(s *LoginServer) {
		s.CompatibleVersions = append(versions)
	}
}

func ServeLogin(c *proto.Conn, settings ...LoginServerSetting) (ret *LoginServer, err error) {
	ret = &LoginServer{Conn: c}
	for _, s := range settings {
		s(ret)
	}
	ret.CompressionThreshold = 256
	err = ret.Start(ret)
	return
}
