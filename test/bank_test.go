package test

import (
	"fmt"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"
	"testing"
)

// 查询账户
func TestQueryAccount(t *testing.T) {
	account, err := txClient.Bank.QueryAccount(address)
	if err != nil {
		fmt.Println(fmt.Errorf("BANK 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("bank 查询成功:", account)
	}

}

// 发送
func TestSend(t *testing.T) {
	amount, _ := types.ParseDecCoins("10upoint")
	toAddr, _, _ := txClient.Key.Add("testToAddress", "12345678")
	result, err := txClient.Bank.Send(toAddr, amount, baseTx)
	if err != nil {
		fmt.Println(fmt.Errorf("BANK 发送失败: %s", err.Error()))
		return
	} else {
		fmt.Println("BANK 发送成功：", result.Hash)
		e := syncTx(result.Hash)
		require.NoError(t, e)
	}
}
