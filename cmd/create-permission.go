// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var CreatePermissionInput grbacpb.CreatePermissionRequest

var CreatePermissionFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(CreatePermissionCmd)

	CreatePermissionInput.Permission = new(grbacpb.Permission)

	CreatePermissionCmd.Flags().StringVar(&CreatePermissionInput.Permission.Name, "permission.name", "", "Required. The resource name of the permission.")

	CreatePermissionCmd.Flags().StringVar(&CreatePermissionFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreatePermissionCmd = &cobra.Command{
	Use:   "create-permission",
	Short: "CreatePermission creates a new permission.",
	Long:  "CreatePermission creates a new permission.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreatePermissionFromFile == "" {

			cmd.MarkFlagRequired("permission.name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreatePermissionFromFile != "" {
			in, err = os.Open(CreatePermissionFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreatePermissionInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "CreatePermission", &CreatePermissionInput)
		}
		resp, err := AccessControlClient.CreatePermission(ctx, &CreatePermissionInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
