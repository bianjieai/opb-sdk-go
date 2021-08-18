package main

import (
	"fmt"
	"github.com/bianjieai/irita-sdk-go/types"
	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
)

func main() {

	// 初始化 SDK 配置
	cfg, err := types.NewClientConfig("localhost:26657", "localhost:9090", "testing")
	if err != nil {
		panic(err)
	}

	// 初始化 OPB 网关账号
	authToken := model.NewAuthToken("TestProjectID", "TestProjectKey", "TestChainAccountAddress")

	// 创建 OPB 客户端
	client := opb.NewClient(cfg, &authToken)

	// 导入私钥
	client.Key.Recover("test_key_name", "test_password", "supreme zero ladder chaos blur lake dinner warm rely voyage scan dilemma future spin victory glance legend faculty join man mansion water mansion exotic")

	// 初始化 Tx 基础参数
	//baseTx := types.BaseTx{
	//	From:     "test_key_name", // 对应上面导入的私钥名称
	//	Password: "test_password", // 对应上面导入的私钥密码
	//	Gas:      20000,		   // 单 Tx 消耗的 Gas 上限
	//	Memo:     "",			   // Tx 备注
	//	Mode:     types.Commit,    // Tx 广播模式
	//}

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：发行 NFT 类别
	// result, err := client.NFT.IssueDenom(nft.IssueDenomRequest{}, baseTx)

	// 使用 Client 选择对应的功能模块，查询链上状态；例：查询账户信息
	acc, _ := client.Bank.QueryAccount("iaa1lxvmp9h0v0dhzetmhstrmw3ecpplp5tljnr35f")
	fmt.Println("账户信息：", acc)
}
