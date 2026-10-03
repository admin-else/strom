package client

import (
	"testing"
	"time"
)

func TestPingMillis(t *testing.T) {
	s := &StatusClient{}
	if got := s.PingMillis(); got != 0 {
		t.Errorf("ping before exchange = %d", got)
	}
	s.PingSendTime = time.UnixMilli(1000)
	s.PingReceiveTime = time.UnixMilli(1042)
	if got := s.PingMillis(); got != 42 {
		t.Errorf("ping = %d", got)
	}
}
