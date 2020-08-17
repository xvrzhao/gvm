package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/funcs"
	e "github.com/xvrzhao/utils/errors"
)

var switchCmd = &cobra.Command{
	Use:     "switch SEMANTIC_VERSION",
	Aliases: []string{"s"},
	Short:   "Switch to the specified Go version",
	Long: `Switch to the specified Go version. You can add the flag '-i' to install 
the version if it's not installed yet, do not forget to add the flag 
'--cn' if you are in China and add '-i'.`,
	PreRun: isRootUser,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("need a version of Go to switch")
		}

		inCn, _ := cmd.Flags().GetBool("cn")
		v, err := funcs.NewVersion(args[0], inCn)
		if err != nil {
			return e.Wrapper(err, "new version error")
		}

		wantToInstall, _ := cmd.Flags().GetBool("install")
		if !v.IsInstalled() && wantToInstall {
			args := []string{"install", v.String()}
			if inCn {
				args = append(args, "--cn")
			}

			rootCmd.SetArgs(args)
			if err = rootCmd.Execute(); err != nil {
				return e.Wrapper(err, "install command executing error")
			}
		}

		fmt.Print("switching version ... ")

		if err := funcs.SwitchVersion(v); err != nil {
			return e.Wrapper(err, "switch version error")
		}

		fmt.Println("done")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().Bool("cn", false, "use https://golang.google.cn to download")
	switchCmd.Flags().BoolP("install", "i", false, "install if the version is not installed")
}
