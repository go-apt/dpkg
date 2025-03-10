package dpkg

import (
	"path/filepath"
	"testing"
)

// TestCalculateHashes tests the calculateHashes function
func TestCalculateHashes(t *testing.T) {
	tests := []struct {
		debFile        string
		expectedMD5    string
		expectedSHA1   string
		expectedSHA256 string
	}{
		{
			"testdata/debs/invalid-package_1.0_amd64.deb",
			"ec649caf62773d8324b6ef96dac1c572",
			"262fc1ab65f071f015912c47b38507f6364aadd8",
			"c00fbabe4192ff18b5c80eff33488cc6015b3600b9fd4fb270a66b05164fa815",
		},
		{
			"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb",
			"70ac9e55bb99b0e1b5d22f105e099ce0",
			"8ea6af742f8b673a27a49e77c7f452b3eed4f631",
			"d40b30835087af1affdd9c949848757b013bde7b140bd033b75e1c1aef9597c9",
		},
		{
			"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb",
			"9c7b57fa57c18835d54945e932fd1120",
			"8fe1a32cadf5115eee1257ec4f24f10738dc014a",
			"c4e554b7b5a25692210b8730d3183f66b2286b8429570949598ce56ab671ca13",
		},
	}

	for _, test := range tests {
		t.Run(filepath.Base(test.debFile), func(t *testing.T) {
			pkg := &DebPackage{
				Fields: map[string]string{
					"Filename": test.debFile,
				},
			}

			err := pkg.CalculateAllHashes()
			if err != nil {
				t.Fatalf("Failed to calculate hashes for %s: %v", test.debFile, err)
			}

			if pkg.Fields["MD5"] != test.expectedMD5 {
				t.Errorf("Expected MD5 hash %s, got %s", test.expectedMD5, pkg.Fields["MD5"])
			}

			if pkg.Fields["SHA1"] != test.expectedSHA1 {
				t.Errorf("Expected SHA1 hash %s, got %s", test.expectedSHA1, pkg.Fields["SHA1"])
			}

			if pkg.Fields["SHA256"] != test.expectedSHA256 {
				t.Errorf("Expected SHA256 hash %s, got %s", test.expectedSHA256, pkg.Fields["SHA256"])
			}
		})
	}
}

// TestMD5sum tests the MD5sum function
func TestMD5sum(t *testing.T) {
	tests := []struct {
		debFile     string
		expectedMD5 string
	}{
		{
			"testdata/debs/invalid-package_1.0_amd64.deb",
			"ec649caf62773d8324b6ef96dac1c572",
		},
		{
			"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb",
			"70ac9e55bb99b0e1b5d22f105e099ce0",
		},
		{
			"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb",
			"9c7b57fa57c18835d54945e932fd1120",
		},
	}

	for _, test := range tests {
		t.Run(filepath.Base(test.debFile), func(t *testing.T) {
			pkg := &DebPackage{
				Fields: map[string]string{
					"Filename": test.debFile,
				},
			}

			md5sum := pkg.MD5sum()
			if md5sum != test.expectedMD5 {
				t.Errorf("Expected MD5 hash %s, got %s", test.expectedMD5, md5sum)
			}
		})
	}
}

// TestSHA1sum tests the SHA1sum function
func TestSHA1sum(t *testing.T) {
	tests := []struct {
		debFile      string
		expectedSHA1 string
	}{
		{
			"testdata/debs/invalid-package_1.0_amd64.deb",
			"262fc1ab65f071f015912c47b38507f6364aadd8",
		},
		{
			"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb",
			"8ea6af742f8b673a27a49e77c7f452b3eed4f631",
		},
		{
			"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb",
			"8fe1a32cadf5115eee1257ec4f24f10738dc014a",
		},
	}

	for _, test := range tests {
		t.Run(filepath.Base(test.debFile), func(t *testing.T) {
			pkg := &DebPackage{
				Fields: map[string]string{
					"Filename": test.debFile,
				},
			}

			sha1sum := pkg.SHA1sum()
			if sha1sum != test.expectedSHA1 {
				t.Errorf("Expected SHA1 hash %s, got %s", test.expectedSHA1, sha1sum)
			}
		})
	}
}

