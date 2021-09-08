package graceful

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// NewGrpcListener listens for incoming gRPC requests.
func NewGrpcListener(ctx context.Context, address string, server *grpc.Server) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	logrus.Infof("gRPC server listening to [%s]", address)
	return ServeWithContext(ctx, server, listener)
}

// ServeWithContext is a wrapper around the Serve function which also implements
// context cancellation and graceful shutdown.
func ServeWithContext(ctx context.Context, server *grpc.Server, listener net.Listener) error {
	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gr := new(errgroup.Group)
	gr.Go(func() error {
		defer cancel()

		if err := server.Serve(listener); err != nil {
			return err
		}

		return nil
	})

	gr.Go(func() error {
		for {
			select {
			case <-serverCtx.Done():
				// ListenAndServe exited already - nothing to do.
				return nil
			case <-ctx.Done():
				// SIGTERM or SIGINT received - initiate graceful shutdown.
				goto shutdown
			}
		}

	shutdown:
		logrus.Info("gracefully shutting down the server - waiting for active connections to close")
		server.GracefulStop()

		return nil
	})

	return gr.Wait()
}
