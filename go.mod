module github.com/bianjieai/opb-sdk-go

go 1.16

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)

require (
	github.com/irisnet/core-sdk-go v0.0.0-20220515104139-554292f91a1a
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.41.0
)
