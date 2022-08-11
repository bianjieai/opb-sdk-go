package test2

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/core-sdk-go"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
)

var iritaClient2 sdk.Client
var addr = ""

const (
	nodeURI  = "http://localhost:26657"
	grpcAddr = "localhost:9090"
	chainID  = "irita-tesnet"
)

var (
	password = "12345678"

	//algo = "sm2"
	//name             = "validator"
	//mnemonic         = "second lobster moment enjoy nasty sight remember cram pave raise father tunnel sort soon carbon excite domain foil design approve boost keep globe sheriff"

	algo     = "secp256k1"
	name     = "test03"
	mnemonic = "shy spawn around wheat target nose kick body letter october banner slide toward trade fog moment life cabbage napkin camera couch range choose vacant"
)

func init() {
	fee, _ := types.ParseDecCoins("1000ugas")
	options := []types.Option{
		types.AlgoOption(algo),
		types.KeyDAOOption(store.NewMemory(nil)),
		types.FeeOption(fee),
		types.TimeoutOption(20),
		types.CachedOption(true),
	}
	config, err := types.NewClientConfig(nodeURI, grpcAddr, chainID, options...)
	if err != nil {
		panic(err)
	}

	iritaClient2 = sdk.NewClient(config)
	add, err := iritaClient2.Key.Recover(name, password, mnemonic)
	if err != nil {
		panic(err)
	}
	addr = add

	//providerAddr, err := iritaClient2.Key.Recover("provider", "12345678", "broccoli bone broom cliff pluck normal jelly guess brush urban depth under become sustain chronic away milk repair prevent latin good rug refuse ill")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("provider:", providerAddr)
	fmt.Println("初始化完成")
}

var baseTx2 = types.BaseTx{
	From:               name,
	Password:           password,
	Mode:               types.Commit,
	Gas:                200000,
	SimulateAndExecute: false,
}

func Test_simBankTx(t *testing.T) {
	coins, err := types.ParseDecCoins("1000ugas")
	if err != nil {
		t.Error(err)
		return
	}
	send, err := iritaClient2.Bank.Send("iaa17y3qs2zuanr93nk844x0t7e6ktchwygnc8fr0g", coins, baseTx2)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("交易信息：", send)
	fmt.Println("交易信息：", string(send.Data))

	//time.Sleep(time.Second * 1)
	//queryTx, err := iritaClient2.QueryTx(send.Hash)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//fmt.Println("交易哈希：", send.Hash)
	//fmt.Println("交易信息：", queryTx.Tx)

	fmt.Println("结束，地址：", addr)
}

func Test_queryTx(t *testing.T) {
	queryTx, err := iritaClient2.QueryTx("7785201D39F2D682189FA1B882EB22CC17D2A934225069EF517A28BBC08F4D4A")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(queryTx)
}

func Test_queryAccount(t *testing.T) {
	account, err := iritaClient2.Bank.QueryAccount("iaa18q8k3mkp97zmhhycugwqwgraw3cgsw2gf06var")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(account)
}
