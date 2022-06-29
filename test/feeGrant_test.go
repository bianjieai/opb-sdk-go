package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/client"
	"github.com/irisnet/core-sdk-go/feegrant"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/stretchr/testify/require"
)

var cfg types.ClientConfig
var txClient client.Client
var baseTx types.BaseTx

func init() {
	fee, _ := types.ParseDecCoins("300000ugas") // 设置文昌链主网的默认费用，10W不够就填20W，30W....
	// 初始化 SDK 配置
	options := []types.Option{
		types.AlgoOption("sm2"),
		types.KeyDAOOption(store.NewMemory(nil)),
		types.FeeOption(fee),
		types.TimeoutOption(10),
		types.CachedOption(true),
	}
	cfg, err := types.NewClientConfig("http://127.0.0.1:26657", "127.0.0.1:9090", "2022", options...)
	if err != nil {
		panic(err)
	}
	// 初始化 OPB 网关账号（测试网环境设置为 nil 即可）
	//authToken := model.NewAuthToken("TestProjectID", "TestProjectKey", "TestChainAccountAddress")

	// 创建 OPB 客户端
	//client := opb.NewClient(cfg, &authToken)
	txClient = client.NewClient(cfg)

	// 开启 TLS 连接
	// 若服务器要求使用安全链接，此处应设为true；若此处设为false可能导致请求出现长时间不响应的情况
	//authToken.SetRequireTransportSecurity(true)
	// 创建 OPB 客户端
	//feeGrantClient = sdk.NewClient(cfg)

	// 导入私钥
	address, _ := txClient.Key.Recover("validator", "12345678", "west farm disease weasel age cram cross second battle brief slim steel network arrive series lab gorilla gun fiction robust skin torch planet burden")
	fmt.Println("address:", address)

	// 初始化 Tx 基础参数
	baseTx = types.BaseTx{
		From:     "validator", // 对应上面导入的私钥名称
		Password: "12345678",  // 对应上面导入的私钥密码
		Gas:      200000,      // 单 Tx 消耗的 Gas 上限
		Memo:     "",          // Tx 备注
		Mode:     types.Sync,  // Tx 广播模式
	}
}

//授权
func TestGrantAllowance(t *testing.T) {
	granter, _ := types.AccAddressFromBech32("iaa193eqcr7zwtfjx7us0wm33ddtdtct38adr942f3")
	grantee, _ := types.AccAddressFromBech32("iaa1l3r7kx3nqa7uymk225s97dfsw46cysf3hqwdj8")
	atom := types.NewCoins(types.NewInt64Coin("ugas", 55500000000))
	threeHours := time.Now().Add(3 * time.Hour)
	basic := &feegrant.BasicAllowance{
		SpendLimit: atom,        //授权额度
		Expiration: &threeHours, //过期时间
	}
	result, err := txClient.Feegrant.GrantAllowance(granter, grantee, basic, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}

//设置交易代扣
func TestFeeGrant(t *testing.T) {
	address, _ := txClient.Key.Recover("account4", "12345678", "such tooth bicycle bonus album west win chunk tuna erosion protect rifle kiss purity marble ketchup spirit material cash fee argue silent column obscure")
	fmt.Println("address:", address)
	feeGranter, _ := types.AccAddressFromBech32("iaa193eqcr7zwtfjx7us0wm33ddtdtct38adr942f3")
	baseTx2 := types.BaseTx{
		From:       "account4",
		Password:   "12345678",
		Gas:        200000,
		Memo:       "",
		Mode:       types.Sync,
		FeeGranter: feeGranter, //设置代扣地址
	}
	nftResult, err := txClient.NFT.IssueDenom(nft.IssueDenomRequest{ID: "testdenom112", Name: "TestDenom112", Schema: "{}"}, baseTx2)
	require.NoError(t, err)
	require.NotEmpty(t, nftResult.Hash)
}

//解除授权
func TestRevokeAllowance(t *testing.T) {
	granter, _ := types.AccAddressFromBech32("iaa193eqcr7zwtfjx7us0wm33ddtdtct38adr942f3")
	grantee, _ := types.AccAddressFromBech32("iaa1l3r7kx3nqa7uymk225s97dfsw46cysf3hqwdj8")
	result, err := txClient.Feegrant.RevokeAllowance(granter, grantee, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}
