package client

import (
	"net"
	"strings"
)

func LookupMinecraftSRV(host string, port uint16) (hostRet string, portRet uint16) {
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
