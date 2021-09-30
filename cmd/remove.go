package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
	e "github.com/xvrzhao/utils/errors"
)

var cmdRemove = &cobra.Command{
	Use:     "remove SEMANTIC_VERSION [SEMANTIC_VERSION...]",
	Aliases: []string{"rm", "uninstall", "ui", "delete", "del"},
	Short:   "Remove one or more Go versions installed by GVM",
	Long:    `Remove one or more Go versions installed by GVM.`,
	PreRun:  isRootUser,
	RunE:    runCmdRemove,
}

func runCmdRemove(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return errors.New("no version to delete")
	}

	curVersion, err := internal.GetCurrentVersionStr()
	if err != nil {
		return e.Wrapper(err, "GetCurrentVersionStr error")
	}

	versions := make([]*internal.Version, 0)
	for _, semVerStr := range args {
		v, err := internal.NewVersion(semVerStr, false)
		if err != nil {
			return e.Wrapper(err, "error when new version %s", semVerStr)
		}

		if !v.IsInstalled() {
			return fmt.Errorf("go%s is not installed", v.String())
		}

		if v.String() == curVersion {
			return errors.New("can not remove current version")
		}

		versions = append(versions, v)
	}

	fmt.Print("remove versions ... ")

	for _, v := range versions {
		if err = internal.RmVersion(v); err != nil {
			return e.Wrapper(err, "RmVersion error")
		}
	}

	fmt.Println("done")
	return nil
}

func init() {
	App.AddCommand(cmdRemove)
}
