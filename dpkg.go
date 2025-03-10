package dpkg

import (
	"os"
	"strings"
)

// Dpkg represents a Debian package manager
type Dpkg struct {
	StatusFileLocation string
}

// NewDpkg creates a new instance of Dpkg
func NewDpkg() *Dpkg {
	return &Dpkg{
		StatusFileLocation: DPKG_DATABASE,
	}
}

// Info retrieves the metadata of a Debian package
func (d *Dpkg) Info(debFile string) (*DebPackage, error) {
	pkg, err := d.readArchive(debFile)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// List lists packages from the default dpkg database
func (d *Dpkg) List() ([]DebPackage, error) {
	return parseStatusFile(d.StatusFileLocation)
}

// ListGrep lists packages from the default dpkg database that match the given package name
func (d *Dpkg) ListGrep(pkgName string) ([]DebPackage, error) {
	packages, err := parseStatusFile(d.StatusFileLocation)
	if err != nil {
		return nil, err
	}

	var filteredPackages []DebPackage
	for _, pkg := range packages {
		if strings.Contains(pkg.Fields["Package"], pkgName) {
			filteredPackages = append(filteredPackages, pkg)
		}
	}

	return filteredPackages, nil
}

// IsDebFile checks if the file is a valid .deb package
// https://manpages.debian.org/buster/dpkg-dev/deb.5.en.html#FORMAT
func (d *Dpkg) IsDebFile(debFile string) bool {
	file, err := os.Open(debFile)
	if err != nil {
		return false
	}
	defer file.Close()

	magicValue := "!<arch>\ndebian-binary"
	magic := make([]byte, len(magicValue))

	_, err = file.Read(magic)
	if err != nil {
		return false
	}

	// Check the "magic value" for the ar format
	return strings.HasPrefix(string(magic), magicValue)
}
