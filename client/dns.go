package client

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
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

type DNSCheckFunc func(addr string) error

// DoDns resolves a DNS name to an IP address
// It should resolve a minecraft compatible address into something net.Dial
func DoDns(in string) (out string, err error) {
	return DoDNSChecked(in, nil)
}

// DoDNSChecked resolves a DNS name to an IP address
// It should resolve a minecraft compatible address into something net.Dial
// It also optionally accepts a function to "check" the address
func DoDNSChecked(in string, f DNSCheckFunc) (out string, err error) { // Bad code should match minecraft
	// https://mcsrc.dev/1/26.1-pre-1/net/minecraft/client/multiplayer/resolver/ServerAddress#L34
	host, port, err := net.SplitHostPort(in)
	if err != nil {
		var addrErr *net.AddrError
		if errors.As(err, &addrErr) && addrErr.Err == "missing port in address" {
			err = nil
			host = in
			port = "25565"
		} else {
			return
		}
	}
	porti, err := strconv.Atoi(port)
	if err != nil {
		return
	}
	if f == nil {
		f = func(addr string) error { return nil }
	}

	err = f(host)
	if err != nil {
		return
	}
	portN := uint16(porti)
	outIP, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return
	}
	err = f(outIP.String())
	if err != nil {
		return
	}
	out = outIP.String() + ":" + port
	if portN != 25565 {

		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) // otherwise this takes 10+ seconds
	defer cancel()

	r := &net.Resolver{}

	_, srvs, _ := r.LookupSRV(ctx, "minecraft", "tcp", host)
	if len(srvs) == 0 {

		return
	}
	host = strings.TrimSuffix(srvs[0].Target, ".")
	portN = srvs[0].Port
	outIP, err = net.ResolveIPAddr("ip", host)
	if err != nil {
		return
	}
	err = f(outIP.String())
	if err != nil {
		return
	}
	out = outIP.String() + ":" + port
	return
}
