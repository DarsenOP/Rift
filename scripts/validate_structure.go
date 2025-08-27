// Script to validate project structure
package main

import (
	"fmt"
	"os"
)

func main() {
	requiredDirs := []string{
		"cmd/rift",
		"internal",
		"pkg",
		"scripts",
	}
	
	requiredFiles := []string{
		"cmd/rift/main.go",
		"go.mod",
		"go.sum", 
		"Makefile",
		".golangci.yml",
		".github/workflows/go.yml",
	}
	
	allGood := true
	
	// Check directories
	for _, dir := range requiredDirs {
		if info, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("xxx Missing directory: %s\n", dir)
			allGood = false
		} else if !info.IsDir() {
			fmt.Printf("xxx Not a directory: %s\n", dir)	
			allGood = false
		}
	}
	
	// Check files
	for _, file := range requiredFiles {
		if info, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("xxx Missing file: %s\n", file)
			allGood = false
		} else if info.IsDir() {
			fmt.Printf("xxx Not a file: %s\n", file)
			allGood = false
		}
	}
	
	if allGood {
		fmt.Println("--- Project structure is valid!")
		os.Exit(0)
	} else {
		fmt.Println("xxx Project structure validation failed!")
		os.Exit(1)
	}
}
