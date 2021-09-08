// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"

	grbacpb "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"os"
)

var DeleteSubjectInput grbacpb.DeleteSubjectRequest

var DeleteSubjectFromFile string

func init() {
	AccessControlServiceCmd.AddCommand(DeleteSubjectCmd)

	DeleteSubjectCmd.Flags().StringVar(&DeleteSubjectInput.Name, "name", "", "Required. The subject to delete.")

	DeleteSubjectCmd.Flags().StringVar(&DeleteSubjectFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteSubjectCmd = &cobra.Command{
	Use:   "delete-subject",
	Short: "DeleteSubject deletes a subject.",
	Long:  "DeleteSubject deletes a subject.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteSubjectFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteSubjectFromFile != "" {
			in, err = os.Open(DeleteSubjectFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteSubjectInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("AccessControl", "DeleteSubject", &DeleteSubjectInput)
		}
		err = AccessControlClient.DeleteSubject(ctx, &DeleteSubjectInput)

		return err
	},
}
