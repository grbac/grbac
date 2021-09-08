// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var CreateRoleInput grbacpb.CreateRoleRequest

var CreateRoleFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(CreateRoleCmd)

	CreateRoleInput.Role = new(grbacpb.Role)

	CreateRoleCmd.Flags().StringVar(&CreateRoleInput.Role.Name, "role.name", "", "Required. The resource name of the role.")

	CreateRoleCmd.Flags().StringSliceVar(&CreateRoleInput.Role.Permissions, "role.permissions", []string{}, "Required. The list of permissions granted by the role.")

	CreateRoleCmd.Flags().BytesHexVar(&CreateRoleInput.Role.Etag, "role.etag", []byte{}, "An etag for concurrency control, ignored during...")

	CreateRoleCmd.Flags().StringVar(&CreateRoleFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateRoleCmd = &cobra.Command{
	Use:   "create-role",
	Short: "CreateRole creates a new role.",
	Long:  "CreateRole creates a new role.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateRoleFromFile == "" {

			cmd.MarkFlagRequired("role.name")

			cmd.MarkFlagRequired("role.permissions")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateRoleFromFile != "" {
			in, err = os.Open(CreateRoleFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateRoleInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "CreateRole", &CreateRoleInput)
		}
		resp, err := AccessControlClient.CreateRole(ctx, &CreateRoleInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
