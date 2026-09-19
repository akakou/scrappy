module example.com/m/v2

go 1.21.4

require (
	github.com/akakou/ecdaa v0.0.2
	github.com/akakou/scrappy v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
)

require (
	github.com/akakou-fork/amcl-go/miracl v0.0.0-20240206094909-344c847a50cc // indirect
	github.com/akakou/fp256bn-amcl-utils v0.0.2 // indirect
	github.com/google/go-tpm v0.9.1-0.20240206213016-638c2b803c16 // indirect
	golang.org/x/sys v0.16.0 // indirect
// indirect
)

replace github.com/akakou/scrappy => ../../core
