// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test tests the fs library.
func Test() error {
	mg.Deps(InstallDeps)

	fmt.Println("Testing...")
	return sh.Run("go", "test", "-v", "-race", "./fs")
}

// Lint lints the fs library.
func Lint() error {
	fmt.Println("Linting...")
	return sh.Run("golangci-lint", "run", "./fs")
}

// LintAll lints the fs library using all linters.
func LintAll() error {
	fmt.Println("Linting...")
	return sh.Run("golangci-lint", "run", "--enable-all", "-D", "goimports", "-D", "gofmt", "./fs")
}

// InstallDeps installs the dependencies.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("go", "mod", "download")
}
