// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var GetRoleInput grbacpb.GetRoleRequest

var GetRoleFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(GetRoleCmd)

	GetRoleCmd.Flags().StringVar(&GetRoleInput.Name, "name", "", "Required. The name of the role to retrieve.")

	GetRoleCmd.Flags().StringVar(&GetRoleFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetRoleCmd = &cobra.Command{
	Use:   "get-role",
	Short: "GetRole returns a role.",
	Long:  "GetRole returns a role.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetRoleFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetRoleFromFile != "" {
			in, err = os.Open(GetRoleFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetRoleInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "GetRole", &GetRoleInput)
		}
		resp, err := AccessControlClient.GetRole(ctx, &GetRoleInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
