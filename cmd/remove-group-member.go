// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var RemoveGroupMemberInput grbacpb.RemoveGroupMemberRequest

var RemoveGroupMemberFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(RemoveGroupMemberCmd)

	RemoveGroupMemberCmd.Flags().StringVar(&RemoveGroupMemberInput.Group, "group", "", "Required. The name of the group to remove an member from.")

	RemoveGroupMemberCmd.Flags().StringVar(&RemoveGroupMemberInput.Member, "member", "", "Required. The member to be removed.")

	RemoveGroupMemberCmd.Flags().StringVar(&RemoveGroupMemberFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RemoveGroupMemberCmd = &cobra.Command{
	Use:   "remove-group-member",
	Short: "RemoveGroupMember removes a member from a group.",
	Long:  "RemoveGroupMember removes a member from a group.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RemoveGroupMemberFromFile == "" {

			cmd.MarkFlagRequired("group")

			cmd.MarkFlagRequired("member")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RemoveGroupMemberFromFile != "" {
			in, err = os.Open(RemoveGroupMemberFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RemoveGroupMemberInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "RemoveGroupMember", &RemoveGroupMemberInput)
		}
		resp, err := AccessControlClient.RemoveGroupMember(ctx, &RemoveGroupMemberInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
