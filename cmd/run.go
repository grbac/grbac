package main

import (
	"context"
	"os"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/grbac/grbac/pkg/graceful"
	"github.com/grbac/grbac/pkg/interrupt"
	"github.com/grbac/grbac/pkg/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RuntimeConfig struct {
	port           string
	dgraphEndpoint string
}

// TODO: Investigate whether mTLS could be useful.
// TODO: Investigate whether fallback server for HTTP/1.1 could be useful.

// See https://github.com/googleapis/gapic-showcase/blob/master/cmd/gapic-showcase/endpoint.go

func init() {
	config := RuntimeConfig{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the API server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			intr := interrupt.New(func(os.Signal) {}, cancel)

			opts := []grpc.ServerOption{}
			server := grpc.NewServer(opts...)

			cfg := &services.AccessControlServerConfig{
				DgraphHostname: config.dgraphEndpoint,
			}

			accessControlServer, err := services.NewAccessControlServer(cfg)
			if err != nil {
				logrus.WithError(err).Fatalf("failed to start the [authorizer] server")
			}
			defer accessControlServer.(*services.AccessControlServerImpl).Close()

			// Register Services to the server.
			grbac.RegisterAccessControlServer(server, accessControlServer)

			// Register reflection service on gRPC server.
			reflection.Register(server)

			if err := intr.Run(func() error { return graceful.NewGrpcListener(ctx, config.port, server) }); err != nil {
				logrus.WithError(err).Fatalf("http server exited with error")
			}
		},
	}

	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(
		&config.port,
		"port",
		"p",
		":9080",
		"The port that this serice will be served on.")

	runCmd.Flags().StringVar(
		&config.dgraphEndpoint,
		"dgraph-endpoint",
		"127.0.0.1:9080",
		"The endpoint of the dgraph database.")
}
