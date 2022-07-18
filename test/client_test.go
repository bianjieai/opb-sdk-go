package test

import (
	"fmt"
	"github.com/avast/retry-go"
	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/client"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/query"
	"github.com/irisnet/core-sdk-go/types/store"
	tendermintTypes "github.com/tendermint/tendermint/abci/types"
	"time"
)

var txClient client.Client
var baseTx types.BaseTx
var address string

// 分页查询参数
var pagination = &query.PageRequest{
	//Key: []byte{1}, // 下标，与Offset二选一
	Offset:     0,     // 偏移量，与Key二选一
	Limit:      0,     // 查询数量：最大值100
	CountTotal: false, // 是否查询总数：目前仅支持false
}

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
	cfg, err := types.NewClientConfig("http://47.100.192.234:26657", "47.100.192.234:9090", "testing", options...)
	if err != nil {
		panic(err)
	}
	// 初始化 OPB 网关账号（测试网环境设置为 nil 即可）
	authToken := model.NewAuthToken("TestProjectID", "TestProjectKey", "TestChainAccountAddress")
	// 开启 TLS 连接
	// 若服务器要求使用安全链接，此处应设为true；若此处设为false可能导致请求出现长时间不响应的情况
	authToken.SetRequireTransportSecurity(false)

	// 创建 OPB 客户端
	txClient = opb.NewClient(cfg, &authToken)

	// 导入私钥
	address, _ = txClient.Key.Recover("test_key_name", "test_password", "supreme zero ladder chaos blur lake dinner warm rely voyage scan dilemma future spin victory glance legend faculty join man mansion water mansion exotic")

	// 初始化 Tx 基础参数
	baseTx = types.BaseTx{
		From:     "test_key_name", // 对应上面导入的私钥名称
		Password: "test_password", // 对应上面导入的私钥密码
		Gas:      200000,          // 单 Tx 消耗的 Gas 上限
		Memo:     "",              // Tx 备注
		Mode:     types.Commit,    // Tx 广播模式
	}
}

// 异步模式上链
func syncTx(txHash string) error {
	err := retry.Do(func() error {
		tx, err := txClient.QueryTx(txHash)
		if err != nil {
			return err
		}
		if tx.Result.Code == tendermintTypes.CodeTypeOK {
			fmt.Println("交易上链成功，交易哈希:", txHash)
		} else {
			fmt.Printf("交易上链失败，交易哈希:%s， 错误:%s. \n", txHash, tx.Result.Log)
		}
		return nil
	}, retry.Attempts(3), retry.Delay(2*time.Second))
	return err
}
