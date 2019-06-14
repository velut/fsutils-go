package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// SubdirOf returns true if the given dirname is a subdirectory
// of the given target directory.
func SubdirOf(dirname, targetname string) (bool, error) {
	dirInfo, err := os.Stat(dirname)
	if err != nil {
		return false, err
	}
	if !dirInfo.IsDir() {
		return false, fmt.Errorf("%q is not a directory", dirname)
	}

	targetInfo, err := os.Stat(targetname)
	if err != nil {
		return false, err
	}
	if !targetInfo.IsDir() {
		return false, fmt.Errorf("%q is not a directory", targetname)
	}

	if os.SameFile(dirInfo, targetInfo) {
		return false, nil
	}

	prevParent := filepath.Clean(dirname)
	nextParent := filepath.Dir(dirname)

	for {
		if nextParent == prevParent {
			return false, nil
		}

		nextParentInfo, _ := os.Stat(nextParent)
		if os.SameFile(nextParentInfo, targetInfo) {
			return true, nil
		}

		prevParent = nextParent
		nextParent = filepath.Dir(nextParent)
	}
}

// SameDir returns true if the given dirnames point to the same directory.
func SameDir(dirname1, dirname2 string) (bool, error) {
	info1, err := os.Stat(dirname1)
	if err != nil {
		return false, err
	}
	if !info1.IsDir() {
		return false, fmt.Errorf("%q is not a directory", dirname1)
	}

	info2, err := os.Stat(dirname2)
	if err != nil {
		return false, err
	}
	if !info2.IsDir() {
		return false, fmt.Errorf("%q is not a directory", dirname2)
	}

	return os.SameFile(info1, info2), nil
}

// AssertDir returns an error if the given dirname is not a directory.
func AssertDir(dirname string) error {
	isDir, err := IsDir(dirname)
	if err != nil || !isDir {
		return fmt.Errorf("%q is not a directory", dirname)
	}
	return nil
}

// IsDir returns true if the given dirname represents a directory.
func IsDir(dirname string) (bool, error) {
	info, err := os.Stat(dirname)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
