package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/server"
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
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

	fmt.Printf("Rift server v%s listening on :6380\n", version.Version)
	fmt.Println("Ready to accept connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
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

		response := server.HandleCommand(value)

		err = resp.WriteValue(conn, response)
		if err != nil {
			log.Printf("Write error to %s: %v", conn.RemoteAddr(), err)
			break
		}
	}
}
