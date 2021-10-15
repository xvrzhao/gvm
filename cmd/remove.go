package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var cmdRemove = &cobra.Command{
	Use:     "remove VERSION [VERSIONS]",
	Aliases: []string{"rm", "uninstall", "ui", "delete", "del"},
	Short:   "Remove versions managed by GVM",

	PreRun:  checkPermission,
	RunE:    runCmdRemove,
	PostRun: printDone,
}

func runCmdRemove(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return internal.ErrNoVersionSpecified
	}

	curVersion, err := internal.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to GetCurrentVersion: %w", err)
	}

	versions := make([]*internal.Version, 0)
	for _, versionName := range args {
		v, err := internal.NewVersion(versionName, false)
		if err != nil {
			return fmt.Errorf("failed to NewVersion: %w", err)
		}

		if !v.IsInstalled() {
			continue
		}

		if v.String() == curVersion {
			return internal.ErrVersionIsInUse
		}

		versions = append(versions, v)
	}

	for _, v := range versions {
		if err = v.Remove(); err != nil {
			return fmt.Errorf("failed to remove go%v: %w", v, err)
		}
	}

	return nil
}

func init() {
	App.AddCommand(cmdRemove)
}
