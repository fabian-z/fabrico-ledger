module github.com/fabian-z/fabrico-ledger

go 1.19

replace github.com/SmartBFT-Go/consensus/v2 => github.com/fabian-z/consensus/v2 v2.0.0-20221128114342-11dd18c41432

require (
	github.com/SmartBFT-Go/consensus/v2 v2.3.0
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/miekg/dns v1.1.43 // indirect
)

require (
	filippo.io/edwards25519 v1.0.1-0.20220803165937-8c58ed0e3550 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/libp2p/zeroconf/v2 v2.2.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
