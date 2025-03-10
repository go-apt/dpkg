package dpkg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"strconv"
	"strings"
)

// DebPackage represents the metadata of a Debian package
// https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/lib/dpkg/parse.c?ref_type=heads#L53
type DebPackage struct {
	Fields map[string]string
}

// readDebFile opens the .deb file associated with the DebPackage
func (dp *DebPackage) readDebFile() (io.ReadCloser, error) {
	if !dp.HasFilename() {
		return nil, ErrNoFilenameAvailable
	}
	return os.Open(dp.Fields["Filename"])
}

// HasFilename checks if the DebPackage has a filename
func (dp *DebPackage) HasFilename() bool {
	return dp.Fields["Filename"] != ""
}

// ShortDescription returns the short description of the package
func (dp *DebPackage) ShortDescription() string {
	return strings.Split(dp.Fields["Description"], "\n")[0]
}

// CalculateAllHashes calculates the MD5, SHA1, and SHA256 hashes of the package content
func (dp *DebPackage) CalculateAllHashes() error {
	r, err := dp.readDebFile()
	if err != nil {
		return err
	}
	defer r.Close()

	// Create hashers
	MD5 := md5.New()
	SHA1 := sha1.New()
	SHA256 := sha256.New()

	// Use a MultiWriter to calculate all hashes simultaneously
	multiWriter := io.MultiWriter(MD5, SHA1, SHA256)
	if _, err := io.Copy(multiWriter, r); err != nil {
		return err
	}

	// Update hash fields
	dp.Fields["MD5"] = hex.EncodeToString(MD5.Sum(nil))
	dp.Fields["SHA1"] = hex.EncodeToString(SHA1.Sum(nil))
	dp.Fields["SHA256"] = hex.EncodeToString(SHA256.Sum(nil))

	return nil
}

// calculateSingleHash calculates a single hash for the package content
func (dp *DebPackage) calculateSingleHash(hashFunc func() hash.Hash, hashField string) error {
	r, err := dp.readDebFile()
	if err != nil {
		return err
	}
	defer r.Close()

	hash := hashFunc()
	if _, err := io.Copy(hash, r); err != nil {
		return err
	}

	dp.Fields[hashField] = hex.EncodeToString(hash.Sum(nil))
	return nil
}

// MD5sum returns the MD5 hash of the package content
func (dp *DebPackage) MD5sum() string {
	if err := dp.calculateSingleHash(md5.New, "MD5"); err != nil {
		return ""
	}
	return dp.Fields["MD5"]
}

// SHA1sum returns the SHA1 hash of the package content
func (dp *DebPackage) SHA1sum() string {
	if err := dp.calculateSingleHash(sha1.New, "SHA1"); err != nil {
		return ""
	}
	return dp.Fields["SHA1"]
}

// SHA256sum returns the SHA256 hash of the package content
func (dp *DebPackage) SHA256sum() string {
	if err := dp.calculateSingleHash(sha256.New, "SHA256"); err != nil {
		return ""
	}

	return dp.Fields["SHA256"]
}

// CalcSize calculates the size of the package file
func (dp *DebPackage) CalcSize() {
	if dp.HasFilename() {
		fileInfo, err := os.Stat(dp.Fields["Filename"])
		if err != nil {
			dp.Fields["Size"] = "0"
			return
		}
		dp.Fields["Size"] = strconv.FormatInt(fileInfo.Size(), 10)
	}
}
