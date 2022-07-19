package client

import (
	"github.com/bianjieai/iritamod-sdk-go/identity"
	"github.com/bianjieai/iritamod-sdk-go/perm"
	"github.com/irisnet/core-sdk-go/bank"
	"github.com/irisnet/core-sdk-go/client"
	keys "github.com/irisnet/core-sdk-go/client"
	"github.com/irisnet/core-sdk-go/common/codec"
	cdctypes "github.com/irisnet/core-sdk-go/common/codec/types"
	cryptocodec "github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/irisnet/core-sdk-go/types"
	txtypes "github.com/irisnet/core-sdk-go/types/tx"
	"github.com/irisnet/irismod-sdk-go/mt"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/irisnet/irismod-sdk-go/random"
	"github.com/irisnet/irismod-sdk-go/record"
	"github.com/irisnet/irismod-sdk-go/service"
	"github.com/irisnet/irismod-sdk-go/token"
	"github.com/tendermint/tendermint/libs/log"
)

type Client struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig

	types.BaseClient
	Bank    bank.Client
	Token   token.Client
	Service service.Client
	Record  record.Client
	Random  random.Client
	NFT     nft.Client
	MT      mt.Client

	Identity identity.Client
	Perm     perm.Client

	Key keys.Client
}

func NewClient(cfg types.ClientConfig) Client {
	encodingConfig := makeEncodingConfig()

	// create a instance of baseClient
	baseClient := client.NewBaseClient(cfg, encodingConfig, nil)
	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	tokenClient := token.NewClient(baseClient, encodingConfig.Marshaler)
	serviceClient := service.NewClient(baseClient, encodingConfig.Marshaler)
	recordClient := record.NewClient(baseClient, encodingConfig.Marshaler)
	nftClient := nft.NewClient(baseClient, encodingConfig.Marshaler)
	randomClient := random.NewClient(baseClient, encodingConfig.Marshaler)
	mtClient := mt.NewClient(baseClient, encodingConfig.Marshaler)

	idClient := identity.NewClient(baseClient, encodingConfig.Marshaler)
	permClient := perm.NewClient(baseClient, encodingConfig.Marshaler)

	keysClient := keys.NewKeysClient(cfg, baseClient)

	sdkClient := &Client{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
		Bank:           bankClient,
		Token:          tokenClient,
		Service:        serviceClient,
		Record:         recordClient,
		Random:         randomClient,
		NFT:            nftClient,
		MT:             mtClient,
		Identity:       idClient,
		Perm:           permClient,
		Key:            keysClient,
	}

	sdkClient.RegisterModule(
		bankClient,
		tokenClient,
		serviceClient,
		recordClient,
		nftClient,
		randomClient,
		mtClient,
		idClient,
		permClient,
	)
	return *sdkClient
}

func (client Client) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client Client) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client Client) AppCodec() codec.Marshaler {
	return client.encodingConfig.Marshaler
}

func (client Client) EncodingConfig() types.EncodingConfig {
	return client.encodingConfig
}

func (client Client) Manager() types.BaseClient {
	return client.BaseClient
}

func (client Client) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
	}
}

func (client Client) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := txtypes.NewTxConfig(marshaler, txtypes.DefaultSignModes)

	encodingConfig := types.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	RegisterLegacyAminoCodec(encodingConfig.Amino)
	RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	cryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.v1beta1.Msg", (*types.Msg)(nil))
	txtypes.RegisterInterfaces(registry)
	cryptocodec.RegisterInterfaces(registry)
}
