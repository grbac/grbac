// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var UpdateRoleInput grbacpb.UpdateRoleRequest

var UpdateRoleFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(UpdateRoleCmd)

	UpdateRoleInput.Role = new(grbacpb.Role)

	UpdateRoleInput.UpdateMask = new(fieldmaskpb.FieldMask)

	UpdateRoleCmd.Flags().StringVar(&UpdateRoleInput.Role.Name, "role.name", "", "Required. The resource name of the role.")

	UpdateRoleCmd.Flags().StringSliceVar(&UpdateRoleInput.Role.Permissions, "role.permissions", []string{}, "Required. The list of permissions granted by the role.")

	UpdateRoleCmd.Flags().BytesHexVar(&UpdateRoleInput.Role.Etag, "role.etag", []byte{}, "An etag for concurrency control, ignored during...")

	UpdateRoleCmd.Flags().StringSliceVar(&UpdateRoleInput.UpdateMask.Paths, "update_mask.paths", []string{}, "The set of field mask paths.")

	UpdateRoleCmd.Flags().StringVar(&UpdateRoleFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateRoleCmd = &cobra.Command{
	Use:   "update-role",
	Short: "UpdateRole updates a role with a field mask.",
	Long:  "UpdateRole updates a role with a field mask.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateRoleFromFile == "" {

			cmd.MarkFlagRequired("role.name")

			cmd.MarkFlagRequired("role.permissions")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if UpdateRoleFromFile != "" {
			in, err = os.Open(UpdateRoleFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &UpdateRoleInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "UpdateRole", &UpdateRoleInput)
		}
		resp, err := AccessControlClient.UpdateRole(ctx, &UpdateRoleInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
