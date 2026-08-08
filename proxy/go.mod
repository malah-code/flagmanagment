module github.com/flagmanagment/proxy

go 1.23

require (
	github.com/go-chi/chi/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/rs/zerolog v1.35.1
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/flagmanagment/backend => ../backend
