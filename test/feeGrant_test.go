package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/core-sdk-go/feegrant"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
)

var granter types.AccAddress
var grantee types.AccAddress

func init() {
	// 导入私钥
	granterAddr, _ := txClient.Key.Recover("test_key_name", "test_password", "supreme zero ladder chaos blur lake dinner warm rely voyage scan dilemma future spin victory glance legend faculty join man mansion water mansion exotic")
	granter, _ = types.AccAddressFromBech32(granterAddr)

	granteeAddr, _, _ := txClient.Key.Add("test_grantee", "12345678")
	grantee, _ = types.AccAddressFromBech32(granteeAddr)

}

//授权
func TestGrantAllowance(t *testing.T) {
	atom := types.NewCoins(types.NewInt64Coin("ugas", 55500000))
	threeHours := time.Now().Add(3 * time.Hour)
	basic := &feegrant.BasicAllowance{
		SpendLimit: atom,        //授权额度
		Expiration: &threeHours, //过期时间
	}
	// 授权
	result, err := txClient.Feegrant.GrantAllowance(granter, grantee, basic, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
	// sync 模式异步上链
	e := syncTx(result.Hash)
	require.NoError(t, e)
}

//设置交易代扣
func TestFeeGrant(t *testing.T) {
	testDenom := fmt.Sprintf("testdenom%d", rand.Int())
	resultTx, err := txClient.NFT.IssueDenom(nft.IssueDenomRequest{
		ID:   testDenom,
		Name: testDenom,
	}, baseTx)
	require.NoError(t, err)
	// sync 模式异步上链
	e := syncTx(resultTx.Hash)
	require.NoError(t, e)

	// 记录授权方现有余额
	granterAcc, err := txClient.Bank.QueryAccount(granter.String())
	require.NoError(t, err)
	originToken := granterAcc.Coins.AmountOf("ugas")

	// 限定此次交易费用
	feeNum := int64(200000)

	baseTx2 := types.BaseTx{
		From:       "test_grantee", // 由被授权方发起交易
		Password:   "12345678",
		Fee:        types.NewDecCoins(types.NewDecCoin("ugas", types.NewInt(feeNum))),
		Memo:       "",
		Mode:       types.Sync,
		FeeGranter: granter, //设置代扣地址
	}
	nftResult, err := txClient.NFT.MintNFT(nft.MintNFTRequest{
		Denom: testDenom,
		ID:    "testnftid",
	}, baseTx2)
	require.NoError(t, err)
	require.NotEmpty(t, nftResult.Hash)
	// sync 模式异步上链
	e = syncTx(nftResult.Hash)
	require.NoError(t, e)

	// 查询授权方账户现有余额
	newGranterAcc, err := txClient.Bank.QueryAccount(granter.String())
	require.NoError(t, err)
	newTokenNum := newGranterAcc.Coins.AmountOf("ugas")

	// 授权方账户应代付了200000ugas
	require.Equal(t, originToken.Add(newTokenNum.Neg()), types.NewInt(feeNum))
}

//解除授权
func TestRevokeAllowance(t *testing.T) {
	result, err := txClient.Feegrant.RevokeAllowance(granter, grantee, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
	// sync 模式异步上链
	e := syncTx(result.Hash)
	require.NoError(t, e)
}
