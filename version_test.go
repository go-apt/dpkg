package dpkg

import (
	"testing"
)

// TestCompareVersions tests the compareVersions function
func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1:1.0-1", "1:1.0-1", 0},
		{"1:1.0-1", "1:1.0-2", -1},
		{"1:1.0-2", "1:1.0-1", 1},
		{"1:1.0-1", "2:1.0-1", -1},
		{"2:1.0-1", "1:1.0-1", 1},
		{"1.0-1", "1.0-1", 0},
		{"1.0-1", "1.0-2", -1},
		{"1.0-2", "1.0-1", 1},
		{"1.0", "1.0-1", -1},
		{"1.0-1", "1.0", 1},
		{"1:7.0-122+1etch5", "2:9.1.1113-1", -1},
		{"2:9.1.1113-1", "1:7.0-122+1etch5", 1},
		{"9.1.1113-1", "7.0-122+1etch5", 1},
		{"7.0-122+1etch5", "7.0-122+1etch5", 0},
		{"7.0-122+1etch5", "7.0-123+1etch5", -1},
		{"7.0-122+2etch5", "7.0-122+1etch5", 1},
		{"7.0-122+1etch6", "7.0-122+1etch5", 1},
		{"7.0-122+1etch5", "7.0-122+2etch5", -1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_"+tt.v2, func(t *testing.T) {
			d := &Dpkg{}
			got := d.CompareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d; want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

// TestSplitVersion tests the splitVersion function
func TestSplitVersion(t *testing.T) {
	tests := []struct {
		version      string
		wantEpoch    int
		wantMain     string
		wantRevision string
	}{
		{"1:1.0-1", 1, "1.0", "1"},
		{"1.0-1", 0, "1.0", "1"},
		{"1.0", 0, "1.0", ""},
		{"2:2.0", 2, "2.0", ""},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			gotEpoch, gotMain, gotRevision := splitVersion(tt.version)
			if gotEpoch != tt.wantEpoch || gotMain != tt.wantMain || gotRevision != tt.wantRevision {
				t.Errorf("splitVersion(%q) = (%d, %q, %q); want (%d, %q, %q)",
					tt.version, gotEpoch, gotMain, gotRevision, tt.wantEpoch, tt.wantMain, tt.wantRevision)
			}
		})
	}
}

// TestCompareDebianVersion tests the compareDebianVersion function
func TestCompareDebianVersion(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0", "1.0", 0},
		{"1.0", "1.1", -1},
		{"1.1", "1.0", 1},
		{"1.0a", "1.0b", -1},
		{"1.0b", "1.0a", 1},
		{"1.0-1", "1.0-2", -1},
		{"1.0-2", "1.0-1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_"+tt.v2, func(t *testing.T) {
			got := compareDebianVersion(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("compareDebianVersion(%q, %q) = %d; want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

// TestSplitVersionParts tests the splitVersionParts function
func TestSplitVersionParts(t *testing.T) {
	tests := []struct {
		version string
		want    []string
	}{
		{"1.0", []string{"1", ".", "0"}},
		{"1.0-1", []string{"1", ".", "0", "-", "1"}},
		{"1.0a", []string{"1", ".", "0", "a"}},
		{"1.0~beta1", []string{"1", ".", "0", "~", "b", "e", "t", "a", "1"}},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := splitVersionParts(tt.version)
			if !equalStringSlices(got, tt.want) {
				t.Errorf("splitVersionParts(%q) = %v; want %v", tt.version, got, tt.want)
			}
		})
	}
}

// TestGetPart tests the getPart function
func TestGetPart(t *testing.T) {
	tests := []struct {
		parts []string
		index int
		want  string
	}{
		{[]string{"1", ".", "0"}, 0, "1"},
		{[]string{"1", ".", "0"}, 1, "."},
		{[]string{"1", ".", "0"}, 2, "0"},
		{[]string{"1", ".", "0"}, 3, ""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := getPart(tt.parts, tt.index)
			if got != tt.want {
				t.Errorf("getPart(%v, %d) = %q; want %q", tt.parts, tt.index, got, tt.want)
			}
		})
	}
}

// Helper function to compare two slices of strings
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
