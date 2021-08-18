package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAuthToken(t *testing.T) {

	projectID := "TestProjectID"
	projectKey := "TestProjectKey"
	chainAccountAddress := "TestChainAccountAddress"

	authToken := NewAuthToken(projectID, projectKey, chainAccountAddress)

	metadata, err := authToken.GetRequestMetadata(context.Background())

	require.NoError(t, err)
	require.Equal(t, metadata["projectIdHeader"], projectID)
	require.Equal(t, metadata["projectKeyHeader"], projectKey)
	require.Equal(t, metadata["chainAccountAddressHeader"], chainAccountAddress)
	require.False(t, authToken.RequireTransportSecurity())
}
