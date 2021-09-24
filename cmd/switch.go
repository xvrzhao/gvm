package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
	e "github.com/xvrzhao/utils/errors"
)

var cmdSwitch = &cobra.Command{
	Use:     "switch SEMANTIC_VERSION",
	Aliases: []string{"s"},
	Short:   "Switch to the specified Go version",
	Long:    internal.CmdDescriptionSwitch,
	PreRun:  isRootUser,
	RunE:    runCmdSwitch,
}

func runCmdSwitch(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return errors.New("need a version of Go to switch")
	}

	inCn, _ := cmd.Flags().GetBool("cn")
	v, err := internal.NewVersion(args[0], inCn)
	if err != nil {
		return e.Wrapper(err, "new version error")
	}

	wantToInstall, _ := cmd.Flags().GetBool("install")
	if !v.IsInstalled() && wantToInstall {
		args := []string{"install", v.String()}
		if inCn {
			args = append(args, "--cn")
		}

		app.SetArgs(args)
		if err = app.Execute(); err != nil {
			return e.Wrapper(err, "install command executing error")
		}
	}

	fmt.Print("switching version ... ")

	if err := internal.SwitchVersion(v); err != nil {
		return e.Wrapper(err, "switch version error")
	}

	fmt.Println("done")
	return nil
}

func init() {
	cmdSwitch.Flags().Bool("cn", false, "use https://golang.google.cn to download")
	cmdSwitch.Flags().BoolP("install", "i", false, "install if the version is not installed")

	app.AddCommand(cmdSwitch)
}
