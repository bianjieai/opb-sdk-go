package test

import (
	"fmt"
	"github.com/bianjieai/iritamod-sdk-go/perm"
	"github.com/stretchr/testify/require"
	"testing"
)

// add role
func TestPermAssignRoles(t *testing.T) {
	roleAddress := "iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj"
	roles := []perm.Role{
		perm.RoleBlacklistAdmin,
	}

	// add role
	rs, err := txClient.Perm.AssignRoles(roleAddress, roles, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, rs.Hash)
}

// query role
func TestPermQueryRoles(t *testing.T) {
	roleAddress := "iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj"
	role, err := txClient.Perm.QueryRoles(roleAddress)
	require.NoError(t, err)
	fmt.Println(role)
}

// remove role
func TestPermUnassignRoles(t *testing.T) {
	roleAddress := "iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj"
	roles := []perm.Role{
		perm.RoleBlacklistAdmin,
	}
	rs, err := txClient.Perm.UnassignRoles(roleAddress, roles, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, rs.Hash)
}

// block account
func TestPermBlockAccount(t *testing.T) {
	accountAddress := "iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj"
	rs, err := txClient.Perm.BlockAccount(accountAddress, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, rs.Hash)
}

// query blacklist
func TestPermQueryAccountBlockList(t *testing.T) {
	bl, err := txClient.Perm.QueryAccountBlockList()
	require.NoError(t, err)
	fmt.Println(bl)
}

// unblock blacklist
func TestPermUnblockAccount(t *testing.T) {
	accountAddress := "iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj"
	rs, err := txClient.Perm.UnblockAccount(accountAddress, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, rs.Hash)
}
