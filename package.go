package dpkg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// DebPackage represents the metadata of a Debian package
// https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/lib/dpkg/parse.c?ref_type=heads#L53
type DebPackage struct {
	Package          string
	Version          string
	Architecture     string
	Maintainer       string
	Essential        string
	Protected        string
	Status           string
	Priority         string
	Section          string
	InstalledSize    string
	Origin           string
	Bugs             string
	MultiArch        string
	Source           string
	ConfigVersion    string
	Replaces         []string
	Provides         []string
	Depends          []string
	PreDepends       []string
	Recommends       []string
	Suggests         []string
	Breaks           []string
	Conflicts        []string
	Enhances         []string
	Conffiles        string
	Filename         string
	Size             string
	MSDOSFilename    string
	Description      string
	ShortDescription string
	TriggersPending  string
	TriggersAwaited  string
	Homepage         string
	MD5Hash          string
	SHA1Hash         string
	SHA256Hash       string
}

// readDebFile opens the .deb file associated with the DebPackage
func (dp *DebPackage) readDebFile() (io.ReadCloser, error) {
	if dp.Filename == "" {
		return nil, ErrNoFilenameAvailable
	}
	return os.Open(dp.Filename)
}

// CalculateAllHashes calculates the MD5, SHA1, and SHA256 hashes of the package content
func (dp *DebPackage) CalculateAllHashes() error {
	if dp.MD5Hash != "" && dp.SHA1Hash != "" && dp.SHA256Hash != "" {
		return nil
	}

	r, err := dp.readDebFile()
	if err != nil {
		return err
	}
	defer r.Close()

	// Create hashers
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	// Use a MultiWriter to calculate all hashes simultaneously
	multiWriter := io.MultiWriter(md5Hash, sha1Hash, sha256Hash)
	if _, err := io.Copy(multiWriter, r); err != nil {
		return err
	}

	// Update hash fields
	dp.MD5Hash = hex.EncodeToString(md5Hash.Sum(nil))
	dp.SHA1Hash = hex.EncodeToString(sha1Hash.Sum(nil))
	dp.SHA256Hash = hex.EncodeToString(sha256Hash.Sum(nil))

	return nil
}

// calculateSingleHash calculates a single hash for the package content
func (dp *DebPackage) calculateSingleHash(hashFunc func() hash.Hash, hashField *string) error {
	r, err := dp.readDebFile()
	if err != nil {
		return err
	}
	defer r.Close()

	hash := hashFunc()
	if _, err := io.Copy(hash, r); err != nil {
		return err
	}

	*hashField = hex.EncodeToString(hash.Sum(nil))
	return nil
}

// MD5sum returns the MD5 hash of the package content
func (dp *DebPackage) MD5sum() string {
	if dp.MD5Hash == "" {
		if err := dp.calculateSingleHash(md5.New, &dp.MD5Hash); err != nil {
			return ""
		}
	}
	return dp.MD5Hash
}

// SHA1sum returns the SHA1 hash of the package content
func (dp *DebPackage) SHA1sum() string {
	if dp.SHA1Hash == "" {
		if err := dp.calculateSingleHash(sha1.New, &dp.SHA1Hash); err != nil {
			return ""
		}
	}
	return dp.SHA1Hash
}

// SHA256sum returns the SHA256 hash of the package content
func (dp *DebPackage) SHA256sum() string {
	if dp.SHA256Hash == "" {
		if err := dp.calculateSingleHash(sha256.New, &dp.SHA256Hash); err != nil {
			return ""
		}
	}

	return dp.SHA256Hash
}
