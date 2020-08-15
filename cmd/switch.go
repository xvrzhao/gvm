package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/funcs"
	e "github.com/xvrzhao/utils/errors"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch SEMANTIC_VERSION",
	Aliases: []string{"s"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  isRootUser,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	switchCmd.Flags().Bool("cn", false, "Use https://golang.google.cn to download.")
	switchCmd.Flags().BoolP("install", "i", false, "Install if the version is not installed.")
}
