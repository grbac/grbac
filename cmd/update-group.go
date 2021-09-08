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

var UpdateGroupInput grbacpb.UpdateGroupRequest

var UpdateGroupFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(UpdateGroupCmd)

	UpdateGroupInput.Group = new(grbacpb.Group)

	UpdateGroupInput.UpdateMask = new(fieldmaskpb.FieldMask)

	UpdateGroupCmd.Flags().StringVar(&UpdateGroupInput.Group.Name, "group.name", "", "Required. The resource name of the group.")

	UpdateGroupCmd.Flags().StringSliceVar(&UpdateGroupInput.Group.Members, "group.members", []string{}, "The list of members of the group. Groups might...")

	UpdateGroupCmd.Flags().BytesHexVar(&UpdateGroupInput.Group.Etag, "group.etag", []byte{}, "An etag for concurrency control, ignored during...")

	UpdateGroupCmd.Flags().StringSliceVar(&UpdateGroupInput.UpdateMask.Paths, "update_mask.paths", []string{}, "The set of field mask paths.")

	UpdateGroupCmd.Flags().StringVar(&UpdateGroupFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateGroupCmd = &cobra.Command{
	Use:   "update-group",
	Short: "UpdateGroup updates a group with a field mask.",
	Long:  "UpdateGroup updates a group with a field mask.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateGroupFromFile == "" {

			cmd.MarkFlagRequired("group.name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if UpdateGroupFromFile != "" {
			in, err = os.Open(UpdateGroupFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &UpdateGroupInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "UpdateGroup", &UpdateGroupInput)
		}
		resp, err := AccessControlClient.UpdateGroup(ctx, &UpdateGroupInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
