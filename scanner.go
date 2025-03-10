package dpkg

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type HASH = string

const (
	MD5    HASH = "MD5"
	SHA1   HASH = "SHA1"
	SHA256 HASH = "SHA256"
)

var (
	AVAILABLE_HASHES []HASH = []HASH{MD5, SHA1, SHA256}
)

// PackagesScanner represents a scanner for Debian packages
type PackagesScanner struct {
	RootDir      string
	Arch         string
	Type         string
	Hashes       []HASH
	Multiversion bool
}

// NewPackagesScanner creates a new instance of PackagesScanner
func NewPackagesScanner(dir string) *PackagesScanner {
	return &PackagesScanner{
		RootDir:      dir,
		Arch:         "",
		Type:         "deb",
		Hashes:       []HASH{MD5, SHA1, SHA256},
		Multiversion: false,
	}
}

func isValidHash(h HASH) bool {
	for _, a := range AVAILABLE_HASHES {
		if strings.EqualFold(a, h) {
			return true
		}
	}
	return false
}

func containsHash(hashes []HASH, h HASH) bool {
	for _, a := range hashes {
		if strings.EqualFold(a, h) {
			return true
		}
	}
	return false
}

// ScanPackages scans the directory for packages matching the criteria
func (ps *PackagesScanner) ScanPackages() ([]byte, error) {
	// Check if the hashes are valid
	for _, h := range ps.Hashes {
		if !isValidHash(h) {
			return nil, fmt.Errorf("invalid hash: %s", h)
		}
	}

	// List the files in the directory that match the filter
	files, err := findFiles(ps.RootDir, ps.createFilter())
	if err != nil {
		return nil, fmt.Errorf("error finding files: %v", err)
	}

	pkgs := make(map[string]DebPackage)
	d := NewDpkg()

	for _, file := range files {
		p, err := d.Info(file)
		if err != nil {
			// Let the user know that there was an error parsing the package
			// but continue to the next package
			fmt.Fprintf(os.Stderr, "go-apt/dpkg: error parsing package: '%s' - %v\n", file, err)
			continue
		}

		if checkMultivalue(p, &pkgs) {
			op := pkgs[p.Fields["Package"]]

			if ps.Multiversion {
				fmt.Fprintf(os.Stderr, "multiversion enabled; adding repeated Package '%s'\n", p.Fields["Package"])
				pkgs[p.Fields["Package"]+p.Fields["Version"]] = *p
				continue
			}

			if d.CompareVersions(p.Fields["Version"], op.Fields["Version"]) > 0 {
				fmt.Fprintf(os.Stderr, "package '%s' (filename '%s') is repeat but newer version; ignored version '%s'!\n", p.Fields["Package"], p.Fields["Filename"], op.Fields["Version"])
				pkgs[p.Fields["Package"]] = *p
			} else {
				fmt.Fprintf(os.Stderr, "package '%s' (filename '%s') is repeat but older; ignored version '%s'!\n", p.Fields["Package"], p.Fields["Filename"], p.Fields["Version"])
				continue
			}
		}

		pkgs[p.Fields["Package"]] = *p
	}

	return generatePackageIndex(&pkgs, &ps.Hashes), nil
}

// createFilter creates a regular expression filter based on architecture and type
func (ps *PackagesScanner) createFilter() *regexp.Regexp {
	if ps.Arch != "" {
		// Create the regular expression with the architecture
		pattern := fmt.Sprintf("_(?:all|%s)\\.%s$", ps.Arch, ps.Type)
		return regexp.MustCompile(pattern)
	} else {
		// Create the regular expression without the architecture
		pattern := fmt.Sprintf("\\.%s$", ps.Type)
		return regexp.MustCompile(pattern)
	}
}

// findFiles finds files that match the filter, including subdirectories
func findFiles(directory string, filter *regexp.Regexp) ([]string, error) {
	var matches []string

	// Use filepath.WalkDir to traverse the directory and subdirectories
	err := filepath.WalkDir(directory, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err // Return errors encountered during the walk
		}

		// Check if the entry is a file (ignore directories)
		if !entry.IsDir() {
			fileName := entry.Name()
			// Check if the file name matches the filter
			if filter.MatchString(fileName) {
				matches = append(matches, path) // Add the full file path
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

// checkMultivalue checks if the package already exists in the map
func checkMultivalue(p *DebPackage, pkgs *map[string]DebPackage) bool {
	_, exists := (*pkgs)[p.Fields["Package"]]
	return exists
}

// generatePackageIndex generates the package index from the map of packages
func generatePackageIndex(pkgs *map[string]DebPackage, hashes *[]HASH) []byte {
	var buffer bytes.Buffer
	priorityFields := []string{
		"Package",
		"Source",
		"Version",
		"Installed-Size",
		"Maintainer",
		"Architecture",
		"Depends",
		"Recommends",
		"Suggests",
		"Homepage",
		"Section",
		"Priority",
		"Provides",
		"Description",
		"Size",
		"Filename",
		"MD5",
		"SHA256",
		"SHA1",
	}

	for _, p := range *pkgs {
		printedFields := make(map[string]bool)

		if containsHash(*hashes, "MD5") && containsHash(*hashes, "SHA1") && containsHash(*hashes, "SHA256") {
			p.CalculateAllHashes()
		}
		if containsHash(*hashes, "MD5") {
			p.MD5sum()
		}
		if containsHash(*hashes, "SHA1") {
			p.SHA1sum()
		}
		if containsHash(*hashes, "SHA256") {
			p.SHA256sum()
		}
		p.CalcSize()

		// Print priority fields first
		for _, field := range priorityFields {
			if value, exists := p.Fields[field]; exists {
				buffer.WriteString(fmt.Sprintf("%s: %s\n", field, value))
				printedFields[field] = true
			}
		}

		// Print remaining fields
		for key, value := range p.Fields {
			if !printedFields[key] {
				buffer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
			}
		}

		buffer.WriteString("\n")
	}
	return buffer.Bytes()
}
