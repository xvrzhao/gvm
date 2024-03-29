package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/xvrzhao/gvm/cmd"
	"github.com/xvrzhao/gvm/internal"
)

func main() {
	handleError(cmd.App.Execute())
}

func handleError(err error) {
	if err == nil {
		return
	}

	var msg string

	var userError internal.UserError
	if errors.As(err, &userError) {
		msg = fmt.Sprintf("Usage Error: %v\n", userError)
	} else if strings.Contains(err.Error(), "unknown command") {
		msg = fmt.Sprintf("Usage Error: %v\n", err.Error())
	} else {
		msg = fmt.Sprintf("GVM Internal Error: %v\n", err)
	}

	os.Stderr.WriteString(msg)
	os.Exit(1)
}
