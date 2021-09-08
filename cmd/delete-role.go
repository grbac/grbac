// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var DeleteRoleInput grbacpb.DeleteRoleRequest

var DeleteRoleFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(DeleteRoleCmd)

	DeleteRoleCmd.Flags().StringVar(&DeleteRoleInput.Name, "name", "", "Required. The resource name of the role to delete.")

	DeleteRoleCmd.Flags().StringVar(&DeleteRoleFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteRoleCmd = &cobra.Command{
	Use:   "delete-role",
	Short: "DeleteRole deletes a role.",
	Long:  "DeleteRole deletes a role.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteRoleFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteRoleFromFile != "" {
			in, err = os.Open(DeleteRoleFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteRoleInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "DeleteRole", &DeleteRoleInput)
		}
		err = AccessControlClient.DeleteRole(ctx, &DeleteRoleInput)

		return err
	},
}
