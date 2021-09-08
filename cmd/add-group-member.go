// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var AddGroupMemberInput grbacpb.AddGroupMemberRequest

var AddGroupMemberFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(AddGroupMemberCmd)

	AddGroupMemberCmd.Flags().StringVar(&AddGroupMemberInput.Group, "group", "", "Required. The name of the group to add a member to.")

	AddGroupMemberCmd.Flags().StringVar(&AddGroupMemberInput.Member, "member", "", "Required. The member to be added.")

	AddGroupMemberCmd.Flags().StringVar(&AddGroupMemberFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var AddGroupMemberCmd = &cobra.Command{
	Use:   "add-group-member",
	Short: "AddGroupMember adds a member to a group.",
	Long:  "AddGroupMember adds a member to a group.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if AddGroupMemberFromFile == "" {

			cmd.MarkFlagRequired("group")

			cmd.MarkFlagRequired("member")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if AddGroupMemberFromFile != "" {
			in, err = os.Open(AddGroupMemberFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &AddGroupMemberInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "AddGroupMember", &AddGroupMemberInput)
		}
		resp, err := AccessControlClient.AddGroupMember(ctx, &AddGroupMemberInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
