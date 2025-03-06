package dpkg

import (
	"io"
	"strings"
)

// DebPackage represents the metadata of a Debian package
type DebPackage struct {
	Package      string
	Version      string
	Architecture string
	Maintainer   string
	Description  string
}

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
			descriptionLines = append(descriptionLines, strings.TrimSpace(strings.TrimPrefix(line, "Description:")))
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
