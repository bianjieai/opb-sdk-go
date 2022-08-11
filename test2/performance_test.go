package test2

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/client"

	"github.com/avast/retry-go"
	coreType "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/query"
	"github.com/irisnet/core-sdk-go/types/store"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/irisnet/irismod-sdk-go/service"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/abci/types"
)

var iritaClient *client.Client

var baseTx = coreType.BaseTx{
	From:     "node0",
	Password: "12345678",
	//Fee:                coreType.ParseDecCoins("3stake"),
	Mode:               coreType.Commit,
	Gas:                5000000,
	SimulateAndExecute: false,
	GasAdjustment:      0,
}

var svcSchema = `{
	"input": {
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "test_svc_req",
		"description": "aaa",
		"type": "object",
		"properties": {
			"goid": {
				"description": "global oid",
				"type": "string"
			}
		}
	},
	"output": {
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "test_svc_resp",
		"description": "aaa",
		"type": "object",
		"properties": {
			"code": {
				"description": "state code",
				"type": "string"
			},
			"message": {
				"description": "respond message",
				"type": "string"
			},
			"svc_name": {
				"description": "服务名称",
				"type": "string"
			},
			"data": {
				"description": "return data",
				"loid": {
					"description": "local oid",
					"type": "string"
				}
			}
		}
	}
}`

func init() {
	//fee, _ := coreType.ParseDecCoins("30stake")
	fee, _ := coreType.ParseDecCoins("4point")
	options := []coreType.Option{
		coreType.AlgoOption("sm2"),
		coreType.KeyDAOOption(store.NewMemory(nil)),
		coreType.FeeOption(fee),
		coreType.TimeoutOption(20),
		coreType.CachedOption(true),
	}
	config, err := coreType.NewClientConfig("http://10.0.0.36:26657", "10.0.0.36:9090", "cschain-otc", options...)
	//config, err := coreType.NewClientConfig("http://localhost:26660", "localhost:9091", "test", options...)
	if err != nil {
		panic(err)
	}
	cli := client.NewClient(config)
	iritaClient = &cli
	//addr, err := iritaClient.Key.Recover("node0", "12345678", "best canal bronze cloth clean winter wife danger exercise chief matter trust junk inch crowd enjoy leg van aisle arrive increase rather laptop scheme")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("account loaded: ", addr)
	addr, err := iritaClient.Key.Recover("node0", "12345678", "mistake hunt fade appear donkey cause toward quick maid tattoo diary turtle defense banner rose beyond develop repair raise pony guide soldier want start")
	if err != nil {
		panic(err)
	}
	fmt.Println("account loaded: ", addr)
	providerAddr, err := iritaClient.Key.Recover("provider", "12345678", "broccoli bone broom cliff pluck normal jelly guess brush urban depth under become sustain chronic away milk repair prevent latin good rug refuse ill")
	if err != nil {
		panic(err)
	}
	fmt.Println("provider:", providerAddr)
}

func TestTx(t *testing.T) {
	result, err := iritaClient.QueryTx("89AA2DF32B254D579C197F313852EE1A552758ACE6DFF5769CEC0D548AE6F582")
	require.NoError(t, err)
	require.Equal(t, result.Result.Code, types.CodeTypeOK, result.Result.Log)
}

func TestService(t *testing.T) {
	resultTx, err := iritaClient.Service.DefineService(service.DefineServiceRequest{
		ServiceName:       "test_svc",
		Description:       "aaa",
		Tags:              []string{"test_svc"},
		AuthorDescription: "nothing",
		Schemas:           svcSchema,
	}, baseTx)
	require.NoError(t, err)
	t.Log(resultTx)
}

func TestBind(t *testing.T) {
	coins, err2 := coreType.ParseDecCoins("30000000upoint")
	require.NoError(t, err2)
	resultTx, err := iritaClient.Service.BindService(service.BindServiceRequest{
		ServiceName: "test_svc",
		Deposit:     coins,
		Pricing:     "{\"price\":\"1upoint\"}",
		QoS:         100,
		Options:     `{"service_name":"test_svc"}`,
		Provider:    "iaa1pcsr5whtrmdjnygz3uwjd0hmrzk9stppx0a4u6",
	}, baseTx)
	require.NoError(t, err)
	t.Log(resultTx)
}

