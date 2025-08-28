package model

import (
	"context"
	"errors"
)

type AuthToken struct {
	projectKey string
	enableTLS  bool
	domain     string
}

func NewAuthToken(projectKey string) AuthToken {
	return AuthToken{
		projectKey: projectKey,
		enableTLS:  true,
	}
}

func (a *AuthToken) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{
		"projectKeyHeader": a.projectKey,
		"x-api-key":        a.projectKey,
	}, nil
}

func (a *AuthToken) RequireTransportSecurity() bool {
	return a.enableTLS
}

func (a *AuthToken) SetRequireTransportSecurity(enabled bool, domain string) error {
	if enabled && domain == "" {
		return errors.New("domain is required while tls is enabled")
	}
	a.enableTLS = enabled
	a.domain = domain
	return nil
}

func (a *AuthToken) GetProjectKey() string {
	return a.projectKey
}

func (a *AuthToken) GetEnableTLS() bool {
	return a.enableTLS
}

func (a *AuthToken) SetDomain(domain string) {
	a.domain = domain
}

func (a *AuthToken) GetDomain() string {
	return a.domain
}
