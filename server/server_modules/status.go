package server_modules

import (
	"encoding/json"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/admin-else/strom/text"
)

// StatusResponse here is the structure I used: https://minecraft.wiki/w/Java_Edition_protocol/Server_List_Ping#Status_Response
type StatusResponse struct {
	Version struct {
		Name     *text.Component `json:"name"`
		Protocol int             `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name *text.Component `json:"name"`
			ID   string          `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description        *text.Component `json:"description"`
	Favicon            string          `json:"favicon"`
	EnforcesSecureChat bool            `json:"enforcesSecureChat"`
}

type StatusServer struct {
	*proto.Conn
	Status string
}

func (p *StatusServer) Default(_ any) (err error) {
	return
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
	err = server.Start(server)
	return
}
