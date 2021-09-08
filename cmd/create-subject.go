// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var CreateSubjectInput grbacpb.CreateSubjectRequest

var CreateSubjectFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(CreateSubjectCmd)

	CreateSubjectInput.Subject = new(grbacpb.Subject)

	CreateSubjectCmd.Flags().StringVar(&CreateSubjectInput.Subject.Name, "subject.name", "", "Required. The resource name of the subject.")

	CreateSubjectCmd.Flags().StringVar(&CreateSubjectFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateSubjectCmd = &cobra.Command{
	Use:   "create-subject",
	Short: "CreateSubject creates a new subject.",
	Long:  "CreateSubject creates a new subject.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateSubjectFromFile == "" {

			cmd.MarkFlagRequired("subject.name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateSubjectFromFile != "" {
			in, err = os.Open(CreateSubjectFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateSubjectInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "CreateSubject", &CreateSubjectInput)
		}
		resp, err := AccessControlClient.CreateSubject(ctx, &CreateSubjectInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
