package sdk

import (
	iritasdk "github.com/bianjieai/irita-sdk-go"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/opb-sdk-go/internal/app/sdk/model"
	"google.golang.org/grpc"
)

// NewClient create a new IRITA OPB client
func NewClient(cfg types.ClientConfig, authToken *model.AuthToken) iritasdk.IRITAClient {

	// overwrite grpcOpts
	grpcOpts := []grpc.DialOption {
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authToken),
	}
	cfg.GRPCOptions = grpcOpts

	return iritasdk.NewIRITAClient(cfg)
}