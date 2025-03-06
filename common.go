package dpkg

import "errors"

var (
	ErrDebHeader     = errors.New("go-apt/dpkg: invalid debian package (ar magic header not matched)")
	ErrOpenDebFile   = errors.New("go-apt/dpkg: failed to open file")
	ErrNoControlFile = errors.New("go-apt/dpkg: failed to find control.tar file")
)
