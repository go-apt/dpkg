# dpkg

`dpkg` is a Go package for managing Debian packages. It provides functionalities to read and parse `.deb` files, list installed packages, and filter packages based on their names.

## Installation

To install the `dpkg` package, use the following command:

```sh
go get github.com/go-apt/dpkg
```

## Main Functions

### Reading the Contents of a `.deb` File

To read the contents of a `.deb` file and retrieve its metadata, use the `Info` function:

```go
// Create a new instance of the Dpkg struct
d := dpkg.NewDpkg()

// Read the contents of the .deb file
pkg, err := d.Info("/path/to/debFile")
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}

// Print package name from the .deb file
fmt.Printf("Package from .deb file: %s\n", pkg.Fields["Package"])
```

### Validating a `.deb` File

To validate if a file is a valid `.deb` package, use the `IsDebFile` function:

```go
debFile := "/path/to/debFile"

// Validate if the file is a .deb package
if !d.IsDebFile(debFile) {
    fmt.Fprintf(os.Stderr, "Error: %s is not a valid .deb file\n", debFile)
    os.Exit(1)
}

fmt.Printf("%s is a valid .deb file\n", debFile)
```

### Listing Installed Packages

To list the installed packages and retrieve their metadata, use the `List` function:

```go
// Create a new instance of the Dpkg struct
d := dpkg.NewDpkg()

// Read the contents of the /var/lib/dpkg/status file
packages, err := d.List()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}

// Print package names from installed packages
for _, p := range packages {
    fmt.Printf("Package from dpkg status file: %s\n", p.Fields["Package"])
}
```

### Filtering Packages by Name

To filter packages by name, use the `ListGrep` function:

```go
// Create a new instance of the Dpkg struct
d := dpkg.NewDpkg()

// Filter packages by name
filteredPackages, err := d.ListGrep("apt")
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}

// Print package names from filtered packages
for _, p := range filteredPackages {
    fmt.Printf("Packages from Grep search: %s\n", p.Fields["Package"])
}
```
