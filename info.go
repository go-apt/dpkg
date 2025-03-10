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

// readArchive reads the .deb file and extracts the control file
func (d *Dpkg) readArchive(debFile string) (*DebPackage, error) {
	file, err := os.Open(debFile)
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
			pkg.Fields["Filename"] = debFile
			return pkg, nil
		}
	}

	return nil, ErrNoControlFile
}

// extractControlFile extracts the control file from the archive
func extractControlFile(filename string, arReader io.Reader) (*DebPackage, error) {
	var err error
	var uncompressedData io.Reader

	// Handle different compression formats
	// https://manpages.debian.org/buster/dpkg-dev/deb.5.en.html#FORMAT
	switch path.Ext(filename) {
	case ".tar":
		uncompressedData = arReader
	case ".gz":
		uncompressedData, err = extractGzipFile(arReader)
	case ".xz":
		uncompressedData, err = extractXzFile(arReader)
	default:
		return nil, fmt.Errorf("go-apt/dpkg: unsupported compression format for %s", filename)
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
		return nil, fmt.Errorf("go-apt/dpkg: failed to create gzip reader: %w", err)
	}
	return gzReader, nil
}

// extractXzFile extracts an xz compressed file
func extractXzFile(compressedReader io.Reader) (io.Reader, error) {
	xzReader, err := xz.NewReader(compressedReader)
	if err != nil {
		return nil, fmt.Errorf("go-apt/dpkg: failed to create xz reader: %w", err)
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
			return nil, fmt.Errorf("go-apt/dpkg: failed to read tar header: %w", err)
		}

		if tarHeader.Name == "./control" || tarHeader.Name == "control" {
			return parseControlFile(tarReader)
		}
	}
	return nil, ErrNoControlFile
}
