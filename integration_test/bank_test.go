package integration_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/core-sdk-go/bank"
	"github.com/irisnet/core-sdk-go/types"
)

func (s IntegrationTestSuite) TestBank() {
	cases := []SubTest{
		{
			"TestQueryAccount",
			queryAccount,
		},
		{
			"TestSend",
			send,
		},
		{
			"TestMultiSend",
			multiSend,
		},
		{
			"TestSimulate",
			simulate,
		},
		{
			"TestSendWitchSpecAccountInfo",
			sendWitchSpecAccountInfo,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() { t.testCase(s) })
	}
}

func queryAccount(s IntegrationTestSuite) {
	account, err := s.Bank.QueryAccount(s.Account().Address.String())
	s.NoError(err)
	s.NotEmpty(account)
	bz, _ := json.Marshal(account)
	fmt.Println(string(bz))
}

func send(s IntegrationTestSuite) {
	coins, err := types.ParseDecCoins("10uirita")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()

	ch := make(chan int)
	s.Bank.SubscribeSendTx(s.Account().Address.String(), to, func(send bank.EventDataMsgSend) {
		ch <- 1
	})
	s.NoError(err)
	baseTx := types.BaseTx{
		From:               s.Account().Name,
		Gas:                200000,
		Memo:               "TEST",
		Mode:               types.Commit,
		Password:           s.Account().Password,
		SimulateAndExecute: false,
		GasAdjustment:      1.5,
	}

	res, err := s.Bank.Send(to, coins, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
	time.Sleep(1 * time.Second)

	resp, err := s.Manager().QueryTx(res.Hash)
	s.NoError(err)
	s.Equal(resp.Result.Code, uint32(0))
	s.Equal(resp.Height, res.Height)

	<-ch
}

func multiSend(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      2000000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	coins, err := types.ParseDecCoins("10uirita")
	s.NoError(err)

	accNum := 11
	acc := make([]string, accNum)
	receipts := make([]bank.Receipt, accNum)
	for i := 0; i < accNum; i++ {
		acc[i] = s.RandStringOfLength(10)
		addr, _, err := s.Add(acc[i], "12345678")

		s.NoError(err)
		s.NotEmpty(addr)

		receipts[i] = bank.Receipt{
			Address: addr,
			Amount:  coins,
		}
	}
	_, err = s.Bank.MultiSend(bank.MultiSendRequest{Receipts: receipts}, baseTx)
	s.NoError(err)

	coins, err = types.ParseDecCoins("2uirita")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()
	begin := time.Now()
	var wait sync.WaitGroup
	for i := 1; i < 5; i++ {
		wait.Add(1)
		index := rand.Intn(accNum)
		go func() {
			defer wait.Done()
			_, err := s.Bank.Send(to, coins, types.BaseTx{
				From:     acc[index],
				Gas:      200000,
				Memo:     "test",
				Mode:     types.Commit,
				Password: "12345678",
			})
			s.NoError(err)
		}()
	}
	wait.Wait()
	end := time.Now()
	fmt.Printf("total senconds:%s\n", end.Sub(begin).String())
}

func simulate(s IntegrationTestSuite) {
	coins, err := types.ParseDecCoins("10uirita")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()
	baseTx := types.BaseTx{
		From:               s.Account().Name,
		Password:           s.Account().Password,
		Gas:                200000,
		Memo:               "test",
		Mode:               types.Commit,
		SimulateAndExecute: true,
	}

	result, err := s.Bank.Send(to, coins, baseTx)
	s.NoError(err)
	s.Greater(result.GasWanted, int64(0))
	fmt.Println(result)
}

func sendWitchSpecAccountInfo(s IntegrationTestSuite) {
	for i := 0; i < 10; i++ {
		coins, err := types.ParseDecCoins("10uirita")
		baseTx := types.BaseTx{
			From:     s.Account().Name,
			Gas:      200000,
			Fee:      coins,
			Memo:     "TEST",
			Mode:     types.Commit,
			Password: s.Account().Password,
		}

		curAccount, err := s.Bank.QueryAccount(s.Account().Address.String())
		require.NoError(s.T(), err)

		accountNumber := curAccount.AccountNumber
		sequence := curAccount.Sequence
		randomAddr := s.GetRandAccount().Address.String()
		amount, _ := types.ParseDecCoins("10uirita")

		res, err := s.Bank.SendWitchSpecAccountInfo(randomAddr, sequence, accountNumber, amount, baseTx)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), res.Hash)
	}
}