func TestCall(t *testing.T) {
	coins, err2 := coreType.ParseDecCoins("2point")
	baseTx.Mode = coreType.Commit
	baseTx.Fee = coins
	require.NoError(t, err2)
	reqCtxId, resultTx, err := iritaClient.Service.InvokeService(service.InvokeServiceRequest{
		ServiceName:   "crx_chain_svc",
		Providers:     []string{"iaa18h9lz30t5su37fmcernu3r9t5p452yrgh7sxvv"},
		Input:         `{"header":{"req_sequence":"requestid411628592022-07-0fasdfd1","id":"svcd4bb0a8348872bfa13c5e57bc35facbd854cb1236cd8fb82348f2ad2"},"body":{}}`,
		ServiceFeeCap: coins,
		Timeout:       100,
	}, baseTx)
	require.NoError(t, err)
	t.Log(reqCtxId, resultTx)

	ctx, err := iritaClient.Service.QueryRequestContext(reqCtxId)
	require.NoError(t, err)

	result, err := iritaClient.Service.QueryRequestsByReqCtx(reqCtxId, ctx.BatchCounter, &query.PageRequest{
		Limit: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result, "no request")
	t.Log(result[0].ID)
}

func TestResp(t *testing.T) {
	baseTx.From = "provider"
	resultTx, err := iritaClient.Service.InvokeServiceResponse(service.InvokeServiceResponseRequest{
		RequestId: "883D6A22B4E9F6C079038C8EAD9D651CE6B39DA4FF2810E6D427FD2C3BEDA6460000000000000000000000000000000100000000026AF28A0000",
		Output:    `{"header":{"req_sequence":"requestid411628592022-07-0fasdfd1","id":"svcd4bb0a8348872bfa13c5e57bc35facbd854cb1236cd8fb82348f2ad2"},"body":{"tx_hash":"","result":"{\"code\":\"100001\",\"msg\":\"unknown this grpc request's method:\",\"result\":null}"}}`,
		Result:    `{"code": 200, "message": "success"}`,
	}, baseTx)
	require.NoError(t, err)
	t.Log(resultTx)
}

func TestReqCtxId(t *testing.T) {
	reqCtxId := "12B076F7E2E968033A23F1AE85C78762C6D9CF3549EB8B380EE27F9418D717550000000000000000"
	reqCtx, err := iritaClient.Service.QueryRequestsByReqCtx(reqCtxId, 1, &query.PageRequest{
		Limit: 100,
	})
	require.NoError(t, err)
	t.Log(reqCtx)

	ctx, s := iritaClient.Service.QueryRequestContext(reqCtxId)
	require.NoError(t, s)
	t.Log(ctx)
	req, err := iritaClient.Service.QueryServiceRequest("AEFE0BDA070228CA4C33170E75E141399DBFF8B377D85104841EAE199B0C6E770000000000000000")
	require.NoError(t, err)
	t.Log(req)
}

func TestIssueDenom(t *testing.T) {
	addr, err := iritaClient.Key.Recover("node1", "12345678", "crime disease champion pull receive drama siren rally else proud arrow list expand want roof feel lucky help ball job clean harsh immense motor")
	require.NoError(t, err)
	t.Log(addr, "log in")
	fee, _ := coreType.ParseDecCoins("30stake")
	baseTx.Fee = fee
	baseTx.From = "node1"
	denomId := fmt.Sprintf("testdenom%d", rand.Int63n(10000000000))
	resultTx, err := iritaClient.NFT.IssueDenom(nft.IssueDenomRequest{
		ID:     denomId,
		Name:   denomId,
		Schema: "",
	}, baseTx)
	if err != nil {
		t.Error(err, resultTx.Hash)
		return
	}
	t.Log(resultTx)
}

func Test_BatchNft(t *testing.T) {
	denomId := "testdenom" + time.Now().Format("060102150405")
	resultTx, err := iritaClient.NFT.IssueDenom(nft.IssueDenomRequest{
		ID:     denomId,
		Name:   denomId,
		Schema: "",
	}, baseTx)
	if err != nil {
		t.Error(err, resultTx.Hash)
		return
	}
	for i := 0; i < 10; i++ {
		accName := fmt.Sprintf("testac%d", i)
		addr, _, err := iritaClient.Key.Add(accName, "12345678")
		if err != nil {
			t.Errorf("create account %s err: %s", accName, err)
			continue
		}
		accBaseTx := baseTx
		accBaseTx.From = accName
		accBaseTx.Mode = coreType.Sync
		go func(accName string, addr string, accBaseTx coreType.BaseTx, baseTx coreType.BaseTx) {
			baseTx.Mode = coreType.Sync
			resultTx, err1 := iritaClient.Bank.Send(addr, coreType.DecCoins{coreType.NewDecCoin("stake", coreType.NewInt(10000))}, baseTx)
			if err1 != nil {
				t.Errorf("send token failed: %s", err1)
				return
			}
			retry.Do(func() error {
				_, err2 := iritaClient.QueryTx(resultTx.Hash)
				return err2
			},
				retry.RetryIf(func(err error) bool {
					return true
				}))

			for j := 0; true; j++ {
				msgs := coreType.Msgs{}
				for k := 0; k < 1000; k++ {
					nftID := fmt.Sprintf("test%d", rand.Int63n(9999999999999999))
					msgs = append(msgs, &nft.MsgMintNFT{
						DenomId:   denomId,
						Id:        nftID,
						Name:      nftID,
						Sender:    addr,
						Recipient: addr,
					})

				}
				_, nftErr := iritaClient.BuildAndSend(msgs, accBaseTx)
				if nftErr != nil {
					fmt.Printf("acc:%s, mint err.", accName)
					return
				}
				time.Sleep(200 * time.Millisecond)
			}

		}(accName, addr, accBaseTx, baseTx)
	}

	select {}
}
