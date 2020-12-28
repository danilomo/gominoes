package faketcp

import (
	"net"
)

// FakeServer aaa
func FakeServer(gen func() net.Conn) net.Listener {
	return &fakeTCP{gen}
}

type fakeTCP struct {
	generator func() net.Conn
}

func (fake *fakeTCP) Accept() (net.Conn, error) {
	return fake.generator(), nil
}

func (fake *fakeTCP) Close() error {
	return nil
}

func (fake *fakeTCP) Addr() net.Addr {
	return nil
}

// Server creates a new fake TCP server listener based on a generator function
func Server(conGenerator func() net.Conn) net.Listener {
	return &fakeTCP{conGenerator}
}
