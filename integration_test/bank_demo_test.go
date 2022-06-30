package integration_test

import (
	"fmt"
	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
	"testing"
)

func TestBankDemo(t *testing.T) {
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
	client := opb.NewClient(cfg, &authToken)

	// 导入私钥
	address, _ := client.Key.Recover("validator", "12345678", "opera vivid pride shallow brick crew found resist decade neck expect apple chalk belt sick author know try tank detail tree impact hand best")
	fmt.Println("address:", address)

	// 查询账户
	account, err := client.Bank.QueryAccount(address)
	if err != nil {
		fmt.Println(fmt.Errorf("BANK 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("bank 查询成功:", account)
	}

	// 发生交易
	// 初始化 Tx 基础参数
	baseTx := types.BaseTx{
		From:     "validator", // 对应上面导入的私钥名称
		Password: "12345678",  // 对应上面导入的私钥密码
		Gas:      200000,      // 单 Tx 消耗的 Gas 上限
		Memo:     "",          // Tx 备注
		Mode:     types.Sync,  // Tx 广播模式
	}
	amount, _ := types.ParseDecCoins("10upoint")
	result, err := client.Bank.Send("toAddress", amount, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("BANK 发送失败: %s", err.Error()))
		return
	} else {
		fmt.Println("BANK 发送成功：", result.Hash)
	}
}
