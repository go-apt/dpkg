package dpkg

import (
	"os"
	"testing"
)

// TestParseControlFile tests the parseControlFile function
func TestParseControlFile(t *testing.T) {
	tests := []struct {
		filePath string
		expected *DebPackage
	}{
		{
			"testdata/control/control-1",
			&DebPackage{
				Package:      "vim-tiny",
				Version:      "1:7.0-122+1etch5",
				Architecture: "amd64",
				Maintainer:   "Debian VIM Maintainers <pkg-vim-maintainers@lists.alioth.debian.org>",
				Description:  "Vi IMproved - enhanced vi editor - compact version\n Vim is an almost compatible version of the UNIX editor Vi.\n .\n Many new features have been added: multi level undo, syntax\n highlighting, command line history, on-line help, filename\n completion, block operations, folding, Unicode support, etc.\n .\n This package contains a minimal version of vim compiled with no\n GUI and a small subset of features in order to keep small the\n package size. This package does not depend on the vim-runtime\n package, but installing it you will get its additional benefits\n (online documentation, plugins, ...).",
			},
		},
		{
			"testdata/control/control-2",
			&DebPackage{
				Package:      "vim-tiny",
				Version:      "2:9.1.1113-1",
				Architecture: "amd64",
				Maintainer:   "Debian Vim Maintainers <team+vim@tracker.debian.org>",
				Description:  "Vi IMproved - enhanced vi editor - compact version\n Vim is an almost compatible version of the UNIX editor Vi.\n .\n This package contains a minimal version of Vim compiled with no GUI and\n a small subset of features. This package's sole purpose is to provide\n the vi binary for base installations.\n .\n If a vim binary is wanted, try one of the following more featureful\n packages: vim, vim-nox, vim-motif, or vim-gtk3.",
			},
		},
	}

	for _, test := range tests {
		file, err := os.Open(test.filePath)
		if err != nil {
			t.Fatalf("Failed to open file %s: %v", test.filePath, err)
		}
		defer file.Close()

		pkg, err := parseControlFile(file)
		if err != nil {
			t.Fatalf("Failed to parse control file from %s: %v", test.filePath, err)
		}

		if pkg.Package != test.expected.Package || pkg.Version != test.expected.Version || pkg.Architecture != test.expected.Architecture || pkg.Maintainer != test.expected.Maintainer || pkg.Description != test.expected.Description {
			t.Errorf("Parsed package metadata from %s does not match expected values", test.filePath)
		}
	}
}

// TestReadPackageBlocks tests the readPackageBlocks function
func TestReadPackageBlocks(t *testing.T) {
	statusFile := "testdata/status"
	blocks, err := readPackageBlocks(statusFile)
	if err != nil {
		t.Fatalf("Failed to read package blocks from %s: %v", statusFile, err)
	}

	if len(blocks) == 0 {
		t.Errorf("Expected non-zero number of blocks, got %d", len(blocks))
	}

	// Additional checks to ensure blocks are not empty
	for i, block := range blocks {
		if len(block) == 0 {
			t.Errorf("Block %d is empty", i)
		}
	}
}

// TestParseStatusFile tests the parseStatusFile function
func TestParseStatusFile(t *testing.T) {
	statusFile := "testdata/status"
	expectedPackages := []DebPackage{
		{Package: "adduser", Version: "3.134"},
		{Package: "apparmor", Version: "3.0.8-3"},
		{Package: "apt", Version: "2.6.1"},
		{Package: "apt-listchanges", Version: "3.24"},
		{Package: "apt-mirror", Version: "0.5.4-2"},
	}

	packages, err := parseStatusFile(statusFile)
	if err != nil {
		t.Fatalf("Failed to parse status file from %s: %v", statusFile, err)
	}

	if len(packages) != len(expectedPackages) {
		t.Fatalf("Expected %d packages, got %d", len(expectedPackages), len(packages))
	}

	for i, pkg := range packages {
		expected := expectedPackages[i]
		if pkg.Package != expected.Package || pkg.Version != expected.Version {
			t.Errorf("Package %d does not match expected values. Got %+v, expected %+v", i, pkg, expected)
		}
	}
}
