// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var CreateResourceInput grbacpb.CreateResourceRequest

var CreateResourceFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(CreateResourceCmd)

	CreateResourceInput.Resource = new(grbacpb.Resource)

	CreateResourceCmd.Flags().StringVar(&CreateResourceInput.Resource.Name, "resource.name", "", "Required. The full resource name that identifies the...")

	CreateResourceCmd.Flags().StringVar(&CreateResourceInput.Resource.Parent, "resource.parent", "", "Required. The full resource name that identifies the parent...")

	CreateResourceCmd.Flags().BytesHexVar(&CreateResourceInput.Resource.Etag, "resource.etag", []byte{}, "An etag for concurrency control, ignored during...")

	CreateResourceCmd.Flags().StringVar(&CreateResourceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateResourceCmd = &cobra.Command{
	Use:   "create-resource",
	Short: "CreateResource creates a new resource.",
	Long:  "CreateResource creates a new resource.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateResourceFromFile == "" {

			cmd.MarkFlagRequired("resource.name")

			cmd.MarkFlagRequired("resource.parent")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateResourceFromFile != "" {
			in, err = os.Open(CreateResourceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateResourceInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "CreateResource", &CreateResourceInput)
		}
		resp, err := AccessControlClient.CreateResource(ctx, &CreateResourceInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
