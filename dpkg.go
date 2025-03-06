package dpkg

import (
	"os"
	"strings"
)

// Dpkg represents a Debian package manager
type Dpkg struct {
	StatusFileLocation string
}

// Info retrieves the metadata of a Debian package
func (d *Dpkg) Info(debFile string) (*DebPackage, error) {
	pkg, err := d.readArchive(debFile)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// List lists packages from default dpkg database
func (d *Dpkg) List() ([]DebPackage, error) {
	return parseStatusFile(DPKG_DATABASE)
}

// List lists packages from custom dpkg database
func (d *Dpkg) ListCustom(statusFile string) ([]DebPackage, error) {
	return parseStatusFile(statusFile)
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
