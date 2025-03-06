package dpkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

// Dpkg represents a Debian package manager
type Dpkg struct{}

// Info retrieves the metadata of a Debian package
func (d *Dpkg) Info(debFile string) (*DebPackage, error) {
	pkg, err := d.readArchive(debFile)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// readArchive reads the .deb file and extracts the control file
func (d *Dpkg) readArchive(debFile string) (*DebPackage, error) {
	file, err := d.openFile(debFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := ar.NewReader(file)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(header.Name, "control.tar") {
			pkg, err := extractControlFile(header.Name, reader)
			if err != nil {
				return nil, err
			}
			return pkg, nil
		}
	}

	return nil, ErrNoControlFile
}

// openFile opens the .deb file for reading
func (d *Dpkg) openFile(debFile string) (*os.File, error) {
	file, err := os.Open(debFile)
	if err != nil {
		return nil, ErrOpenDebFile
	}

	if !IsDebFile(file) {
		return nil, ErrDebHeader
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// IsDebFile checks if the file is a valid .deb package
// https://manpages.debian.org/buster/dpkg-dev/deb.5.en.html#FORMAT
func IsDebFile(file *os.File) bool {
	magicValue := "!<arch>\ndebian-binary"
	magic := make([]byte, len(magicValue))
	_, err := file.Read(magic)
	if err != nil {
		return false
	}

	// Check the "magic value" for the ar format
	return strings.HasPrefix(string(magic), magicValue)
}

// extractControlFile extracts the control file from the archive
func extractControlFile(filename string, arReader *ar.Reader) (*DebPackage, error) {
	var err error
	var uncompressedData io.Reader

	// https://manpages.debian.org/buster/dpkg-dev/deb.5.en.html#FORMAT
	switch path.Ext(filename) {
	case ".tar":
		uncompressedData = arReader
	case ".gz":
		uncompressedData, err = extractGzipFile(arReader)
	case ".xz":
		uncompressedData, err = extractXzFile(arReader)
	default:
		return nil, fmt.Errorf("go-apt/dpkg: does not know how to handle compression format for %s", filename)
	}

	if err != nil {
		return nil, err
	}

	return extractControlFromTarFile(uncompressedData)
}

// extractGzipFile extracts a gzip compressed file
func extractGzipFile(compressedReader io.Reader) (io.Reader, error) {
	gzReader, err := gzip.NewReader(compressedReader)
	if err != nil {
		return nil, err
	}
	return gzReader, nil
}

// extractXzFile extracts an xz compressed file
func extractXzFile(compressedReader io.Reader) (io.Reader, error) {
	xzReader, err := xz.NewReader(compressedReader)
	if err != nil {
		return nil, err
	}
	return xzReader, nil
}

// extractControlFromTarFile extracts the control file from a tar archive
func extractControlFromTarFile(uncompressedReader io.Reader) (*DebPackage, error) {
	tarReader := tar.NewReader(uncompressedReader)
	for {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if tarHeader.Name == "./control" || tarHeader.Name == "control" {
			return parseControlFile(tarReader)
		}
	}
	return nil, ErrNoControlFile
}
