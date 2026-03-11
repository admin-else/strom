package client

import (
	"net"
	"strconv"
	"strings"
)

func LookupMinecraftSRV(host string, port uint16) (hostRet string, portRet uint16) {
	hostRet = host // in case of error
	portRet = port
	if port != 25565 {
		return
	}
	_, srvs, err := net.LookupSRV("minecraft", "tcp", host)
	if err != nil || len(srvs) == 0 {
		return
	}
	target := strings.TrimSuffix(srvs[0].Target, ".")
	hostRet = target
	portRet = srvs[0].Port
	return
}

func DoDNS(in string) (out string, err error) { // Bad code should match minecraft
	// https://mcsrc.dev/1/26.1-pre-1/net/minecraft/client/multiplayer/resolver/ServerAddress#L34
	host, port, _ := net.SplitHostPort(in)
	var portN uint16 = 25565
	if port != "" {
		var porti int
		porti, err = strconv.Atoi(port)
		if err != nil {
			return
		}
		portN = uint16(porti)
	}
	host, portN = LookupMinecraftSRV(host, portN)
	out = net.JoinHostPort(host, strconv.Itoa(int(portN)))
	return
}
