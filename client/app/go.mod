module example.com/m/v2

go 1.18

require github.com/akakou/ecdaa v0.0.0

require miracl v0.0.0 // indirect

require core v0.0.0

require (
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/gin-contrib/sessions v0.0.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.0.0-20220926161630-eccd6366d1be // indirect
	golang.org/x/net v0.0.0-20220927171203-f486391704dc // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/akakou/ecdaa => ../../thirdparty/ecdaa

replace miracl => ../../thirdparty/ecdaa/thirdparty/miracl

replace core => ../../core

replace github.com/google/go-tpm => ../../thirdparty/ecdaa/thirdparty/go-tpm
