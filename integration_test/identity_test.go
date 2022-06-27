package integration_test

import (
	iritaidentity "github.com/bianjieai/iritamod-sdk-go/identity"

	"github.com/irisnet/core-sdk-go/common/uuid"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/sm2"
)

func (s IntegrationTestSuite) TestIdentity() {
	baseTx := sdk.BaseTx{
		From:          s.Account().Name,
		Password:      s.Account().Password,
		Gas:           200000,
		Mode:          sdk.Commit,
		GasAdjustment: 1.5,
	}

	var (
		uuidGenerator, _ = uuid.NewV4()
		test1PubKeySm2   = sm2.GenPrivKey().PubKeySm2()
		test2PubKeySm2   = sm2.GenPrivKey().PubKeySm2()
		testCredentials  = "https://kyc.com/user/10001"
		testCertificate  = ""
	)

	var (
		id = sdk.HexStringFrom(uuidGenerator.Bytes())

		pubKeyInfo1 = iritaidentity.PubKeyInfo{
			PubKey:    sdk.HexStringFrom(test1PubKeySm2[:]),
			Algorithm: iritaidentity.SM2,
		}

		pubKeyInfo2 = iritaidentity.PubKeyInfo{
			PubKey:    sdk.HexStringFrom(test2PubKeySm2[:]),
			Algorithm: iritaidentity.SM2,
		}

		// request for creation
		req1 = iritaidentity.CreateIdentityRequest{
			Id:          id,
			PubKeyInfo:  &pubKeyInfo1,
			Certificate: testCertificate,
			Credentials: &testCredentials,
		}

		// request for update
		req2 = iritaidentity.UpdateIdentityRequest{
			Id:          id,
			PubKeyInfo:  &pubKeyInfo2,
			Certificate: testCertificate,
			Credentials: &testCredentials,
		}
	)

	res1, err := s.Identity.QueryIdentity(id)
	require.Empty(s.T(), res1)
	require.Error(s.T(), err)

	//create and query
	res2, err := s.Identity.CreateIdentity(req1, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res2.Hash)

	res3, err := s.Identity.QueryIdentity(id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), res3.Credentials, testCredentials)
	require.Contains(s.T(), res3.PubKeyInfos, pubKeyInfo1)

	// update and query
	res4, err := s.Identity.UpdateIdentity(req2, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res4.Hash)

	res5, err := s.Identity.QueryIdentity(id)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res5)
	require.Contains(s.T(), res5.PubKeyInfos, pubKeyInfo2)
}
