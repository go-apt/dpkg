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
		panic(err)
	}

	fmt.Printf("Package: %s\n", pkg.Package)
	fmt.Printf("Version: %s\n", pkg.Version)
	fmt.Printf("Architecture: %s\n", pkg.Architecture)
	fmt.Printf("Maintainer: %s\n", pkg.Maintainer)
	fmt.Printf("Description: %s\n", pkg.Description)
}
