package server

import (
	"context"
	"net"
	"sync"
)

// GracefulListener wraps a net.Listener and counts active connections.
type GracefulListener struct {
	net.Listener
	wg        sync.WaitGroup
	closeOnce sync.Once
}

// NewGracefulListener returns the wrapped listener ready for use.
func NewGracefulListener(ln net.Listener) *GracefulListener {
	return &GracefulListener{Listener: ln}
}

// Accept waits for and returns the next connection to the listener.
// It increments an internal wait-group before returning.
func (g *GracefulListener) Accept() (net.Conn, error) {
	conn, err := g.Listener.Accept()
	if err != nil {
		return nil, err
	}
	g.wg.Add(1)
	return &gracefulConn{Conn: conn, wg: &g.wg}, nil
}

// Shutdown gracefully closes the underlying listener and waits
// for all tracked connections to finish.
func (g *GracefulListener) Shutdown(ctx context.Context) error {
	// stop accepting new connections
	g.closeOnce.Do(func() { _ = g.Listener.Close() })

	// wait for active handlers or until context cancelled
	done := make(chan struct{})
	go func() { g.wg.Wait(); close(done) }()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// gracefulConn decrements the wait-group on Close.
type gracefulConn struct {
	net.Conn
	wg   *sync.WaitGroup
	once sync.Once
}

func (c *gracefulConn) Close() error {
	c.once.Do(c.wg.Done)
	return c.Conn.Close()
}
