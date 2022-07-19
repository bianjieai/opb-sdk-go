package main

import (
	"fmt"
	"github.com/irisnet/irismod-sdk-go/record"
	"time"

	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
	"github.com/irisnet/irismod-sdk-go/mt"
	"github.com/irisnet/irismod-sdk-go/nft"
	tendermintTypes "github.com/tendermint/tendermint/abci/types"
)

func main() {
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
	address, _ := client.Key.Recover("test_key_name", "test_password", "supreme zero ladder chaos blur lake dinner warm rely voyage scan dilemma future spin victory glance legend faculty join man mansion water mansion exotic")
	fmt.Println("address:", address)

	// 初始化 Tx 基础参数
	baseTx := types.BaseTx{
		From:     "test_key_name", // 对应上面导入的私钥名称
		Password: "test_password", // 对应上面导入的私钥密码
		Gas:      200000,          // 单 Tx 消耗的 Gas 上限
		Memo:     "",              // Tx 备注
		Mode:     types.Sync,      // Tx 广播模式
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

	// 创建 NFT
	mintNFT, err := client.NFT.MintNFT(nft.MintNFTRequest{Denom: "testdenom", ID: "OpbTestName_1", Name: "aaa", URI: "www.baidu.com", Data: "test", Recipient: address}, baseTx)
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

	// 使用 Client 选择对应的功能模块，构造、签名并发送交易；例：BANK 发送交易
	result, err := client.Bank.Send("", fee, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("BANK 发送失败: %s", err.Error()))
	} else {
		fmt.Println("BANK 发送成功：", result.Hash)
		hashArray = append(hashArray, result.Hash)
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

	//创建存证
	req := record.CreateRecordRequest{
		Contents: []record.Content{
			{
				Digest:     "digest", //存证元数据摘要
				DigestAlgo: "sha256", //存证元数据摘要的生成算法
				URI:        "1231321",
				Meta:       "tx", //源数据
			},
		},
	}
	recordResp, err := client.Record.CreateRecord(req, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("存证创建失败: %s", err.Error()))
	}
	fmt.Println(recordResp.RecordId)
	fmt.Println(recordResp.Hash)
	fmt.Println(recordResp)

}
