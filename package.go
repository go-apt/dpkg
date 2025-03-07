package dpkg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"reflect"
	"strings"
)

// DebPackage represents the metadata of a Debian package
// https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/lib/dpkg/parse.c?ref_type=heads#L53
type DebPackage struct {
	Package       string   `apt:"Package"`
	Version       string   `apt:"Version"`
	Architecture  string   `apt:"Architecture"`
	Maintainer    string   `apt:"Maintainer"`
	Essential     string   `apt:"Essential"`
	Protected     string   `apt:"Protected"`
	Status        string   `apt:"Status"`
	Priority      string   `apt:"Priority"`
	Section       string   `apt:"Section"`
	InstalledSize string   `apt:"Installed-Size"`
	Origin        string   `apt:"Origin"`
	MultiArch     string   `apt:"Multi-Arch"`
	Source        string   `apt:"Source"`
	Replaces      []string `apt:"Replaces"`
	Provides      []string `apt:"Provides"`
	Depends       []string `apt:"Depends"`
	PreDepends    []string `apt:"Pre-Depends"`
	Recommends    []string `apt:"Recommends"`
	Suggests      []string `apt:"Suggests"`
	Breaks        []string `apt:"Breaks"`
	Conflicts     []string `apt:"Conflicts"`
	Enhances      []string `apt:"Enhances"`
	Conffiles     string   `apt:"Conffiles"`
	Filename      string   `apt:"Filename"`
	Size          int64    `apt:"Size"`
	Description   string   `apt:"Description"`
	Homepage      string   `apt:"Homepage"`
	MD5Hash       string   `apt:"MD5sum"`
	SHA1Hash      string   `apt:"SHA1"`
	SHA256Hash    string   `apt:"SHA256"`
}

// readDebFile opens the .deb file associated with the DebPackage
func (dp *DebPackage) readDebFile() (io.ReadCloser, error) {
	if !dp.HasFilename() {
		return nil, ErrNoFilenameAvailable
	}
	return os.Open(dp.Filename)
}

// HasFilename checks if the DebPackage has a filename
func (dp *DebPackage) HasFilename() bool {
	return dp.Filename != ""
}

// ShortDescription returns the short description of the package
func (dp *DebPackage) ShortDescription() string {
	return strings.Split(dp.Description, "\n")[0]
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

// CalcSize calculates the size of the package file
func (dp *DebPackage) CalcSize() {
	if dp.HasFilename() {
		fileInfo, err := os.Stat(dp.Filename)
		if err != nil {
			dp.Size = 0
			return
		}
		dp.Size = fileInfo.Size()
	}
}

// GetAptTag retrieves the 'apt' tag value for a given field name in the DebPackage struct
func (dp DebPackage) GetAptTag(field interface{}) string {
	// Get the value of the field using reflection
	fieldValue := reflect.ValueOf(field)

	// Get the type of the struct
	t := reflect.TypeOf(dp)

	// Iterate over the fields of the struct to find the corresponding field
	for i := 0; i < t.NumField(); i++ {
		fieldStruct := t.Field(i)

		// Get the value of the struct field
		fieldValueStruct := reflect.ValueOf(dp).Field(i)

		// Compare the values
		if areEqual(fieldValueStruct.Interface(), fieldValue.Interface()) {
			// Return the value of the tag if the field is found
			return fieldStruct.Tag.Get("apt")
		}
	}

	// Return empty string if the field is not found
	return ""
}

// areEqual is a helper function to compare two values safely
func areEqual(a, b interface{}) bool {
	// If both are slices, compare element by element
	if reflect.TypeOf(a).Kind() == reflect.Slice && reflect.TypeOf(b).Kind() == reflect.Slice {
		return reflect.DeepEqual(a, b)
	}

	// Otherwise, compare directly
	return a == b
}
