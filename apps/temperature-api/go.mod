module temperature-api

go 1.25.0

require (
	github.com/gofiber/contrib/v3/swaggerui v1.0.11
	github.com/gofiber/fiber/v3 v3.5.0
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/go-openapi/runtime/server-middleware v0.33.1 // indirect
	github.com/gofiber/schema v1.8.4 // indirect
	github.com/gofiber/utils/v2 v2.4.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.73.0 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// golangci-lint 2.6 собран на Go 1.25 и падает на //go:build go1.26
// в свежих x/crypto и x/net (fiber тянет их транзитивно).
replace golang.org/x/crypto => golang.org/x/crypto v0.45.0

replace golang.org/x/net => golang.org/x/net v0.43.0
