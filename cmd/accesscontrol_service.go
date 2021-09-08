// Code generated. DO NOT EDIT.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	gapic "github.com/animeapis/api-go-client/grbac/v1alpha1"
)

var AccessControlConfig *viper.Viper
var AccessControlClient *gapic.AccessControlClient
var AccessControlSubCommands []string = []string{
	"test-iam-policy",
	"get-iam-policy",
	"set-iam-policy",
	"get-resource",
	"create-resource",
	"transfer-resource",
	"delete-resource",
	"create-subject",
	"delete-subject",
	"get-group",
	"create-group",
	"update-group",
	"add-group-member",
	"remove-group-member",
	"delete-group",
	"create-permission",
	"delete-permission",
	"get-role",
	"create-role",
	"update-role",
	"delete-role",
}

func init() {
	rootCmd.AddCommand(AccessControlServiceCmd)

	AccessControlConfig = viper.New()
	AccessControlConfig.SetEnvPrefix("GRBAC_ACCESSCONTROL")
	AccessControlConfig.AutomaticEnv()

	AccessControlServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use GRBAC_ACCESSCONTROL_INSECURE. Must be used with \"address\" option")
	AccessControlConfig.BindPFlag("insecure", AccessControlServiceCmd.PersistentFlags().Lookup("insecure"))
	AccessControlConfig.BindEnv("insecure")

	AccessControlServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use GRBAC_ACCESSCONTROL_ADDRESS.")
	AccessControlConfig.BindPFlag("address", AccessControlServiceCmd.PersistentFlags().Lookup("address"))
	AccessControlConfig.BindEnv("address")

	AccessControlServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use GRBAC_ACCESSCONTROL_TOKEN.")
	AccessControlConfig.BindPFlag("token", AccessControlServiceCmd.PersistentFlags().Lookup("token"))
	AccessControlConfig.BindEnv("token")

	AccessControlServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use GRBAC_ACCESSCONTROL_API_KEY.")
	AccessControlConfig.BindPFlag("api_key", AccessControlServiceCmd.PersistentFlags().Lookup("api_key"))
	AccessControlConfig.BindEnv("api_key")
}

var AccessControlServiceCmd = &cobra.Command{
	Use:       "accesscontrol",
	Short:     "AccessControl is the internal service used by...",
	Long:      "AccessControl is the internal service used by Animeshon to enforce RBAC rules.",
	ValidArgs: AccessControlSubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := AccessControlConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if AccessControlConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := AccessControlConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := AccessControlConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		AccessControlClient, err = gapic.NewAccessControlClient(ctx, opts...)
		return
	},
}
