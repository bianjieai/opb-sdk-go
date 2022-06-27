package integration_test

import (
	"fmt"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/mt"
	"github.com/stretchr/testify/require"
	"strings"
)

func (s IntegrationTestSuite) TestMT() {
	dec := sdk.NewDecCoin("upoint", sdk.NewInt(100000))
	fee := sdk.NewDecCoins(dec)
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      1000000000,
		Memo:     "test",
		Fee:      fee,
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	denomName := strings.ToLower(s.RandStringOfLength(4))
	dataStr := strings.ToLower(s.RandStringOfLength(10))
	denomData := []byte(dataStr)
	issueReq := mt.IssueDenomRequest{
		Name: denomName,
		Data: denomData,
	}

	res, err := s.MT.IssueDenom(issueReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
	fmt.Println(res.Hash)
	denomID, err2 := res.Events.GetValue("issue_denom", "denom_id")
	require.NoError(s.T(), err2)
	require.NotEmpty(s.T(), denomID)

	mintMTData := []byte(strings.ToLower(s.RandStringOfLength(7)))
	mtRecipient := s.Account().Address.String()
	amount := uint64(10)
	mintReq := mt.MintMTRequest{
		DenomID:   denomID,
		Amount:    amount,
		Data:      mintMTData,
		Recipient: mtRecipient,
	}
	res, err = s.MT.MintMT(mintReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
	mtID, err2 := res.Events.GetValue("mint_mt", "mt_id")
	require.NoError(s.T(), err2)

	editMTData := []byte(strings.ToLower(s.RandStringOfLength(8)))
	editReq := mt.EditMTRequest{
		DenomID: denomID,
		ID:      mtID,
		Data:    editMTData,
	}
	res, err = s.MT.EditMT(editReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	transferMTRecipient := s.randAccounts[3].Address.String()
	transferAmount := uint64(5)
	transferMTReq := mt.TransferMTRequest{
		ID:        mtID,
		DenomID:   mintReq.DenomID,
		Amount:    transferAmount,
		Recipient: transferMTRecipient,
	}
	res, err = s.MT.TransferMT(transferMTReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	//supply, err := s.MT.QuerySupply(mintReq.DenomID, transferMTRecipient)
	//require.NoError(s.T(), err)
	//require.Equal(s.T(), uint64(1), supply)

	denom, err := s.MT.QueryDenom(denomID)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), denom)

	denoms, err := s.MT.QueryDenoms()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), denoms)

	mtSupply, err := s.MT.QueryMTSupply(denomID, mtID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(10), mtSupply)

	mtSingle, err := s.MT.QueryMT(denomID, mtID)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), mtSingle)
	require.Equal(s.T(), mtSingle.Data, editMTData)

	mts, err := s.MT.QueryMTs(denomID)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), mts)

	balances, err := s.MT.QueryBalances(denomID, transferMTRecipient)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), balances)
	require.Equal(s.T(), balances[0].Amount, amount-transferAmount)

	//burnMTReq := s.Account().Address.String()
	burnMTReq := mt.BurnMTRequest{
		ID:      mtID,
		DenomID: denomID,
		Amount:  amount - transferAmount,
	}
	res, err = s.MT.BurnMT(burnMTReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	mtSupplyCheck, err := s.MT.QueryMTSupply(denomID, mtID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), transferAmount, mtSupplyCheck)

	recipient := s.randAccounts[2].Address.String()
	transferReq := mt.TransferDenomRequest{
		ID:        denomID,
		Recipient: recipient,
	}
	res, err = s.MT.TransferDenom(transferReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
	fmt.Println(res.Hash)
}
