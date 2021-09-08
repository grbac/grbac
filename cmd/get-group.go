// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var GetGroupInput grbacpb.GetGroupRequest

var GetGroupFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(GetGroupCmd)

	GetGroupCmd.Flags().StringVar(&GetGroupInput.Name, "name", "", "Required. The name of the group to retrieve.")

	GetGroupCmd.Flags().StringVar(&GetGroupFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetGroupCmd = &cobra.Command{
	Use:   "get-group",
	Short: "GetGroup returns a group.",
	Long:  "GetGroup returns a group.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetGroupFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetGroupFromFile != "" {
			in, err = os.Open(GetGroupFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetGroupInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "GetGroup", &GetGroupInput)
		}
		resp, err := AccessControlClient.GetGroup(ctx, &GetGroupInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
