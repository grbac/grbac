// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var TestIamPolicyInput grbacpb.TestIamPolicyRequest

var TestIamPolicyFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(TestIamPolicyCmd)

	TestIamPolicyInput.AccessTuple = new(grbacpb.AccessTuple)

	TestIamPolicyCmd.Flags().StringVar(&TestIamPolicyInput.AccessTuple.Principal, "access_tuple.principal", "", "Required. The member, or principal, whose access you want...")

	TestIamPolicyCmd.Flags().StringVar(&TestIamPolicyInput.AccessTuple.FullResourceName, "access_tuple.full_resource_name", "", "Required. The full resource name that identifies the...")

	TestIamPolicyCmd.Flags().StringVar(&TestIamPolicyInput.AccessTuple.Permission, "access_tuple.permission", "", "Required. The IAM permission to check for the specified...")

	TestIamPolicyCmd.Flags().StringVar(&TestIamPolicyFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var TestIamPolicyCmd = &cobra.Command{
	Use:   "test-iam-policy",
	Short: "Checks whether a member has a specific permission...",
	Long:  "Checks whether a member has a specific permission for a specific resource.  If not allowed an Unauthorized (403) error will be returned.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if TestIamPolicyFromFile == "" {

			cmd.MarkFlagRequired("access_tuple.principal")

			cmd.MarkFlagRequired("access_tuple.full_resource_name")

			cmd.MarkFlagRequired("access_tuple.permission")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if TestIamPolicyFromFile != "" {
			in, err = os.Open(TestIamPolicyFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &TestIamPolicyInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "TestIamPolicy", &TestIamPolicyInput)
		}
		err = AccessControlClient.TestIamPolicy(ctx, &TestIamPolicyInput)

		return err
	},
}
