module example.com/m/v2

go 1.18

require (
	github.com/akakou/ecdaa v0.0.0
	github.com/google/go-tpm v0.3.3
	github.com/mattn/go-sqlite3 v1.14.15
)

require miracl v0.0.0

require (
	golang.org/x/mobile v0.0.0-20230427221453-e8d11dd0ba41 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/tools v0.1.12 // indirect
)

replace github.com/akakou/ecdaa => ../thirdparty/ecdaa32

replace miracl => ../thirdparty/ecdaa32/thirdparty/miracl

replace github.com/google/go-tpm => ../thirdparty/ecdaa32/thirdparty/go-tpm
