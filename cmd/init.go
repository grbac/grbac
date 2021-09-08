package main

import (
	"context"

	"github.com/grbac/grbac/pkg/bootstrap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	type RuntimeConfig struct {
		dgraphEndpoint string
	}

	config := RuntimeConfig{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Runs the API server initializer",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := bootstrap.Schema(ctx, config.dgraphEndpoint); err != nil {
				logrus.Fatalf("failed to migrate the schema: %v", err)
			}

			logrus.Info("finished migrating the schema")
		},
	}

	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(
		&config.dgraphEndpoint,
		"dgraph-endpoint",
		"127.0.0.1:9080",
		"The endpoint of the dgraph database.")
}
