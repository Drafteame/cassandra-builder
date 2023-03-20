// nolint
package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Lint Runs golangci-lint checks over the code.
func Lint() error {
	command := "revive"
	args := []string{"-config=revive.toml", "-formatter=friendly", "-exclude=magefiles/...", "./..."}

	out, err := sh.Output(command, args...)

	fmt.Println(out)
	return err
}

// Format Runs gofmt over the code.
func Format() error {
	outImp, err := sh.Output("goimports-reviser", "-format", "./...")
	if err != nil {
		return err
	}

	fmt.Println(outImp)

	return nil
}
