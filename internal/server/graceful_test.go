package server

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

func TestGracefulListener_Shutdown(t *testing.T) {
	// Create a test listener
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer func() { _ = ln.Close() }()

	gracefulLn := NewGracefulListener(ln)

	// Test shutdown without active connections
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = gracefulLn.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown without connections failed: %v", err)
	}
}

func TestGracefulListener_ShutdownWithActiveConnections(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer func() { _ = ln.Close() }()

	gracefulLn := NewGracefulListener(ln)

	// Simulate active connection
	// conn := &mockConn{}
	gracefulLn.wg.Add(1)

	// Start shutdown in background
	shutdownDone := make(chan error)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		shutdownDone <- gracefulLn.Shutdown(ctx)
	}()

	// Wait a bit then "close" the connection
	time.Sleep(100 * time.Millisecond)
	gracefulLn.wg.Done()

	// Check shutdown completed successfully
	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Errorf("Shutdown with active connections failed: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Shutdown timed out")
	}
}

func TestGracefulListener_ShutdownTimeout(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer func() { _ = ln.Close() }()

	gracefulLn := NewGracefulListener(ln)

	// Add connection that never closes
	gracefulLn.wg.Add(1)

	// Try shutdown with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = gracefulLn.Shutdown(ctx)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded, got %v", err)
	}
}

func TestGracefulConn_Close(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	conn := &gracefulConn{
		Conn: &mockConn{},
		wg:   wg,
	}

	// First close should decrement waitgroup
	err := conn.Close()
	if err != nil {
		t.Errorf("First close failed: %v", err)
	}

	// Second close should be no-op
	err = conn.Close()
	if err != nil {
		t.Errorf("Second close failed: %v", err)
	}
}

// mockConn implements net.Conn for testing
type mockConn struct{}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return &mockAddr{} }
func (m *mockConn) RemoteAddr() net.Addr               { return &mockAddr{} }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// mockAddr implements net.Addr for testing
type mockAddr struct{}

func (m *mockAddr) Network() string { return "tcp" }
func (m *mockAddr) String() string  { return "localhost:0" }
