package integration_test

import (
	"fmt"
	iritaidentity "github.com/bianjieai/iritamod-sdk-go/identity"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/sm2"
	"testing"
)

func TestQueryIdentity(t *testing.T) {
	res, err := txClient.Identity.QueryIdentity("Test_ID")
	if err != nil {
		fmt.Println(fmt.Errorf("Identity 查询失败: %s", err.Error()))
		return
	} else {
		fmt.Println("Identity 查询成功：", res)
	}
}

func TestCreateIdentidy(t *testing.T) {
	testPubKeySm2 := sm2.GenPrivKey().PubKeySm2()
	pubKeyInfo := iritaidentity.PubKeyInfo{
		PubKey:    sdk.HexStringFrom(testPubKeySm2[:]),
		Algorithm: iritaidentity.SM2,
	}

	// request for creation
	Credentials := "https://kyc.com/user/10001"
	req := iritaidentity.CreateIdentityRequest{
		Id:          "Test_ID",
		PubKeyInfo:  &pubKeyInfo,
		Certificate: "",
		Credentials: &Credentials,
	}

	res, err := txClient.Identity.CreateIdentity(req, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}

func TestUpdateIdentity(t *testing.T) {
	testPubKeySm2 := sm2.GenPrivKey().PubKeySm2()
	pubKeyInfo := iritaidentity.PubKeyInfo{
		PubKey:    sdk.HexStringFrom(testPubKeySm2[:]),
		Algorithm: iritaidentity.SM2,
	}

	// request for creation
	Credentials := "https://kyc.com/user/10001"
	req := iritaidentity.UpdateIdentityRequest{
		Id:          "Test_ID",
		PubKeyInfo:  &pubKeyInfo,
		Certificate: "",
		Credentials: &Credentials,
	}

	res, err := txClient.Identity.UpdateIdentity(req, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
}
