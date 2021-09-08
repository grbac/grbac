// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var CreateGroupInput grbacpb.CreateGroupRequest

var CreateGroupFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(CreateGroupCmd)

	CreateGroupInput.Group = new(grbacpb.Group)

	CreateGroupCmd.Flags().StringVar(&CreateGroupInput.Group.Name, "group.name", "", "Required. The resource name of the group.")

	CreateGroupCmd.Flags().StringSliceVar(&CreateGroupInput.Group.Members, "group.members", []string{}, "The list of members of the group. Groups might...")

	CreateGroupCmd.Flags().BytesHexVar(&CreateGroupInput.Group.Etag, "group.etag", []byte{}, "An etag for concurrency control, ignored during...")

	CreateGroupCmd.Flags().StringVar(&CreateGroupFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateGroupCmd = &cobra.Command{
	Use:   "create-group",
	Short: "CreateGroup creates a new group.",
	Long:  "CreateGroup creates a new group.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateGroupFromFile == "" {

			cmd.MarkFlagRequired("group.name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateGroupFromFile != "" {
			in, err = os.Open(CreateGroupFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateGroupInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "CreateGroup", &CreateGroupInput)
		}
		resp, err := AccessControlClient.CreateGroup(ctx, &CreateGroupInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
