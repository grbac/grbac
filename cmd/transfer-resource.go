// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"

	"strings"
)

var TransferResourceInput grbacpb.TransferResourceRequest

var TransferResourceFromFile string

var TransferResourceInputSubstitutions []string

func init() {
	AccessControlServiceCmd.AddCommand(TransferResourceCmd)

	TransferResourceCmd.Flags().StringVar(&TransferResourceInput.Name, "name", "", "Required. The full resource name that identifies the...")

	TransferResourceCmd.Flags().StringVar(&TransferResourceInput.TargetParent, "target_parent", "", "Required. The full resource name that identifies the new...")

	TransferResourceCmd.Flags().StringArrayVar(&TransferResourceInputSubstitutions, "substitutions", []string{}, "key=value pairs. The map of substitutions to apply to the full...")

	TransferResourceCmd.Flags().StringVar(&TransferResourceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var TransferResourceCmd = &cobra.Command{
	Use:   "transfer-resource",
	Short: "TransferResource transfers a resource to a new...",
	Long:  "TransferResource transfers a resource to a new parent.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if TransferResourceFromFile == "" {

			cmd.MarkFlagRequired("name")

			cmd.MarkFlagRequired("target_parent")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if TransferResourceFromFile != "" {
			in, err = os.Open(TransferResourceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &TransferResourceInput)
			if err != nil {
				return err
			}

		}

		if len(TransferResourceInputSubstitutions) > 0 {
			TransferResourceInput.Substitutions = make(map[string]string)
		}
		for _, item := range TransferResourceInputSubstitutions {
			split := strings.Split(item, "=")
			if len(split) < 2 {
				err = fmt.Errorf("Invalid map item: %q", item)
				return
			}

			TransferResourceInput.Substitutions[split[0]] = split[1]
		}

		if Verbose {
			printVerboseInput("AccessControl", "TransferResource", &TransferResourceInput)
		}
		resp, err := AccessControlClient.TransferResource(ctx, &TransferResourceInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
