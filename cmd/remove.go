/*
Copyright © 2020 Xavier Zhao <xvrzhao@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	e "github.com/xvrzhao/utils/errors"
	"gvm/funcs"
)

var removeCmd = &cobra.Command{
	Use:     "remove SEMANTIC_VERSION [SEMANTIC_VERSION...]",
	Aliases: []string{"rm", "uninstall"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  isRootUser,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("no version to delete")
		}

		curVersion, err := funcs.GetCurrentVersionStr()
		if err != nil {
			return e.Wrapper(err, "GetCurrentVersionStr error")
		}

		versions := make([]*funcs.Version, 0)
		for _, semVerStr := range args {
			v, err := funcs.NewVersion(semVerStr, false)
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
			if err = funcs.RmVersion(v); err != nil {
				return e.Wrapper(err, "RmVersion error")
			}
		}

		fmt.Println("done")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
