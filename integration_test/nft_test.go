package integration_test

import (
	"fmt"
	"github.com/irisnet/core-sdk-go/types/query"
	"strings"

	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irismod-sdk-go/nft"
)

func (s IntegrationTestSuite) TestNFT() {

	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	denomID := strings.ToLower(s.RandStringOfLength(4))
	denomName := strings.ToLower(s.RandStringOfLength(4))
	schema := strings.ToLower(s.RandStringOfLength(10))
	issueReq := nft.IssueDenomRequest{
		ID:     denomID,
		Name:   denomName,
		Schema: schema,
	}

	msg := &nft.MsgIssueDenom{
		Id:     denomID,
		Name:   denomName,
		Schema: schema,
		Sender: addr,
	}
	txhash, err := s.BuildTxHash([]sdk.Msg{msg}, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), txhash)
	fmt.Println(txhash)

	res, err := s.NFT.IssueDenom(issueReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
	fmt.Println(res.Hash)

	tokenID := strings.ToLower(s.RandStringOfLength(7))
	tokenName := strings.ToLower(s.RandStringOfLength(7))
	tokenData := strings.ToLower(s.RandStringOfLength(7))
	mintReq := nft.MintNFTRequest{
		Denom: denomID,
		ID:    tokenID,
		Name:  tokenName,
		URI:   fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
		Data:  tokenData,
	}
	res, err = s.NFT.MintNFT(mintReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	editReq := nft.EditNFTRequest{
		Denom: mintReq.Denom,
		ID:    mintReq.ID,
		URI:   fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
	}
	res, err = s.NFT.EditNFT(editReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	nftRes, err := s.NFT.QueryNFT(mintReq.Denom, mintReq.ID)
	require.NoError(s.T(), err)
	fmt.Println(nftRes)
	require.Equal(s.T(), editReq.URI, nftRes.URI)

	supply, err := s.NFT.QuerySupply(mintReq.Denom, nftRes.Creator)
	require.NoError(s.T(), err)
	fmt.Println(supply)
	require.Equal(s.T(), uint64(1), supply)
	pagination := query.PageRequest{
		//Key:        []byte{1},
		Offset:     0,
		Limit:      10,
		CountTotal: true,
	}
	owner, err := s.NFT.QueryOwner(nftRes.Creator, mintReq.Denom, &pagination)
	require.NoError(s.T(), err)
	fmt.Println(owner)
	require.Len(s.T(), owner.IDCs, 1)
	require.Len(s.T(), owner.IDCs[0].TokenIDs, 1)
	require.Equal(s.T(), tokenID, owner.IDCs[0].TokenIDs[0])

	uName := s.RandStringOfLength(10)
	pwd := "11111111"

	recipient, _, err := s.Add(uName, pwd)
	require.NoError(s.T(), err)

	transferReq := nft.TransferNFTRequest{
		Recipient: recipient,
		Denom:     mintReq.Denom,
		ID:        mintReq.ID,
		URI:       fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
	}
	res, err = s.NFT.TransferNFT(transferReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	owner, err = s.NFT.QueryOwner(transferReq.Recipient, mintReq.Denom, &pagination)
	require.NoError(s.T(), err)
	require.Len(s.T(), owner.IDCs, 1)
	require.Len(s.T(), owner.IDCs[0].TokenIDs, 1)
	require.Equal(s.T(), tokenID, owner.IDCs[0].TokenIDs[0])

	supply, err = s.NFT.QuerySupply(mintReq.Denom, transferReq.Recipient)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(1), supply)

	denoms, err := s.NFT.QueryDenoms(&pagination)
	require.NoError(s.T(), err)
	fmt.Println(denoms)
	require.NotEmpty(s.T(), denoms)

	d, err := s.NFT.QueryDenom(denomID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), denomID, d.ID)
	require.Equal(s.T(), denomName, d.Name)
	require.Equal(s.T(), schema, d.Schema)

	col, err := s.NFT.QueryCollection(denomID, &pagination)
	require.NoError(s.T(), err)
	fmt.Println(col)
	require.EqualValues(s.T(), d, col.Denom)
	require.Len(s.T(), col.NFTs, 1)
	require.Equal(s.T(), mintReq.ID, col.NFTs[0].ID)

	burnReq := nft.BurnNFTRequest{
		Denom: mintReq.Denom,
		ID:    mintReq.ID,
	}

	amount, e := sdk.ParseDecCoins("10uirita")
	require.NoError(s.T(), e)
	_, err = s.Bank.Send(recipient, amount, baseTx)
	require.NoError(s.T(), err)

	baseTx.From = uName
	baseTx.Password = pwd
	res, err = s.NFT.BurnNFT(burnReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	supply, err = s.NFT.QuerySupply(mintReq.Denom, transferReq.Recipient)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(0), supply)

	//test TransferDenom
	//transferDenomReq := nft.TransferDenomRequest{
	//	Recipient: recipient,
	//	ID:        mintReq.ID,
	//}
	//res, err = s.NFT.TransferDenom(transferDenomReq, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), res.Hash)
}
