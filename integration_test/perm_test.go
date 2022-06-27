package integration_test

import (
	"github.com/irisnet/core-sdk-go/types"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/iritamod-sdk-go/perm"
)

func (s IntegrationTestSuite) TestPerm() {
	baseTx := types.BaseTx{
		From:          s.Account().Name,
		Password:      s.Account().Password,
		Gas:           200000,
		Mode:          types.Commit,
		GasAdjustment: 1.5,
	}

	acc := s.GetRandAccount()
	roles := []perm.Role{
		perm.RoleBlacklistAdmin,
	}

	// add role
	rs, err := s.Perm.AssignRoles(acc.Address.String(), roles, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// query role added
	roles2, err := s.Perm.QueryRoles(acc.Address.String())
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), roles2)
	require.EqualValues(s.T(), roles, roles2)

	// remove role
	rs, err = s.Perm.UnassignRoles(acc.Address.String(), roles, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// query role removed
	roles2, err = s.Perm.QueryRoles(acc.Address.String())
	require.NoError(s.T(), err)
	require.Empty(s.T(), roles2)

	// block account
	rs, err = s.Perm.BlockAccount(acc.Address.String(), baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// query blacklist
	bl, err := s.Perm.QueryAccountBlockList()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), bl)
	require.Contains(s.T(), bl, acc.Address.String())

	// unblock blacklist
	rs, err = s.Perm.UnblockAccount(acc.Address.String(), baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// query blacklist
	bl, err = s.Perm.QueryAccountBlockList()
	require.NoError(s.T(), err)
	require.NotContains(s.T(), bl, acc.Address.String())

	// contract ??
}
