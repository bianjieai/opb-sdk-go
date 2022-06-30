module github.com/bianjieai/opb-sdk-go

go 1.16

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/prometheus/common => github.com/prometheus/common v0.26.0
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
	github.com/tharsis/ethermint v0.8.1 => github.com/bianjieai/ethermint v0.8.2-0.20220211020007-9ec25dde74d4
)

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/bianjieai/iritamod-sdk-go v0.0.0-20220622091247-de18248d9580
	github.com/cosmos/cosmos-sdk v0.45.5 // indirect
	github.com/irisnet/core-sdk-go v0.0.0-20220515104139-554292f91a1a
	github.com/irisnet/irismod-sdk-go/mt v0.0.0-20220630035504-d1da9a44093b
	github.com/irisnet/irismod-sdk-go/nft v0.0.0-20220630071928-f8a9c5d21264
	github.com/irisnet/irismod-sdk-go/random v0.0.0-20220630035504-d1da9a44093b
	github.com/irisnet/irismod-sdk-go/record v0.0.0-20220630035504-d1da9a44093b
	github.com/irisnet/irismod-sdk-go/service v0.0.0-20220630035504-d1da9a44093b
	github.com/irisnet/irismod-sdk-go/token v0.0.0-20220630035504-d1da9a44093b
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tendermint v0.34.19
	google.golang.org/grpc v1.45.0
)
