// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var DeletePermissionInput grbacpb.DeletePermissionRequest

var DeletePermissionFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(DeletePermissionCmd)

	DeletePermissionCmd.Flags().StringVar(&DeletePermissionInput.Name, "name", "", "Required. The resource name of the permission to delete.")

	DeletePermissionCmd.Flags().StringVar(&DeletePermissionFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeletePermissionCmd = &cobra.Command{
	Use:   "delete-permission",
	Short: "DeletePermission deletes a permission.",
	Long:  "DeletePermission deletes a permission.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeletePermissionFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeletePermissionFromFile != "" {
			in, err = os.Open(DeletePermissionFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeletePermissionInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "DeletePermission", &DeletePermissionInput)
		}
		err = AccessControlClient.DeletePermission(ctx, &DeletePermissionInput)

		return err
	},
}
