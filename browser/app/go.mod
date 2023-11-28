module example.com/m/v2

go 1.18

require (
	github.com/akakou/ecdaa v0.0.0
	github.com/akakou/scrappy v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
)

require (
	github.com/akakou/mcl_utils v0.0.0 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	golang.org/x/sys v0.5.0 // indirect
	miracl v0.0.0 // indirect
)

replace github.com/akakou/ecdaa => ../../thirdparty/ecdaa64

replace miracl => ../../thirdparty/ecdaa64/thirdparty/miracl

replace github.com/akakou/scrappy => ../../core

replace github.com/akakou/mcl_utils => ../../thirdparty/ecdaa64/mcl_utils

replace github.com/google/go-tpm => ../../thirdparty/ecdaa64/thirdparty/go-tpm
