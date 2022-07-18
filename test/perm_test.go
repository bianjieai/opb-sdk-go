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
	res, err := txClient.Perm.AssignRoles(roleAddress, roles, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
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
	res, err := txClient.Perm.UnassignRoles(roleAddress, roles, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}

// block account
func TestPermBlockAccount(t *testing.T) {
	accountAddress := "iaa1ctagfms5nnn4r8tgvk8cy742jgecpvpnle2ktj"
	res, err := txClient.Perm.BlockAccount(accountAddress, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
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
	res, err := txClient.Perm.UnblockAccount(accountAddress, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, res.Hash)
	// sync 模式异步上链
	e := syncTx(res.Hash)
	require.NoError(t, e)
}
