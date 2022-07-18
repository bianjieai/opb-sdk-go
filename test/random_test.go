package test

import (
	"fmt"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/random"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

// request	请求具有可选块间隔的随机数
func TestRandomRequest(t *testing.T) {
	serviceFeeCap, err := types.ParseCoins("10uirita")
	require.NoError(t, err)

	req := random.RequestRandomRequest{
		BlockInterval: 100,
		Oracle:        false,
		ServiceFeeCap: serviceFeeCap,
	}

	resp, res, err := txClient.Random.RequestRandom(req, baseTx)
	require.NoError(t, err)
	fmt.Println(res)
	fmt.Println(resp.ReqID)
	fmt.Println(resp.Height)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}

// 使用ID查询链上生成的随机数。
func TestRandomQuery(t *testing.T) {
	reqId := "99fc7a859a028082ff818b98bec2376230581e547b52677e9da661ed49a5ded8"
	res, err := txClient.Random.QueryRandom(reqId)
	require.NoError(t, err)
	fmt.Println(res.RequestTxHash)
	value, _ := strconv.ParseFloat(res.Value, 10)
	fmt.Println(value)
}

// 查询随机数请求队列，支持可选的高度。
func TestRandomQueryQueue(t *testing.T) {
	res, err := txClient.Random.QueryRandomRequestQueue(21801)
	require.NoError(t, err)
	fmt.Println(res)
}
