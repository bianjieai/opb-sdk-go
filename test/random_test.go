package test

import (
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/random"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomRequestRandom(t *testing.T) {
	serviceFeeCap, err := types.ParseCoins("10uirita")
	require.NoError(t, err)

	req := random.RequestRandomRequest{
		BlockInterval: 0,
		Oracle:        false,
		ServiceFeeCap: serviceFeeCap,
	}

	resp, res, err := txClient.Random.RequestRandom(req, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	require.Len(t, resp.ReqID, 64)
	require.Greater(t, resp.Height, int64(0))

	//integration.testRandom.reqId = resp.ReqID
	//integration.testRandom.generateHeight = resp.Height
}
