package dpkg

import "errors"

const (
	DPKG_DATABASE = "/var/lib/dpkg/status"
)

var (
	ErrDebHeader        = errors.New("go-apt/dpkg: invalid debian package (ar magic header not matched)")
	ErrNoControlFile    = errors.New("go-apt/dpkg: failed to find control.tar file")
	ErrNoDpkgStatusFile = errors.New("go-apt/dpkg: failed to find " + DPKG_DATABASE + " file")
)
