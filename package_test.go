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
				Description:  "Vi IMproved - enhanced vi editor - compact version\nVim is an almost compatible version of the UNIX editor Vi.\n.\nMany new features have been added: multi level undo, syntax\nhighlighting, command line history, on-line help, filename\ncompletion, block operations, folding, Unicode support, etc.\n.\nThis package contains a minimal version of vim compiled with no\nGUI and a small subset of features in order to keep small the\npackage size. This package does not depend on the vim-runtime\npackage, but installing it you will get its additional benefits\n(online documentation, plugins, ...).",
			},
		},
		{
			"testdata/control/control-2",
			&DebPackage{
				Package:      "vim-tiny",
				Version:      "2:9.1.1113-1",
				Architecture: "amd64",
				Maintainer:   "Debian Vim Maintainers <team+vim@tracker.debian.org>",
				Description:  "Vi IMproved - enhanced vi editor - compact version\nVim is an almost compatible version of the UNIX editor Vi.\n.\nThis package contains a minimal version of Vim compiled with no GUI and\na small subset of features. This package's sole purpose is to provide\nthe vi binary for base installations.\n.\nIf a vim binary is wanted, try one of the following more featureful\npackages: vim, vim-nox, vim-motif, or vim-gtk3.",
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
