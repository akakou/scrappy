module example.com/m/v2

go 1.18

require github.com/akakou/ecdaa v0.0.0 // indirect

require core v0.0.0

require miracl v0.0.0 // indirect

require (
	github.com/google/go-tpm v0.3.3 // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
)

replace github.com/akakou/ecdaa => ../thirdparty/ecdaa

replace miracl => ../thirdparty/ecdaa/thirdparty/miracl

replace core => ../core

replace github.com/google/go-tpm => ../thirdparty/ecdaa/thirdparty/go-tpm
