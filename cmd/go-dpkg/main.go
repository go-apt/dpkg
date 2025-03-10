package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-apt/dpkg"
)

func main() {
	// Define flags
	infoFlag := flag.Bool("I", false, "show information about a package")
	listFlag := flag.Bool("l", false, "list packages matching given pattern")
	helpFlag := flag.Bool("?", false, "show this help message")

	// Parse flags
	flag.Parse()

	// Check if help flag is activated
	if *helpFlag {
		printUsage()
		return
	}

	// Check if some flag was requested
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Create a new instance of the Dpkg struct
	d := dpkg.NewDpkg()

	// Check if info flag is activated
	if *infoFlag {
		// Check if a .deb file is provided as an argument
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}
		debFile := os.Args[2]

		// Validate if the file is a .deb package
		if !d.IsDebFile(debFile) {
			fmt.Fprintf(os.Stderr, "Error: %s is not a valid .deb file\n", debFile)
			os.Exit(1)
		}

		// Read the contents of the .deb file
		pkg, err := d.Info(debFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Print package name from the .deb file
		fmt.Printf("Package from .deb file: %s\n", pkg)
	}

	// Check if list flag is activated
	if *listFlag {
		if len(os.Args) == 2 {
			// Read the contents of the /var/lib/dpkg/status file
			packages, err := d.List()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Print package name from installed packages
			for _, p := range packages {
				fmt.Printf("Package from dpkg status file: %s\n", p)
			}
		} else if len(os.Args) == 3 {
			patternName := os.Args[2]

			// List packages matching the given pattern
			filteredPackages, err := d.ListGrep(patternName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Print package name from installed packages
			for _, p := range filteredPackages {
				fmt.Printf("Packages from Grep search: %s\n", p)
			}
		}
	}
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Usage: go-dpkg [<option>...]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}
