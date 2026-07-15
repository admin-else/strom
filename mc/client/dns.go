package client

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

type DNSCheckFunc func(addr string) error

// DoDNSChecked resolves a DNS name to an IP address
// It should resolve a minecraft compatible address into something net.Dial
// It also optionally accepts a function to "check" the address
func DoDNSChecked(ctx context.Context, in string, f DNSCheckFunc) (connectTo string, host string, portN uint16, err error) {
	if f == nil {
		f = func(addr string) error { return nil }
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// https://mcsrc.dev/1/26.1-pre-1/net/minecraft/client/multiplayer/resolver/ServerAddress#L34
	host, port, err := net.SplitHostPort(in)
	if err != nil {
		if addrErr, ok := errors.AsType[*net.AddrError](err); ok && addrErr.Err == "missing port in address" {
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

	err = f(host)
	if err != nil {
		return
	}
	portN = uint16(porti)
	outIP, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return
	}
	err = f(outIP.String())
	if err != nil {
		return
	}
	connectTo = outIP.String() + ":" + port
	if portN != 25565 {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second) // otherwise this takes 10+ seconds
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
	connectTo = outIP.String() + ":" + port
	return
}
