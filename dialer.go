package torgo

import (
	"context"
	"net"
)

type dialer interface {
	// Dial connects to the given address via the proxy.
	Dial(network, addr string) (c net.Conn, err error)
}

type comboDialer interface {
	dialer
	ContextDialer
}

// A ContextDialer dials using a context.
type ContextDialer interface {
	// DialContext connects to the address on the named network using
	// the provided context.
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type comboDialAdapter func(ctx context.Context, network, address string) (net.Conn, error)

func (f comboDialAdapter) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return f(ctx, network, address)
}

func (f comboDialAdapter) Dial(network, address string) (net.Conn, error) {
	return f(context.Background(), network, address)
}
