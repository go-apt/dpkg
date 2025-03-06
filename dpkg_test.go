package dpkg

import (
	"testing"
)

// TestIsDebFile tests the IsDebFile function
func TestIsDebFile(t *testing.T) {
	d := Dpkg{}

	tests := []struct {
		filePath string
		expected bool
	}{
		{"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb", true},
		{"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb", true},
		{"testdata/debs/invalid-package_1.0_amd64.deb", false},
	}

	for _, test := range tests {
		result := d.IsDebFile(test.filePath)
		if result != test.expected {
			t.Errorf("IsDebFile(%s) = %v; want %v", test.filePath, result, test.expected)
		}
	}
}
