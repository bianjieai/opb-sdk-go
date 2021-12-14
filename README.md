# OPB SDK GO

IRITA 开放联盟链 SDK（Golang）

## 快速开始

### 引入依赖

编辑 go.mod

```
replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)

require (
	github.com/bianjieai/irita-sdk-go v1.1.1-0.20211214032850-7c9cd100e6bd
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.40.0
)
```

### 创建和使用 OPB 客户端

参考 [示例代码](./examples/main.go)

```go
package main

import (
	"fmt"
	"github.com/bianjieai/irita-sdk-go/modules/nft"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/store"
	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"time"
)

func main() {
	fee, _ := types.ParseCoin("100000uirita")  // 设置文昌链主网的默认费用，10W不够就填20W，30W....
	// 初始化 SDK 配置
	options := []types.Option{
		types.KeyDAOOption(store.NewMemory(nil)),
		types.FeeOption(types.NewDecCoinsFromCoins(fee)),
	}
	cfg, err := types.NewClientConfig("https://opbt.bsngate.com:18602/api/b26b8b53e1b74aa6974ab207e25310f0/rpc", "wss://opbt.bsngate.com:18602/api/b26b8b53e1b74aa6974ab207e25310f0/ws", "opbt.bsngate.com:18603", "wenchangchain", options...)
	if err != nil {
		panic(err)
	}

	// 初始化 OPB 网关账号（测试网环境设置为 nil 即可）
	authToken := model.NewAuthToken("TestProjectID", "TestProjectKey", "TestChainAccountAddress")

	// 创建 OPB 客户端
	client := opb.NewClient(cfg, &authToken)

	// 导入私钥
	client.Key.Recover("test_key_name", "test_password", "supreme zero ladder chaos blur lake dinner warm rely voyage scan dilemma future spin victory glance legend faculty join man mansion water mansion exotic")

	// 初始化 Tx 基础参数
	baseTx := types.BaseTx{
		From:     "test_key_name", // 对应上面导入的私钥名称
		Password: "test_password", // 对应上面导入的私钥密码
		Gas:      200000,		   // 单 Tx 消耗的 Gas 上限
		Memo:     "",			   // Tx 备注
		Mode:     types.Commit,    // Tx 广播模式
	}

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：创建 NFT 类别
	result, err := client.NFT.IssueDenom(nft.IssueDenomRequest{ID: "testdenom", Name: "TestDenom", Schema: "{}"}, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("NFT 类别创建失败: %s", err.Error()))
	} else {
		fmt.Println("NFT 类别创建成功：", result.Hash)
	}

	// 使用 Client 选择对应的功能模块，查询链上状态；例：查询账户信息
	acc, err := client.Bank.QueryAccount("iaa1lxvmp9h0v0dhzetmhstrmw3ecpplp5tljnr35f")
	if err != nil {
		fmt.Println(fmt.Errorf("账户查询失败: %s", err.Error()))
	} else {
		fmt.Println("账户信息查询成功：", acc)
	}

	// 使用 Client 订阅事件通知，例：订阅区块
	subs, err := client.SubscribeNewBlock(types.NewEventQueryBuilder(), func(block types.EventDataNewBlock) {
		fmt.Println(block)
	})

	if err != nil {
		fmt.Println(fmt.Errorf("区块订阅失败: %s", err.Error()))
	} else {
		fmt.Println("区块订阅成功：", subs.ID)
	}
	time.Sleep(time.Second * 20)
}
```