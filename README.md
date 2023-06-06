# OPB SDK GO

IRITA 开放联盟链 SDK（Golang）

## 快速开始

### 引入依赖

编辑 go.mod

```
replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
	github.com/prometheus/common => github.com/prometheus/common v0.26.0
)

require (
	github.com/irisnet/core-sdk-go v0.0.0-20220720085949-4d825adb8054
	github.com/irisnet/irismod-sdk-go v0.0.0-20220825063058-6bf8e60b42c8
	github.com/bianjieai/iritamod-sdk-go v0.0.0-20220708032705-9e8e301da3a8
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
	"time"

	"github.com/irisnet/irismod-sdk-go/mt"
	"github.com/irisnet/irismod-sdk-go/nft"

	"github.com/irisnet/irismod-sdk-go/record"

	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
	tendermintTypes "github.com/tendermint/tendermint/abci/types"
)

func main() {
	//能量值费用：「创建NFT/MT类别」、「发行NFT/MT」为 40 万能量值（等值人民币：0.1元），其他交易类型为 20 万能量值（等值人民币：0.05元）
	fee, _ := types.ParseDecCoins("400000ugas") 
	// 初始化 SDK 配置
	options := []types.Option{
		types.AlgoOption(algo),
		types.KeyDAOOption(store.NewMemory(nil)),
		types.FeeOption(fee),
		types.TimeoutOption(10),
		types.CachedOption(true),
		types.WSAddrOption(wsAddress),
	}
	cfg, err := types.NewClientConfig(rpcAddress, grpcAddress, chainID, options...)
	if err != nil {
		panic(err)
	}

	// 初始化 OPB 网关账号（测试网环境设置为 nil 即可）
	authToken := model.NewAuthToken(projectId, projectKey, chainAccountAddr)

	// 开启 TLS 连接
	// 若服务器要求使用安全链接，此处应设为true；若此处设为false可能导致请求出现长时间不响应的情况
	// 若开启 TLS 连接，则必须设置验证证书的主机名
	if err :=  authToken.SetRequireTransportSecurity(false, ""); err != nil {
		fmt.Println(fmt.Errorf("开启TLS失败: %s", err.Error()))
		return
    }
	// 创建 OPB 客户端
	client := opb.NewClient(cfg, &authToken)

	// 导入私钥
	address, err := client.Key.Recover(name, password, mnemonic)
	if err != nil {
		fmt.Println(fmt.Errorf("导入私钥失败: %s", err.Error()))
		return
	}
	fmt.Println("address:", address)

	// 初始化 Tx 基础参数
	baseTx := types.BaseTx{
		From:     name,       // 对应上面导入的私钥名称
		Password: password,   // 对应上面导入的私钥密码
		Gas:      400000,     // 单 Tx 消耗的 Gas 上限
		Memo:     "",         // Tx 备注
		Mode:     types.Sync, // Tx 广播模式
	}
	// 初始化交易哈希查询队列
	var hashArray []string

	// 使用 Client 选择对应的功能模块，查询链上状态；例：查询账户信息
	acc, err := client.Bank.QueryAccount(address)
	if err != nil {
		fmt.Println(fmt.Errorf("账户查询失败: %s", err.Error()))
	} else {
		fmt.Println("账户信息查询成功：", acc)
	}

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：创建 NFT 类别
	nftResult, err := client.NFT.IssueDenom(nft.IssueDenomRequest{ID: "testdenom", Name: "TestDenom", Schema: "{}"}, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("NFT 类别创建失败: %s", err.Error()))
	} else {
		fmt.Println("NFT 类别创建成功 TxHash：", nftResult.Hash)
		hashArray = append(hashArray, nftResult.Hash)
	}

	// 例：创建 NFT
	mintNFT, err := client.NFT.MintNFT(nft.MintNFTRequest{Denom: "testdenom", ID: "testnft1", Name: "aaa", URI: "www.test.com", Data: "test", Recipient: address}, baseTx)
	if err != nil {
		e := err.(types.Error)
		if e.Codespace() == nft.ErrInvalidTokenID.Codespace() {
			fmt.Println("Err code: ", e.Code())
		}
		fmt.Println(fmt.Errorf("NFT 创建失败: %s", err))
	} else {
		fmt.Println("NFT 创建成功 TxHash：", mintNFT.Hash)
		hashArray = append(hashArray, mintNFT.Hash)
	}

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：创建 MT 类别
	mtResult, err := client.MT.IssueDenom(mt.IssueDenomRequest{Name: "TestDenom", Data: []byte("TestData")}, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("MT 类别创建失败: %s", err.Error()))
	} else {
		fmt.Println("MT 类别创建成功 TxHash：", mtResult.Hash)
		hashArray = append(hashArray, mtResult.Hash)
	}

	// 例：增发 MT
	addMT, err := client.MT.AddMT(mt.AddMTRequest{ID: "c54e89be44edfd421678d4a504f6c5f110878f52883d19935fb412107168015f", DenomID: "a6a8dabe077c23054a582f8ff9847e52f95385c342aa80b8b662eeb5b8f24b19", Amount: 100}, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("MT 增发失败: %s", err.Error()))
	} else {
		fmt.Println("MT 增发成功 TxHash：", addMT.Hash)
		hashArray = append(hashArray, addMT.Hash)
	}

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：BANK 发送交易
	result, err := client.Bank.Send("", fee, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("BANK 发送失败: %s", err.Error()))
	} else {
		fmt.Println("BANK 发送成功：", result.Hash)
		hashArray = append(hashArray, result.Hash)
	}

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：创建存证
	req := record.CreateRecordRequest{
		Contents: []record.Content{
			{
				Digest:     "digest", //存证元数据摘要
				DigestAlgo: "sha256", //存证元数据摘要的生成算法
				URI:        "www.baidu.com",
				Meta:       "tx", //源数据
			},
		},
	}
	recordResult, err := client.Record.CreateRecord(req, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("存证创建失败: %s", err.Error()))
	} else {
		fmt.Println("存证发送成功：", recordResult.Hash)
		hashArray = append(hashArray, recordResult.Hash)
	}

	// 等待十秒后查询交易
	time.Sleep(time.Second * 10)
	for _, hash := range hashArray {
		tx, err := client.QueryTx(hash)
		if err != nil {
			fmt.Println("查询交易错误：", err)
			continue
		}
		if tx.Result.Code == tendermintTypes.CodeTypeOK {
			fmt.Println("交易上链成功，交易哈希:", hash)
		} else {
			fmt.Printf("交易上链失败，交易哈希:%s， 错误:%s. \n", hash, tx.Result.Log)
		}
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

// 主网使用的配置
//var (
//	wsAddress   = fmt.Sprintf("%s/api/%s/ws", "wss://opbningxia.bsngate.com:18602", projectId)
//	rpcAddress  = fmt.Sprintf("%s/api/%s/rpc", "https://opbningxia.bsngate.com:18602", projectId)
//	grpcAddress = "opbningxia.bsngate.com:18603"
//	chainID     = "wenchangchain"
//
//	algo             = ""
//	projectId        = ""
//	projectKey       = ""
//	chainAccountAddr = ""
//	name             = ""
//	password         = ""
//	mnemonic         = ""
//)

// 测试链使用的配置
var (
	wsAddress   = ""
	rpcAddress  = "http://testnet.bianjie.ai:26657"
	grpcAddress = "testnet.bianjie.ai:9090"
	chainID     = "testing"
	tlsServiceName = "grpcs.testnet.bianjie.ai"

	algo             = "sm2"
	projectId        = "TestProjectID"
	projectKey       = "TestProjectKey"
	chainAccountAddr = "TestChainAccountAddress"
	name             = "test_key_name"
	password         = "test_password"
	mnemonic         = "supreme zero ladder chaos blur lake dinner warm rely voyage scan dilemma future spin victory glance legend faculty join man mansion water mansion exotic"
)

```
