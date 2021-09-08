// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	iampb "google.golang.org/genproto/googleapis/iam/v1"

	"os"
)

var SetIamPolicyInput iampb.SetIamPolicyRequest

var SetIamPolicyFromFile string

var SetIamPolicyInputPolicyBindings []string

func init() {
	AccessControlServiceCmd.AddCommand(SetIamPolicyCmd)

	SetIamPolicyInput.Policy = new(iampb.Policy)

	SetIamPolicyCmd.Flags().StringVar(&SetIamPolicyInput.Resource, "resource", "", "Required. REQUIRED: The resource for which the policy is...")

	SetIamPolicyCmd.Flags().Int32Var(&SetIamPolicyInput.Policy.Version, "policy.version", 0, "Specifies the format of the policy.   Valid...")

	SetIamPolicyCmd.Flags().StringArrayVar(&SetIamPolicyInputPolicyBindings, "policy.bindings", []string{}, "Associates a list of `members` to a `role`....")

	SetIamPolicyCmd.Flags().BytesHexVar(&SetIamPolicyInput.Policy.Etag, "policy.etag", []byte{}, "`etag` is used for optimistic concurrency control...")

	SetIamPolicyCmd.Flags().StringVar(&SetIamPolicyFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var SetIamPolicyCmd = &cobra.Command{
	Use:   "set-iam-policy",
	Short: "Sets the IAM policy that is attached to a generic...",
	Long:  "Sets the IAM policy that is attached to a generic resource.  Note: the full resource name that identifies the resource must be provided.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if SetIamPolicyFromFile == "" {

			cmd.MarkFlagRequired("resource")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if SetIamPolicyFromFile != "" {
			in, err = os.Open(SetIamPolicyFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &SetIamPolicyInput)
			if err != nil {
				return err
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range SetIamPolicyInputPolicyBindings {
			tmp := iampb.Binding{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			SetIamPolicyInput.Policy.Bindings = append(SetIamPolicyInput.Policy.Bindings, &tmp)
		}

		if Verbose {
			printVerboseInput("AccessControl", "SetIamPolicy", &SetIamPolicyInput)
		}
		resp, err := AccessControlClient.SetIamPolicy(ctx, &SetIamPolicyInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
