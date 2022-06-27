package integration_test

import (
	"strings"

	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irismod-sdk-go/token"
)

func (s IntegrationTestSuite) TestToken() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "TEST",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	issueTokenReq := token.IssueTokenRequest{
		Symbol:        strings.ToLower(s.RandStringOfLength(3)),
		Name:          s.RandStringOfLength(8),
		Scale:         0,
		MinUnit:       strings.ToLower(s.RandStringOfLength(3)),
		InitialSupply: 10000,
		MaxSupply:     100000,
		Mintable:      true,
	}

	//test issue token
	rs, err := s.Token.IssueToken(issueTokenReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	//test mint token
	receipt := s.GetRandAccount().Address.String()
	rs, err = s.Token.MintToken(issueTokenReq.Symbol, 1000, receipt, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	account, err := s.Bank.QueryAccount(receipt)
	require.NoError(s.T(), err)

	amt := sdk.NewIntWithDecimal(1000, int(issueTokenReq.Scale))
	require.Equal(s.T(), amt, account.Coins.AmountOf(issueTokenReq.MinUnit))

	editTokenReq := token.EditTokenRequest{
		Symbol:    issueTokenReq.Symbol,
		Name:      "ethereum network",
		MaxSupply: 20000000,
		Mintable:  false,
	}

	//test edit token
	rs, err = s.Token.EditToken(editTokenReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	//test transfer token
	rs, err = s.Token.TransferToken(receipt, issueTokenReq.Symbol, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	t1, er := s.Token.QueryToken(issueTokenReq.Symbol)
	require.NoError(s.T(), er)
	require.Equal(s.T(), t1.Name, editTokenReq.Name)
	require.Equal(s.T(), t1.MaxSupply, editTokenReq.MaxSupply)
	require.Equal(s.T(), t1.Mintable, editTokenReq.Mintable)
	require.Equal(s.T(), receipt, t1.Owner)

	tokens, er := s.Token.QueryTokens(receipt)
	require.NoError(s.T(), er)
	require.Contains(s.T(), tokens, t1)

	feeToken, er := s.Token.QueryFees(issueTokenReq.Symbol)
	require.NoError(s.T(), er)
	require.Equal(s.T(), true, feeToken.Exist)
	require.Equal(s.T(), "60000000000uirita", feeToken.IssueFee.String())
	require.Equal(s.T(), "6000000000uirita", feeToken.MintFee.String())

	res, er := s.Token.QueryParams()
	require.NoError(s.T(), er)
	require.Equal(s.T(), "0.100000000000000000", res.MintTokenFeeRatio)
	require.Equal(s.T(), "0.400000000000000000", res.TokenTaxRate)
	require.Equal(s.T(), "60000uirita", res.IssueTokenBaseFee)
}
