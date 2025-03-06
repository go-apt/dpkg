package dpkg

// DebPackage represents the metadata of a Debian package
type DebPackage struct {
	Package          string
	Version          string
	Architecture     string
	Maintainer       string
	Description      string
	ShortDescription string
}
