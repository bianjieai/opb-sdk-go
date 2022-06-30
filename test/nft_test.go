package test

import (
	"fmt"
	"github.com/irisnet/core-sdk-go/types/query"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/stretchr/testify/require"
	"testing"
)

var pagination = query.PageRequest{
	//Key:        []byte{1},
	Offset:     0,
	Limit:      10,
	CountTotal: true,
}

// 发行资产。
func TestNftIssue(t *testing.T) {
	denomID := "testnftdenomid"
	denomName := "test_nft_denomName"
	schema := ""
	issueReq := nft.IssueDenomRequest{
		ID:     denomID,
		Name:   denomName,
		Schema: schema,
	}
	res, err := txClient.NFT.IssueDenom(issueReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	fmt.Println(res.Hash) //B43434AFD14780D0024AC13BDA08AE41CF6906BD16F61192EE360E6A262CF54B
}

// 创建指定类别的具体资产。
func TestNftMint(t *testing.T) {
	tokenID := "testnfttokenid"
	tokenName := "test_nft_tokenName"
	tokenData := ""
	denomID := "testnftdenomid"
	mintReq := nft.MintNFTRequest{
		Denom: denomID,
		ID:    tokenID,
		Name:  tokenName,
		URI:   "https://www.baidu.com",
		Data:  tokenData,
	}
	res, err := txClient.NFT.MintNFT(mintReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}

// 编辑指定的资产。可更新的属性包括：资产元数据、元数据 URI、URI 的哈希
func TestNftEdit(t *testing.T) {
	tokenID := "testnfttokenid"
	denomID := "testnftdenomid"
	editReq := nft.EditNFTRequest{
		Denom: denomID,
		ID:    tokenID,
		URI:   "https://www.baidu.com",
	}
	res, err := txClient.NFT.EditNFT(editReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}

// 转移指定资产。
func TestNftTransfer(t *testing.T) {
	recipient := address
	tokenID := "testnfttokenid"
	denomID := "testnftdenomid"
	transferReq := nft.TransferNFTRequest{
		Recipient: recipient,
		Denom:     denomID,
		ID:        tokenID,
		URI:       "https://www.baidu.com",
	}
	res, err := txClient.NFT.TransferNFT(transferReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}

// 销毁指定资产。
func TestNftBurn(t *testing.T) {
	tokenID := "testnfttokenid"
	denomID := "testnftdenomid"
	burnReq := nft.BurnNFTRequest{
		Denom: denomID,
		ID:    tokenID,
	}
	res, err := txClient.NFT.BurnNFT(burnReq, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}

// 查询指定类别和 ID 的资产。
func TestNftToken(t *testing.T) {
	tokenID := "testnfttokenid"
	denomID := "testnftdenomid"
	nftRes, err := txClient.NFT.QueryNFT(denomID, tokenID)
	require.NoError(t, err)
	fmt.Println(nftRes)
}

// 查询指定类别的资产信息。
func TestNftDenom(t *testing.T) {
	denomID := "testnftdenomid"
	denomRes, err := txClient.NFT.QueryDenom(denomID)
	require.NoError(t, err)
	fmt.Println(denomRes)
}

// 查询所有类别的资产信息。
func TestNftDenoms(t *testing.T) {
	denoms, err := txClient.NFT.QueryDenoms(&pagination)
	require.NoError(t, err)
	require.NotEmpty(t, denoms)
	fmt.Println(denoms)
}

// 查询指定类别资产的总量。如 owner 被指定，则查询此 owner 所拥有的该类别资产的总量。
func TestNftSupply(t *testing.T) {
	recipient := address
	denomID := "testnftdenomid"
	supply, err := txClient.NFT.QuerySupply(denomID, recipient)
	require.NoError(t, err)
	require.Equal(t, uint64(1), supply)
	fmt.Println(supply)
}

// 查询指定账户的资产列表。如提供 denom，则查询该账户指定 denom 的资产列表。
func TestNftOwner(t *testing.T) {
	recipient := address
	denomID := "testnftdenomid"
	owner, err := txClient.NFT.QueryOwner(recipient, denomID, &pagination)
	require.NoError(t, err)
	fmt.Println(owner)
}

// 查询指定类别的所有资产。
func TestNftCollection(t *testing.T) {
	denomID := "testnftdenomid"
	col, err := txClient.NFT.QueryCollection(denomID, &pagination)
	require.NoError(t, err)
	fmt.Println(col)
}
