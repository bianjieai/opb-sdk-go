package integration_test

import (
	"fmt"
	opb "github.com/bianjieai/opb-sdk-go/pkg/app/sdk"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/query"
	"github.com/irisnet/core-sdk-go/types/store"
	"github.com/irisnet/irismod-sdk-go/nft"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNftDemo(t *testing.T) {
	fee, _ := types.ParseDecCoins("1upoint") // 设置文昌链主网的默认费用，10W不够就填20W，30W....
	// 初始化 SDK 配置
	options := []types.Option{
		types.AlgoOption("sm2"),
		types.KeyDAOOption(store.NewMemory(nil)),
		types.FeeOption(fee),
		types.TimeoutOption(10),
		types.CachedOption(true),
	}
	//cfg, err := types.NewClientConfig("http://47.100.192.234:26657", "47.100.192.234:9090", "testing", options...)
	cfg, err := types.NewClientConfig("http://localhost:26657", "localhost:9090", "irita-test", options...)
	if err != nil {
		panic(err)
	}

	//// 初始化 OPB 网关账号（测试网环境设置为 nil 即可）
	//authToken := model.NewAuthToken("TestProjectID", "TestProjectKey", "TestChainAccountAddress")
	//
	//// 开启 TLS 连接
	//// 若服务器要求使用安全链接，此处应设为true；若此处设为false可能导致请求出现长时间不响应的情况
	//authToken.SetRequireTransportSecurity(false)
	// 创建 OPB 客户端
	//client := opb.NewClient(cfg, &authToken)
	client := opb.NewClient(cfg, nil)

	// 导入私钥
	address, _ := client.Key.Recover("validator", "12345678", "tool clinic access federal flight egg sand lonely book aunt plastic lemon either diagram water betray life gaze spice upon kingdom much coyote flat")
	fmt.Println("address:", address)
	//address := "iaa14hdjjrpljywqrpmjefmn4e9umcu93s3eku7es7"
	// 初始化 Tx 基础参数
	baseTx := types.BaseTx{
		From:     "validator", // 对应上面导入的私钥名称
		Password: "12345678",  // 对应上面导入的私钥密码
		Gas:      200000,      // 单 Tx 消耗的 Gas 上限
		Memo:     "",          // Tx 备注
		Mode:     types.Sync,  // Tx 广播模式
	}
	rand.Seed(time.Now().UnixNano())
	denomID := "testdenom" + strconv.Itoa(rand.Intn(10000)) //资产的类别，全局唯一；长度为3到64，字母数字字符，以字母开始
	denomName := "test_name"
	schema := "no schema"
	fmt.Println(denomID)
	// 发行资产类别
	issueReq := nft.IssueDenomRequest{
		ID:          denomID,
		Name:        denomName,
		Schema:      schema,
		Symbol:      "symbol",
		Description: "test_description",
		Uri:         "https://www.baidu.com",
		Data:        "any data",
	}
	nftDenomResult, err := client.NFT.IssueDenom(issueReq, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("NFT 类别创建失败: %s", err.Error()))
		return
	} else {
		fmt.Println("NFT 类别创建成功 TxHash：", nftDenomResult.Hash)
	}

	nftID := "test" + strconv.Itoa(rand.Intn(10000)) //资产的唯一 ID，如 UUID
	fmt.Println(nftID)
	mintReq := nft.MintNFTRequest{
		Denom: denomID,
		ID:    nftID,
		Name:  "test_nftName",
		URI:   "https://www.baidu.com",
		Data:  "any data",
	}
	nftResult, err := client.NFT.MintNFT(mintReq, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("NFT 创建失败: %s", err.Error()))
		return
	} else {
		fmt.Println("NFT 创建成功 TxHash：", nftResult.Hash)
	}

	//editReq := nft.EditNFTRequest{
	//	Denom: denomID,
	//	ID:    nftID,
	//	URI:   "https://www.baidu.com",
	//}
	//editNftResult, err := client.NFT.EditNFT(editReq, baseTx)
	//if err != nil {
	//	fmt.Println(fmt.Errorf("NFT 更新失败: %s", err.Error()))
	//	return
	//} else {
	//	fmt.Println("NFT 更新成功 TxHash：", editNftResult.Hash)
	//}

	QueryNFTResult, err := client.NFT.QueryNFT(denomID, nftID)
	if err != nil {
		fmt.Println(fmt.Errorf("NFT 查询失败: %s", err.Error()))
		//return
	} else {
		fmt.Println("NFT 查询成功：", QueryNFTResult)
	}

	supply, err := client.NFT.QuerySupply(denomID, "validator")
	if err != nil {
		fmt.Println(fmt.Errorf("supply 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("supply 查询成功：", supply)
	}
	// 分页查询
	pagination := query.PageRequest{
		//Key: []byte{1},  //键值：与offset二选一
		Offset:     0,     //偏移量：与key二选一
		Limit:      5,     //默认100，最大值100
		CountTotal: false, //是否查询总数量，目前仅支持false
	}
	owner, err := client.NFT.QueryOwner("validator", denomID, &pagination)
	if err != nil {
		fmt.Println(fmt.Errorf("owner 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("owner 查询成功：", owner)
	}

	// 新建账户：接收者
	uName := "test"
	pwd := "12345678"

	recipient, _, err := client.Add(uName, pwd)
	if err != nil {
		fmt.Println(fmt.Errorf("key 添加失败: %s", err.Error()))
		return
	} else {
		fmt.Println("key 添加成功：", recipient)
	}

	transferReq := nft.TransferNFTRequest{
		Recipient: recipient,
		Denom:     denomID,
		ID:        nftID,
		URI:       "",
	}
	res, err := client.NFT.TransferNFT(transferReq, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("转移失败: %s", err.Error()))
		return
	} else {
		fmt.Println("转移成功：", res.Hash)
	}

	owner, err = client.NFT.QueryOwner(recipient, denomID, &pagination)
	if err != nil {
		fmt.Println(fmt.Errorf("owner 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("owner 查询成功：", owner)
	}

	supply, err = client.NFT.QuerySupply(denomID, recipient)
	if err != nil {
		fmt.Println(fmt.Errorf("supply 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("supply 查询成功：", supply)
	}

	denoms, err := client.NFT.QueryDenoms(&pagination)
	if err != nil {
		fmt.Println(fmt.Errorf("denoms 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("denoms 查询成功：", denoms)
	}

	d, err := client.NFT.QueryDenom(denomID)
	if err != nil {
		fmt.Println(fmt.Errorf("denom 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("denom 查询成功：", d)
	}

	col, err := client.NFT.QueryCollection(denomID, &pagination)
	if err != nil {
		fmt.Println(fmt.Errorf("Collection 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("Collection 查询成功：", col)
	}

	//burnReq := nft.BurnNFTRequest{
	//	Denom: mintReq.Denom,
	//	ID:    mintReq.ID,
	//}
	//
	//amount, e := types.ParseDecCoins("10uirita")
	//if e!=nil{
	//	fmt.Println(e)
	//}
	//_, err = client.Bank.Send(recipient, amount, baseTx)
	//if err!=nil{
	//	fmt.Println(err)
	//}

	//baseTx.From = uName
	//baseTx.Password = pwd
	//res, err = client.NFT.BurnNFT(burnReq, baseTx)
	//if err != nil {
	//	fmt.Println(fmt.Errorf("nft 销毁失败: %s", err.Error()))
	//} else {
	//	fmt.Println("nft 销毁成功：", res)
	//}
	//
	//supply, err = client.NFT.QuerySupply(mintReq.Denom, transferReq.Recipient)
	//if err != nil {
	//	fmt.Println(fmt.Errorf("supply 查询失败: %s", err.Error()))
	//} else {
	//	fmt.Println("supply 查询成功：", supply)
	//}

	//test TransferDenom
	//transferDenomReq := nft.TransferDenomRequest{
	//	Recipient: recipient,
	//	ID:        mintReq.ID,
	//}
	//res, err = s.NFT.TransferDenom(transferDenomReq, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), res.Hash)
}
