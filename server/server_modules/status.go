package server_modules

import (
	"encoding/json"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

// StatusResponse here is the structure I used: https://minecraft.wiki/w/Java_Edition_protocol/Server_List_Ping#Status_Response
type StatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon            string `json:"favicon"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat"`
}

type StatusServer struct {
	*proto.Conn
	Status string
}

func (p *StatusServer) Default(_ any) (err error) {
	return
}

func (p *StatusServer) OnStatusRequest(_ v1_21_8.StatusToServerPacketPingStart) (err error) {
	return p.Send(v1_21_8.StatusToClientPacketServerInfo{Response: p.Status})
}

func (p *StatusServer) OnStatusPing(packet v1_21_8.StatusToServerPacketPing) (err error) {
	return p.Send(v1_21_8.StatusToClientPacketPing(packet))
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
