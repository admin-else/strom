package server

import (
	"encoding/json"

	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/text"
)

// These are separate types so a status definition looks less ugly

type StatusResponseVersion struct {
	Name     *text.Component `json:"name"`
	Protocol int32           `json:"protocol"`
}

type StatusResponseSample struct {
	Name *text.Component `json:"name"`
	ID   string          `json:"id"`
}
type StatusResponsePlayers struct {
	Max    int                    `json:"max"`
	Online int                    `json:"online"`
	Sample []StatusResponseSample `json:"sample"`
}

// StatusResponse here is the structure I used: https://minecraft.wiki/w/Java_Edition_protocol/Server_List_Ping#Status_Response
type StatusResponse struct {
	Version            StatusResponseVersion `json:"version"`
	Players            StatusResponsePlayers `json:"players"`
	Description        *text.Component       `json:"description"`
	Favicon            string                `json:"favicon"`
	EnforcesSecureChat bool                  `json:"enforcesSecureChat"`
}

// MustMarshal is a utility for marshalling StatusResponse to JSON
func (s StatusResponse) MustMarshal() []byte {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return b
}

type StatusServer struct {
	*proto.Conn
	Status string
}

func (p *StatusServer) OnStatusRequest(_ *v1_21_8.StatusToServerPacketPingStart) (err error) {
	return p.Send(&v1_21_8.StatusToClientPacketServerInfo{Response: p.Status})
}

func (p *StatusServer) OnStatusPing(packet *v1_21_8.StatusToServerPacketPing) (err error) {
	return p.Send(&v1_21_8.StatusToClientPacketPing{Time: packet.Time})
}

func ServeStatus(c *proto.Conn, s StatusResponse) (err error) {
	status, err := json.Marshal(s)
	if err != nil {
		return
	}
	server := &StatusServer{
		Conn:   c,
		Status: string(status),
	}
	server.RegisterUntilLatest(server.OnStatusRequest)
	server.RegisterUntilLatest(server.OnStatusPing)
	err = server.StartLoop()
	return
}
