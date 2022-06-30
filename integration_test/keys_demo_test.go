package integration_test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeys(t *testing.T) {
	name, password := "test_name", "12345678"

	address, mnemonic, err := txClient.Key.Add(name, password)
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotEmpty(t, mnemonic)

	address1, err := txClient.Key.Show(name, password)
	require.NoError(t, err)
	require.Equal(t, address, address1)

	privKeyArmor, err := txClient.Key.Export(name, password)
	require.NoError(t, err)

	err = txClient.Key.Delete(name, password)
	require.NoError(t, err)

	address2, err := txClient.Key.Import(name, password, privKeyArmor)
	require.NoError(t, err)
	require.Equal(t, address, address2)

	err = txClient.Key.Delete(name, password)
	require.NoError(t, err)

	// test Recover
	//address3, err := s.Key.Recover(name, password, mnemonic)
	//require.NoError(t, err)
	//require.Equal(t, address, address3)

	// test Recover With HD Path
	address4, err := txClient.Key.RecoverWithHDPath(name, password, mnemonic, "")
	require.NoError(t, err)
	require.Equal(t, address, address4)
}
