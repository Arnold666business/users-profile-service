package rpc

import (
	"context"
	"os"
	"time"
	userpb "users-profile-service/api/generation/users/v1/service"

	"users-profile-service/internal/api/rpc/users"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	Instance *grpc.Server
	Port     string
	Logger   *zap.SugaredLogger
}

func Build(Logger *zap.SugaredLogger, serviceImpl *users.UserService) *GrpcServer {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authInterceptor,
			loggingInterceptor(Logger),
		),
	)
	userpb.RegisterUserProfileServiceServer(grpcServer, serviceImpl)
	return &GrpcServer{Instance: grpcServer, Port: os.Getenv("GRPC_PORT"), Logger: Logger}
}

func loggingInterceptor(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		logger.Infow("gRPC request started",
			"method", info.FullMethod,
		)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			logger.Errorw("gRPC request failed",
				"method", info.FullMethod,
				"duration", duration,
				"error", err.Error(),
				"code", st.Code().String(),
			)
		} else {
			logger.Infow("gRPC request completed",
				"method", info.FullMethod,
				"duration", duration,
			)
		}

		return resp, err
	}
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "x-Token is not provided")
	}

	xToken := md["x-Token"]
	if len(xToken) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "x-Token is not provided")
	}

	if !validateXToken(xToken[0]) {
		return nil, status.Errorf(codes.Unauthenticated, "x-Token is not valid")
	}

	return handler(ctx, req)
}

func validateXToken(token string) bool {
	return os.Getenv("HTTP2_TOKEN") == token
}
