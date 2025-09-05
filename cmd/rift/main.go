package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/server"
	"github.com/DarsenOP/Rift/internal/storage"

	"github.com/DarsenOP/Rift/pkg/version"
)

var humanMode bool

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	humanFlag := flag.Bool("human", false, "Switch to the human format RESP")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Rift version %s\n", version.Version)
		os.Exit(0)
	}

	if *humanFlag {
		humanMode = true
	} else {
		humanMode = false
	}

	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	gracefulLn := server.NewGracefulListener(listener)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println(">>> Shutdown signal received, draining connections")
		_ = gracefulLn.Listener.Close()
	}()

	fmt.Printf("Rift server v%s listening on :6380\n", version.Version)
	fmt.Println("Ready to accept connections...")

	for {
		conn, err := gracefulLn.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}

			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gracefulLn.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown error: %v", err)
	}
	log.Println(">>> Server exited")
}

func handleConnection(conn net.Conn) {
	store := storage.New()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	fmt.Printf("New connection from: %s\n", conn.RemoteAddr())

	for {
		var value resp.Value
		var err error
		if humanMode {
			value, err = resp.ParseHuman(reader)
		} else {
			value, err = resp.ParseRESP(reader)
		}
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("Connection closed by: %s\n", conn.RemoteAddr())
				break
			}

			log.Printf("Parse error from %s: %v", conn.RemoteAddr(), err)
			if writeErr := resp.WriteError(conn, err.Error()); writeErr != nil {
				log.Printf("Failed to send error response: %v", writeErr)
			}
			break
		}

		response := server.HandleCommand(store, value)

		err = resp.WriteValue(conn, response)
		if err != nil {
			log.Printf("Write error to %s: %v", conn.RemoteAddr(), err)
			break
		}
	}
}
