package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	iritasdk "github.com/bianjieai/irita-sdk-go"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/opb-sdk-go/pkg/app/sdk/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"strings"
)

// NewClient create a new IRITA OPB client
func NewClient(cfg types.ClientConfig, authToken *model.AuthToken) iritasdk.IRITAClient {

	httpHeader := http.Header{}
	if authToken != nil {
		if authToken.GetEnableTLS() {
			certificateList,err := getGateWayTlsCertPool(cfg.RPCAddr)
			if err != nil {
				panic(err)
			}
			roots := x509.NewCertPool()
			for i ,_ := range certificateList {
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
			grpcOpts := []grpc.DialOption {
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(authToken),
			}
			cfg.GRPCOptions = grpcOpts
		}

		if projectKey := authToken.GetProjectKey(); projectKey != "" {
			httpHeader.Set("x-api-key", authToken.GetProjectKey())
		}
	}

	cfg.RPCHeader = httpHeader
	cfg.WSHeader = httpHeader

	return iritasdk.NewIRITAClient(cfg)
}

func getGateWayTlsCertPool(gateWayUrl string) ([]*x509.Certificate,error) {

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
