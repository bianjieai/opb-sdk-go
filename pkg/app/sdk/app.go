package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"strings"

	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/client"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"github.com/irisnet/core-sdk-go/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewClient create a new IRITA OPB client
func NewClient(cfg types.ClientConfig, authToken *model.AuthToken) client.Client {

	httpHeader := http.Header{}
	if authToken != nil {
		if authToken.GetEnableTLS() {
			certificateList, err := getGateWayTlsCertPool(cfg.RPCAddr)
			if err != nil {
				panic(err)
			}
			roots := x509.NewCertPool()
			for i := range certificateList {
				roots.AddCert(certificateList[i])
			}
			cert := credentials.NewClientTLSFromCert(roots, "bsngate.com")
			// overwrite grpcOpts
			grpcOpts := []grpc.DialOption{
				grpc.WithPerRPCCredentials(authToken),
				grpc.WithTransportCredentials(cert),
			}
			cfg.GRPCOptions = grpcOpts
		} else {
			// overwrite grpcOpts
			grpcOpts := []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(authToken),
			}
			cfg.GRPCOptions = grpcOpts
		}

		if projectKey := authToken.GetProjectKey(); projectKey != "" {
			httpHeader.Set("x-api-key", authToken.GetProjectKey())
		}
	}

	return client.NewClient(cfg)
}

func getGateWayTlsCertPool(gateWayUrl string) ([]*x509.Certificate, error) {

	if !strings.Contains(strings.ToLower(gateWayUrl), "https://") {
		return nil, errors.New("tls is enabled, but the address is http")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(gateWayUrl)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return resp.TLS.PeerCertificates, nil

}
