package rpc

import (
	"net"
	"os"
	"os/signal"
)

func (server *GrpcServer) Start() {
	go func() {
		lis, err := net.Listen("tcp", server.Port)
		if err != nil {
			server.Logger.Fatalw("failed to listen", "port", server.Port, "error", err)
		}
		server.Logger.Infow("grpc server listening", "port", server.Port)
		if err := server.Instance.Serve(lis); err != nil {
			server.Logger.Fatalw("failed to serve", "port", server.Port, "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	server.stop()
}

func (server *GrpcServer) stop() {
	server.Logger.Info("Shutting down gRPC server...")
	server.Instance.GracefulStop()
}
