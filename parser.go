package dpkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// parseControlFile parses the control file content into a DebPackage struct
func parseControlFile(reader io.Reader) (*DebPackage, error) {
	pkg := &DebPackage{}

	// Read the content of the control file
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Convert the content to a string and split into lines
	lines := strings.Split(string(content), "\n")

	// Iterate over the lines to fill the struct
	var descriptionLines []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch {
		case strings.HasPrefix(line, "Package:"):
			pkg.Package = strings.TrimSpace(strings.TrimPrefix(line, "Package:"))
		case strings.HasPrefix(line, "Version:"):
			pkg.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
		case strings.HasPrefix(line, "Architecture:"):
			pkg.Architecture = strings.TrimSpace(strings.TrimPrefix(line, "Architecture:"))
		case strings.HasPrefix(line, "Maintainer:"):
			pkg.Maintainer = strings.TrimSpace(strings.TrimPrefix(line, "Maintainer:"))
		case strings.HasPrefix(line, "Description:"):
			firstLine := strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
			pkg.ShortDescription = firstLine
			descriptionLines = append(descriptionLines, firstLine)
			// Continue reading the following lines as part of the description
			for i+1 < len(lines) && (strings.HasPrefix(lines[i+1], " ") || strings.HasPrefix(lines[i+1], "\t")) {
				i++
				descriptionLines = append(descriptionLines, strings.TrimSpace(lines[i]))
			}
			pkg.Description = strings.Join(descriptionLines, "\n")
		}
	}

	return pkg, nil
}

// readPackageBlocks reads the statusFile and returns a slice of []byte,
// where each []byte contains the content of a package block.
func readPackageBlocks(statusFile string) ([][]byte, error) {
	file, err := os.Open(statusFile)
	if err != nil {
		return nil, ErrNoDpkgStatusFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var blocks [][]byte
	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()

		// Check if it's the end of the block (blank line)
		if line != "" {
			// Add the line to the buffer
			buffer.WriteString(line)
			buffer.WriteString("\n") // Keep line breaks
		} else {
			// Create a copy of the buffer and add it to the slice of blocks
			block := make([]byte, buffer.Len())
			copy(block, buffer.Bytes()) // Copy the buffer content to a new slice
			blocks = append(blocks, block)
			buffer.Reset() // Clear the buffer for the next block
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading the status file: %w", err)
	}

	// Add the last block if the file does not end with a blank line
	if buffer.Len() > 0 {
		block := make([]byte, buffer.Len())
		copy(block, buffer.Bytes()) // Copy the buffer content to a new slice
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// parseStatusFile reads the statusFile file and returns a list of packages
func parseStatusFile(statusFile string) ([]DebPackage, error) {
	blocks, err := readPackageBlocks(statusFile)
	if err != nil {
		return nil, err
	}

	var packages []DebPackage

	for _, block := range blocks {
		// Create a new io.Reader for each block
		reader := bytes.NewReader(block)
		pkg, err := parseControlFile(reader)
		if err != nil {
			return nil, fmt.Errorf("error processing the block: %w", err)
		}
		packages = append(packages, *pkg)
	}

	return packages, nil
}
