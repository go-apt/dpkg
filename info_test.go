package dpkg

import (
	"fmt"
	"os"
	"testing"
)

// TestExtractControlFile tests the extractControlFile function
func TestExtractControlFile(t *testing.T) {
	tests := []struct {
		filePath string
		expected *DebPackage
	}{
		{
			"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb",
			&DebPackage{
				Fields: map[string]string{
					"Package":      "vim-tiny",
					"Version":      "2:9.1.1113-1",
					"Architecture": "amd64",
					"Maintainer":   "Debian Vim Maintainers <team+vim@tracker.debian.org>",
					"Description":  "Vi IMproved - enhanced vi editor - compact version\n Vim is an almost compatible version of the UNIX editor Vi.\n .\n This package contains a minimal version of Vim compiled with no GUI and\n a small subset of features. This package's sole purpose is to provide\n the vi binary for base installations.\n .\n If a vim binary is wanted, try one of the following more featureful\n packages: vim, vim-nox, vim-motif, or vim-gtk3.",
				},
			},
		},
		{
			"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb",
			&DebPackage{
				Fields: map[string]string{
					"Package":      "vim-tiny",
					"Version":      "1:7.0-122+1etch5",
					"Architecture": "amd64",
					"Maintainer":   "Debian VIM Maintainers <pkg-vim-maintainers@lists.alioth.debian.org>",
					"Description":  "Vi IMproved - enhanced vi editor - compact version\n Vim is an almost compatible version of the UNIX editor Vi.\n .\n Many new features have been added: multi level undo, syntax\n highlighting, command line history, on-line help, filename\n completion, block operations, folding, Unicode support, etc.\n .\n This package contains a minimal version of vim compiled with no\n GUI and a small subset of features in order to keep small the\n package size. This package does not depend on the vim-runtime\n package, but installing it you will get its additional benefits\n (online documentation, plugins, ...).",
				},
			},
		},
	}

	for _, test := range tests {
		d := Dpkg{}
		pkg, err := d.Info(test.filePath)
		if err != nil {
			t.Fatalf("Failed to extract control file from %s: %v", test.filePath, err)
		}

		if pkg.Fields["Package"] != test.expected.Fields["Package"] || pkg.Fields["Version"] != test.expected.Fields["Version"] || pkg.Fields["Architecture"] != test.expected.Fields["Architecture"] || pkg.Fields["Maintainer"] != test.expected.Fields["Maintainer"] || pkg.Fields["Description"] != test.expected.Fields["Description"] {
			fmt.Println(pkg.Fields["Description"])
			fmt.Println(test.expected.Fields["Description"])
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
