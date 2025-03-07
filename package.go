package dpkg

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
	MD5sum           string
	MSDOSFilename    string
	Description      string
	ShortDescription string
	TriggersPending  string
	TriggersAwaited  string
	Homepage         string
}
