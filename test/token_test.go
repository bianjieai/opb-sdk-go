package test

import (
	"fmt"
	"github.com/irisnet/irismod-sdk-go/token"
	"github.com/stretchr/testify/require"
	"testing"
)

// 发行工分。
func TestTokenIssue(t *testing.T) {
	issueTokenReq := token.IssueTokenRequest{
		Symbol:        "testtokensymbol",
		Name:          "testtokenname",
		Scale:         0,
		MinUnit:       "testtokenunit",
		InitialSupply: 10000,
		MaxSupply:     100000,
		Mintable:      true,
	}

	res, err := txClient.Token.IssueToken(issueTokenReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}

// 编辑存在的工分。可编辑的属性包括：名称、最大供应以及可增发性
func TestTokenEdit(t *testing.T) {
	editTokenReq := token.EditTokenRequest{
		Symbol:    "testtokensymbol",
		Name:      "testtokenname",
		MaxSupply: 20000000,
		Mintable:  true,
	}

	res, err := txClient.Token.EditToken(editTokenReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}

// 增发工分到指定地址。
func TestTokenMint(t *testing.T) {
	to := address
	res, err := txClient.Token.MintToken("testtokensymbol", 100, to, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}

// 转让工分所有权。
func TestTokenTransfer(t *testing.T) {
	res, err := txClient.Token.TransferToken("iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj", "testtokensymbol", baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}

// 查询指定的工分。
func TestTokenToken(t *testing.T) {
	queryToken, er := txClient.Token.QueryToken("testtokensymbol")
	require.NoError(t, er)
	fmt.Println(queryToken)
}

// 查询已发行的所有工分。
func TestTokenTokens(t *testing.T) {
	owner := address
	tokens, er := txClient.Token.QueryTokens(owner, pagination)
	require.NoError(t, er)
	fmt.Println(tokens)
}

// 查询工分相关的费用，包括发行和增发。
func TestTokenFee(t *testing.T) {
	feeToken, er := txClient.Token.QueryFees("testtokensymbol")
	require.NoError(t, er)
	fmt.Println(feeToken)
}

func TestParams(t *testing.T) {
	res, er := txClient.Token.QueryParams()
	require.NoError(t, er)
	fmt.Println(res)
}
