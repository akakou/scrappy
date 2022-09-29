module example.com/m/v2

go 1.18

require github.com/akakou/ecdaa v0.0.0

require miracl v0.0.0 // indirect

require (
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sys v0.0.0-20210629170331-7dc0b73dc9fb // indirect
)

replace github.com/akakou/ecdaa => ../thirdparty/ecdaa

replace miracl => ../thirdparty/ecdaa/thirdparty/miracl

replace github.com/google/go-tpm => ../thirdparty/ecdaa/thirdparty/go-tpm
