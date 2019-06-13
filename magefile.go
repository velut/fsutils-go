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

// InstallDeps installs the dependencies.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("go", "mod", "download")
}
