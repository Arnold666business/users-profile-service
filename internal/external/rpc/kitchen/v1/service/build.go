package service

import (
	"time"
	kitchenpb "users-profile-service/api/generation/kitchen/v1/service"
	"users-profile-service/internal/external/rpc/kitchen/v1/service/interceptor"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GrpcKitchenClient struct {
	conn   *grpc.ClientConn
	client kitchenpb.KitchenServiceClient
}

func Build(logger *zap.SugaredLogger) (*GrpcKitchenClient, error) {
	err := NewRpcServiceConfig()
	if err != nil {
		return nil, err
	}
	accessToken := cfgInstance.AccessToken
	serviceConfig := `{"loadBalancingPolicy": "least_connection"}`
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(5),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(150 * time.Millisecond)),
		grpc_retry.WithCodes(
			codes.Unavailable,
			codes.DeadlineExceeded,
			codes.ResourceExhausted,
		),
	}
	conn, err := grpc.NewClient(cfgInstance.Target, grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(retryOpts...),
			interceptor.AuthUnaryInterceptor(accessToken),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()), //todo
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                25 * time.Second,
			Timeout:             8 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &GrpcKitchenClient{conn, kitchenpb.NewKitchenServiceClient(conn)}, nil
}

func (g *GrpcKitchenClient) Close() error {
	return g.conn.Close()
}
