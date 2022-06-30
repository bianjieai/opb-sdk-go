package integration_test

import (
	"fmt"
	"github.com/irisnet/irismod-sdk-go/mt"
	"github.com/stretchr/testify/require"
	"testing"
)

// 发行资产类别
func TestMtIssue(t *testing.T) {
	denomName := "test_denomname" // 资产类别的名称
	dataStr := "test_data"        // 资产元数据
	denomData := []byte(dataStr)
	issueReq := mt.IssueDenomRequest{
		Name: denomName,
		Data: denomData,
	}

	res, err := txClient.MT.IssueDenom(issueReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	fmt.Println(res.Hash)
	denomID, err2 := res.Events.GetValue("issue_denom", "denom_id")
	require.NoError(t, err2)
	require.NotEmpty(t, denomID)
	fmt.Println(denomID) //259edd57e552854d42bc4a0d98dc7a48fddeae343ad428c9df1d4b09e0ab525a
}

// 转移指定资产类别
func TestMtTransferDenom(t *testing.T) {

}

// 创建指定类别的具体资产
func TestMtMint(t *testing.T) {
	mintMTData := []byte("test_mintMTData")
	mtRecipient := address
	amount := uint64(10)
	mintReq := mt.MintMTRequest{
		DenomID:   "259edd57e552854d42bc4a0d98dc7a48fddeae343ad428c9df1d4b09e0ab525a",
		Amount:    amount,
		Data:      mintMTData,
		Recipient: mtRecipient,
	}
	res, err := txClient.MT.MintMT(mintReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	mtID, err2 := res.Events.GetValue("mint_mt", "mt_id")
	require.NoError(t, err2)
	fmt.Println(mtID) //f28fea81a5cdbb341979a9bb9b0c620226b8ba5077c49fb4a60d630b3a53a161
}

// 编辑指定的资产。
// 可更新的属性包括：资产元数据。
func TestMtEdit(t *testing.T) {
	editMTData := []byte("test_editMTData")
	editReq := mt.EditMTRequest{
		DenomID: "259edd57e552854d42bc4a0d98dc7a48fddeae343ad428c9df1d4b09e0ab525a",
		ID:      "f28fea81a5cdbb341979a9bb9b0c620226b8ba5077c49fb4a60d630b3a53a161",
		Data:    editMTData,
	}
	res, err := txClient.MT.EditMT(editReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}

// 转移指定资产；可指定转移数量。
func TestTransfer(t *testing.T) {
	name, password := "test_name", "12345678"

	transferMTRecipient, _, err := txClient.Key.Add(name, password)
	require.NoError(t, err)
	require.NotEmpty(t, address)
	//transferMTRecipient := ""
	transferAmount := uint64(5)
	transferMTReq := mt.TransferMTRequest{
		ID:        "f28fea81a5cdbb341979a9bb9b0c620226b8ba5077c49fb4a60d630b3a53a161",
		DenomID:   "259edd57e552854d42bc4a0d98dc7a48fddeae343ad428c9df1d4b09e0ab525a",
		Amount:    transferAmount,
		Recipient: transferMTRecipient,
	}
	res, err := txClient.MT.TransferMT(transferMTReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}
