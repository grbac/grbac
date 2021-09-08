// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var GetResourceInput grbacpb.GetResourceRequest

var GetResourceFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(GetResourceCmd)

	GetResourceCmd.Flags().StringVar(&GetResourceInput.Name, "name", "", "Required. The full resource name of the resource to...")

	GetResourceCmd.Flags().StringVar(&GetResourceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetResourceCmd = &cobra.Command{
	Use:   "get-resource",
	Short: "GetResource returns a resource.",
	Long:  "GetResource returns a resource.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetResourceFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetResourceFromFile != "" {
			in, err = os.Open(GetResourceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetResourceInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "GetResource", &GetResourceInput)
		}
		resp, err := AccessControlClient.GetResource(ctx, &GetResourceInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
