// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var DeleteResourceInput grbacpb.DeleteResourceRequest

var DeleteResourceFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(DeleteResourceCmd)

	DeleteResourceCmd.Flags().StringVar(&DeleteResourceInput.Name, "name", "", "Required. The full resource name that identifies the...")

	DeleteResourceCmd.Flags().StringVar(&DeleteResourceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteResourceCmd = &cobra.Command{
	Use:   "delete-resource",
	Short: "DeleteResource deletes a resource.",
	Long:  "DeleteResource deletes a resource.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteResourceFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteResourceFromFile != "" {
			in, err = os.Open(DeleteResourceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteResourceInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "DeleteResource", &DeleteResourceInput)
		}
		err = AccessControlClient.DeleteResource(ctx, &DeleteResourceInput)

		return err
	},
}
