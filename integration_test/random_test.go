package integration_test

import (
	"strconv"
	"time"

	"github.com/irisnet/core-sdk-go/types"

	"github.com/irisnet/irismod-sdk-go/random"
)

type TestRandom struct {
	reqId          string
	generateHeight int64
}

var testRandom TestRandom

func (s IntegrationTestSuite) TestRandom() {

	cases := []SubTest{
		{
			"TestRequestRandom",
			requestRandom,
		},
		{
			"TestQueryRandom",
			queryRandom,
		},
		{
			"TestQueryRandomRequestQueue",
			queryRandomRequestQueue,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

func requestRandom(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
	}
	serviceFeeCap, err := types.ParseCoins("10uirita")
	s.NoError(err)

	req := random.RequestRandomRequest{
		BlockInterval: 0,
		Oracle:        false,
		ServiceFeeCap: serviceFeeCap,
	}

	resp, res, err := s.Random.RequestRandom(req, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
	s.Len(resp.ReqID, 64)
	s.Greater(resp.Height, int64(0))

	testRandom.reqId = resp.ReqID
	testRandom.generateHeight = resp.Height
}

func queryRandom(s IntegrationTestSuite) {
	// Wait for the transaction to be packaged into the block
	time.Sleep(10 * time.Second)
	res, err := s.Random.QueryRandom(testRandom.reqId)
	s.NoError(err)
	s.NotEmpty(res.RequestTxHash)
	value, _ := strconv.ParseFloat(res.Value, 10)
	s.Greater(value, float64(0))
}

func queryRandomRequestQueue(s IntegrationTestSuite) {
	_, err := s.Random.QueryRandomRequestQueue(testRandom.generateHeight)
	s.NoError(err)
	//s.NotEmpty(queue)
}
