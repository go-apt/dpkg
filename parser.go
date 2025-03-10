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
	pkg.Fields = make(map[string]string)

	// Read the content of the control file
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Convert the content to a string and split into lines
	lines := strings.Split(string(content), "\n")

	// Iterate over the lines to fill the struct
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Skip empty lines
		if line == "" {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			// If there's no ":", treat the entire line as a key with an empty value
			pkg.Fields[strings.TrimSpace(parts[0])] = ""
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Check if the next lines are continuations (start with space or tab)
		for i+1 < len(lines) && (strings.HasPrefix(lines[i+1], " ") || strings.HasPrefix(lines[i+1], "\t")) {
			i++
			value += "\n" + lines[i]
		}

		// Store the key-value pair in the map
		pkg.Fields[key] = value
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
