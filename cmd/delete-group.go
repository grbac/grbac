// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var DeleteGroupInput grbacpb.DeleteGroupRequest

var DeleteGroupFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(DeleteGroupCmd)

	DeleteGroupCmd.Flags().StringVar(&DeleteGroupInput.Name, "name", "", "Required. The resource name of the group to delete.")

	DeleteGroupCmd.Flags().StringVar(&DeleteGroupFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteGroupCmd = &cobra.Command{
	Use:   "delete-group",
	Short: "DeleteGroup deletes a group.",
	Long:  "DeleteGroup deletes a group.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteGroupFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteGroupFromFile != "" {
			in, err = os.Open(DeleteGroupFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteGroupInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "DeleteGroup", &DeleteGroupInput)
		}
		err = AccessControlClient.DeleteGroup(ctx, &DeleteGroupInput)

		return err
	},
}
