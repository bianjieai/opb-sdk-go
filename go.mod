module github.com/bianjieai/opb-sdk-go

go 1.16

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/bianjieai/iritamod-sdk-go v0.0.0-20230526172348-b7e022cb4f01
	github.com/cosmos/cosmos-sdk v0.45.5 // indirect
	github.com/irisnet/core-sdk-go v0.1.1-0.20230320094234-f12edcab435e
	github.com/irisnet/irismod-sdk-go/mt v0.0.0-20221014104619-6f27c71cd5e4
	github.com/irisnet/irismod-sdk-go/nft v0.0.0-20221014104619-6f27c71cd5e4
	github.com/irisnet/irismod-sdk-go/random v0.0.0-20221014104619-6f27c71cd5e4
	github.com/irisnet/irismod-sdk-go/record v0.0.0-20230321015746-6d76d0435f7c
	github.com/irisnet/irismod-sdk-go/service v0.0.0-20221014104619-6f27c71cd5e4
	github.com/irisnet/irismod-sdk-go/token v0.0.0-20221014104619-6f27c71cd5e4
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tendermint v0.34.19
	google.golang.org/grpc v1.45.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/prometheus/common => github.com/prometheus/common v0.26.0
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)
