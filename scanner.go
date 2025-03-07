package dpkg

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// PackagesScanner represents a scanner for Debian packages
type PackagesScanner struct {
	RootDir      string
	Arch         string
	Type         string
	Hashes       []string
	Multiversion bool
}

// NewPackagesScanner creates a new instance of PackagesScanner
func NewPackagesScanner(dir string) *PackagesScanner {
	return &PackagesScanner{
		RootDir:      dir,
		Arch:         "",
		Type:         "deb",
		Hashes:       []string{"md5", "sha1", "sha256"},
		Multiversion: false,
	}
}

// ScanPackages scans the directory for packages matching the criteria
func (ps *PackagesScanner) ScanPackages() ([]byte, error) {
	// List the files in the directory that match the filter
	files, err := findFiles(ps.RootDir, ps.createFilter())
	if err != nil {
		return nil, fmt.Errorf("error finding files: %v", err)
	}

	pkgs := make(map[string]DebPackage)
	d := NewDpkg()

	for _, file := range files {
		p, err := d.Info(file)
		if err != nil {
			// Let the user know that there was an error parsing the package
			// but continue to the next package
			fmt.Printf("go-apt/dpkg: error parsing package: '%s' - %v\n", file, err)
			continue
		}

		if checkMultivalue(p, &pkgs) {
			op := pkgs[p.Package]

			if ps.Multiversion {
				fmt.Printf("multiversion enabled; adding repeated Package '%s'\n", p.Package)
				pkgs[p.Package+p.Version] = *p
				continue
			}

			if d.CompareVersions(p.Version, op.Version) > 0 {
				fmt.Printf("package '%s' (filename '%s') is repeat but newer version; ignored version '%s'!\n", p.Package, p.Filename, op.Version)
				pkgs[p.Package] = *p
			} else {
				fmt.Printf("package '%s' (filename '%s') is repeat but older; ignored version '%s'!\n", p.Package, p.Filename, p.Version)
				continue
			}
		}

		pkgs[p.Package] = *p
	}

	return generatePackageIndex(&pkgs), nil
}

// createFilter creates a regular expression filter based on architecture and type
func (ps *PackagesScanner) createFilter() *regexp.Regexp {
	if ps.Arch != "" {
		// Create the regular expression with the architecture
		pattern := fmt.Sprintf("_(?:all|%s)\\.%s$", ps.Arch, ps.Type)
		return regexp.MustCompile(pattern)
	} else {
		// Create the regular expression without the architecture
		pattern := fmt.Sprintf("\\.%s$", ps.Type)
		return regexp.MustCompile(pattern)
	}
}

// findFiles finds files that match the filter, including subdirectories
func findFiles(directory string, filter *regexp.Regexp) ([]string, error) {
	var matches []string

	// Use filepath.WalkDir to traverse the directory and subdirectories
	err := filepath.WalkDir(directory, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err // Return errors encountered during the walk
		}

		// Check if the entry is a file (ignore directories)
		if !entry.IsDir() {
			fileName := entry.Name()
			// Check if the file name matches the filter
			if filter.MatchString(fileName) {
				matches = append(matches, path) // Add the full file path
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

// checkMultivalue checks if the package already exists in the map
func checkMultivalue(p *DebPackage, pkgs *map[string]DebPackage) bool {
	_, exists := (*pkgs)[p.Package]
	return exists
}

// generatePackageIndex generates the package index from the map of packages
func generatePackageIndex(pkgs *map[string]DebPackage) []byte {
	var buffer bytes.Buffer
	for _, p := range *pkgs {
		p.CalculateAllHashes()
		p.CalcSize()

		t := reflect.TypeOf(p)

		for i := 0; i < t.NumField(); i++ {
			fieldValueStruct := reflect.ValueOf(p).Field(i)
			if hasValue(fieldValueStruct.Interface()) {
				var value interface{}
				if fieldValueStruct.Kind() == reflect.Slice && fieldValueStruct.Type().Elem().Kind() == reflect.String {
					value = strings.Join(fieldValueStruct.Interface().([]string), ",")
				} else {
					value = fieldValueStruct.Interface()
				}
				buffer.WriteString(p.GetAptTag(fieldValueStruct.Interface()))
				buffer.WriteString(": ")
				buffer.WriteString(fmt.Sprintf("%v", value))
				buffer.WriteString("\n")
			}
		}
	}
	return buffer.Bytes()
}

// hasValue checks if the attribute has a value
func hasValue(attr interface{}) bool {
	// Get the reflected value of the attribute
	value := reflect.ValueOf(attr)

	// Check if the value is zero (has no value)
	switch value.Kind() {
	case reflect.Invalid:
		return false // nil or invalid value
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return !value.IsNil() // Check if it is nil
	case reflect.String:
		return value.String() != "" // Check if the string is empty
	case reflect.Array, reflect.Struct:
		return true // Arrays and structs always have value
	default:
		// For basic types (int, float, bool, etc.), check if it is the zero value
		return value.Interface() != reflect.Zero(value.Type()).Interface()
	}
}
