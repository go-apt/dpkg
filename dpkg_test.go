package dpkg

import (
	"os"
	"testing"
)

// TestIsDebFile tests the IsDebFile function
func TestIsDebFile(t *testing.T) {
	tests := []struct {
		filePath string
		expected bool
	}{
		{"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb", true},
		{"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb", true},
		{"testdata/debs/invalid-package_1.0_amd64.deb", false},
	}

	for _, test := range tests {
		file, err := os.Open(test.filePath)
		if err != nil {
			t.Fatalf("Failed to open file %s: %v", test.filePath, err)
		}
		defer file.Close()

		result := IsDebFile(file)
		if result != test.expected {
			t.Errorf("IsDebFile(%s) = %v; want %v", test.filePath, result, test.expected)
		}
	}
}

// TestExtractControlFile tests the extractControlFile function
func TestExtractControlFile(t *testing.T) {
	tests := []struct {
		filePath string
		expected *DebPackage
	}{
		{
			"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb",
			&DebPackage{
				Package:      "vim-tiny",
				Version:      "2:9.1.1113-1",
				Architecture: "amd64",
				Maintainer:   "Debian Vim Maintainers <team+vim@tracker.debian.org>",
				Description:  "Vi IMproved - enhanced vi editor - compact version\nVim is an almost compatible version of the UNIX editor Vi.\n.\nThis package contains a minimal version of Vim compiled with no GUI and\na small subset of features. This package's sole purpose is to provide\nthe vi binary for base installations.\n.\nIf a vim binary is wanted, try one of the following more featureful\npackages: vim, vim-nox, vim-motif, or vim-gtk3.",
			},
		},
		{
			"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb",
			&DebPackage{
				Package:      "vim-tiny",
				Version:      "1:7.0-122+1etch5",
				Architecture: "amd64",
				Maintainer:   "Debian VIM Maintainers <pkg-vim-maintainers@lists.alioth.debian.org>",
				Description:  "Vi IMproved - enhanced vi editor - compact version\nVim is an almost compatible version of the UNIX editor Vi.\n.\nMany new features have been added: multi level undo, syntax\nhighlighting, command line history, on-line help, filename\ncompletion, block operations, folding, Unicode support, etc.\n.\nThis package contains a minimal version of vim compiled with no\nGUI and a small subset of features in order to keep small the\npackage size. This package does not depend on the vim-runtime\npackage, but installing it you will get its additional benefits\n(online documentation, plugins, ...).",
			},
		},
	}

	for _, test := range tests {
		d := Dpkg{}
		pkg, err := d.Info(test.filePath)
		if err != nil {
			t.Fatalf("Failed to extract control file from %s: %v", test.filePath, err)
		}

		if pkg.Package != test.expected.Package || pkg.Version != test.expected.Version || pkg.Architecture != test.expected.Architecture || pkg.Maintainer != test.expected.Maintainer || pkg.Description != test.expected.Description {
			t.Errorf("Extracted package metadata from %s does not match expected values", test.filePath)
		}
	}
}

// TestExtractGzipFile tests the extractGzipFile function
func TestExtractGzipFile(t *testing.T) {
	filePath := "testdata/compressed/control.tar.gz"
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = extractGzipFile(file)
	if err != nil {
		t.Fatalf("Failed to extract gzip file from %s: %v", filePath, err)
	}
}

// TestExtractXzFile tests the extractXzFile function
func TestExtractXzFile(t *testing.T) {
	filePath := "testdata/compressed/control.tar.xz"
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = extractXzFile(file)
	if err != nil {
		t.Fatalf("Failed to extract xz file from %s: %v", filePath, err)
	}
}
