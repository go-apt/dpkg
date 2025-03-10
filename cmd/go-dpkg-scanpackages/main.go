package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-apt/dpkg"
)

func main() {
	// Define flags
	typeFlag := flag.String("t", "deb", "scan for <type> packages (default is 'deb')")
	archFlag := flag.String("a", "", "architecture to scan for")
	hashFlag := flag.String("h", "md5,sha1,sha256", "only generate hashes for the specified comma separated list")
	multiversionFlag := flag.Bool("m", false, "allow multiple versions of a single package")
	helpFlag := flag.Bool("?", false, "show this help message")
	versionFlag := flag.Bool("version", false, "show the version")

	// Parse flags
	flag.Parse()

	// Check if help flag is activated
	if *helpFlag {
		printUsage()
		return
	}

	// Check if version flag is activated
	if *versionFlag {
		printVersion()
		return
	}

	// Check if binary path is provided
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: binary-path is required")
		printUsage()
		os.Exit(1)
	}

	binaryPath := args[0]
	// pathPrefix := ""

	// if len(args) > 1 {
	// 	pathPrefix = args[1]
	// }

	// Create a new PackagesScanner instance
	sp := dpkg.NewPackagesScanner(binaryPath)
	sp.Type = *typeFlag
	sp.Arch = *archFlag
	sp.Hashes = strings.Split(*hashFlag, ",")
	sp.Multiversion = *multiversionFlag

	// Scan packages
	packages, err := sp.ScanPackages()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print the content of the packages variable to stdout using a buffer
	buffer := bytes.NewBuffer(packages)
	if _, err := buffer.WriteTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
		os.Exit(1)
	}
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Usage: go-dpkg-scanpackages [<option>...] <binary-path> [<path-prefix>] > Packages")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}

// printVersion prints the version information
func printVersion() {
	fmt.Println("go-dpkg-scanpackages version 1.0")
}