// TestSHA256sum tests the SHA256sum function
func TestSHA256sum(t *testing.T) {
	tests := []struct {
		debFile        string
		expectedSHA256 string
	}{
		{
			"testdata/debs/invalid-package_1.0_amd64.deb",
			"c00fbabe4192ff18b5c80eff33488cc6015b3600b9fd4fb270a66b05164fa815",
		},
		{
			"testdata/debs/vim-tiny_7.0-122+1etch5_amd64.deb",
			"d40b30835087af1affdd9c949848757b013bde7b140bd033b75e1c1aef9597c9",
		},
		{
			"testdata/debs/vim-tiny_9.1.1113-1_amd64.deb",
			"c4e554b7b5a25692210b8730d3183f66b2286b8429570949598ce56ab671ca13",
		},
	}

	for _, test := range tests {
		t.Run(filepath.Base(test.debFile), func(t *testing.T) {
			pkg := &DebPackage{
				Fields: map[string]string{
					"Filename": test.debFile,
				},
			}

			sha256sum := pkg.SHA256sum()
			if sha256sum != test.expectedSHA256 {
				t.Errorf("Expected SHA256 hash %s, got %s", test.expectedSHA256, sha256sum)
			}
		})
	}
}

func BenchmarkCalcHashesIndividually(b *testing.B) {
	test := struct {
		debFile        string
		expectedMD5    string
		expectedSHA1   string
		expectedSHA256 string
	}{
		debFile:        "testdata/debs/invalid-package_1.0_amd64.deb",
		expectedMD5:    "ec649caf62773d8324b6ef96dac1c572",
		expectedSHA1:   "262fc1ab65f071f015912c47b38507f6364aadd8",
		expectedSHA256: "c00fbabe4192ff18b5c80eff33488cc6015b3600b9fd4fb270a66b05164fa815",
	}

	pkg := &DebPackage{
		Fields: map[string]string{
			"Filename": test.debFile,
		},
	}

	for i := 0; i < b.N; i++ {
		pkg.MD5sum()
		pkg.SHA1sum()
		pkg.SHA256sum()

		if pkg.Fields["MD5"] != test.expectedMD5 {
			b.Errorf("Expected MD5 hash %s, got %s", test.expectedMD5, pkg.Fields["MD5"])
		}

		if pkg.Fields["SHA1"] != test.expectedSHA1 {
			b.Errorf("Expected SHA1 hash %s, got %s", test.expectedSHA1, pkg.Fields["SHA1"])
		}

		if pkg.Fields["SHA256"] != test.expectedSHA256 {
			b.Errorf("Expected SHA256 hash %s, got %s", test.expectedSHA256, pkg.Fields["SHA256"])
		}
	}
}

func BenchmarkCalcAllHashes(b *testing.B) {
	test := struct {
		debFile        string
		expectedMD5    string
		expectedSHA1   string
		expectedSHA256 string
	}{
		debFile:        "testdata/debs/invalid-package_1.0_amd64.deb",
		expectedMD5:    "ec649caf62773d8324b6ef96dac1c572",
		expectedSHA1:   "262fc1ab65f071f015912c47b38507f6364aadd8",
		expectedSHA256: "c00fbabe4192ff18b5c80eff33488cc6015b3600b9fd4fb270a66b05164fa815",
	}

	pkg := &DebPackage{
		Fields: map[string]string{
			"Filename": test.debFile,
		},
	}

	for i := 0; i < b.N; i++ {
		pkg.CalculateAllHashes()
		if pkg.Fields["MD5"] != test.expectedMD5 {
			b.Errorf("Expected MD5 hash %s, got %s", test.expectedMD5, pkg.Fields["MD5"])
		}

		if pkg.Fields["SHA1"] != test.expectedSHA1 {
			b.Errorf("Expected SHA1 hash %s, got %s", test.expectedSHA1, pkg.Fields["SHA1"])
		}

		if pkg.Fields["SHA256"] != test.expectedSHA256 {
			b.Errorf("Expected SHA256 hash %s, got %s", test.expectedSHA256, pkg.Fields["SHA256"])
		}
	}
}
