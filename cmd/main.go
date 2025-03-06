package main

import (
	"fmt"
	"os"

	"github.com/go-apt/dpkg"
)

func main() {
	// Check if a .deb file is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <arquivo.deb>")
		return
	}

	debFile := os.Args[1]

	// Create a new instance of the Dpkg struct
	d := dpkg.Dpkg{}

	// Read the contents of the .deb file
	pkg, err := d.Info(debFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print package name from the .deb file
	fmt.Printf("Package from .deb file: %s\n", pkg.Package)

	// Read the contents of the /var/lib/dpkg/status file
	packages, err := d.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print package name from installed packages
	for _, p := range packages {
		fmt.Printf("Package from dpkg status file: %s\n", p.Package)
	}
}
