package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DarsenOP/Rift/pkg/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Rift version %s\n", version.Version)
		os.Exit(0)
	}
}
