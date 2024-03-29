package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
)

// ReadDirOptions represents the options available for reading a directory.
type ReadDirOptions struct {
	// IncludeSubdirs, if true, specifies that subdirectories,
	// at any nesting level, should also be read.
	IncludeSubdirs bool

	// MaxFiles specifies the maximum number of files to read.
	// If MaxFiles is 0, all files are read.
	MaxFiles int
}

// ReadDir reads the directory named by the given dirname
// following the given options and returns a list of FileInfo instances,
// sorted by lexical path order, representing the regular files found.
// Eventual filesystem errors are ignored.
func ReadDir(dirname string, options *ReadDirOptions) ([]*FileInfo, error) {
	if options == nil {
		return nil, NoReadDirOptionsErr
	}

	if err := AssertDir(dirname); err != nil {
		return nil, err
	}

	return readDir(dirname, options)
}

func readDir(dirname string, options *ReadDirOptions) ([]*FileInfo, error) {
	dirname = filepath.Clean(dirname)
	skipSubdirs := !options.IncludeSubdirs
	maxFiles := options.MaxFiles
	limitFiles := maxFiles > 0

	fileInfos := make([]*FileInfo, 0, 1000)
	_ = godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			osPathname = filepath.Clean(osPathname)

			skipDir := skipSubdirs && de.IsDir() && osPathname != dirname
			if skipDir {
				return filepath.SkipDir
			}

			if de.IsRegular() {
				fi, _ := ReadFileInfo(osPathname)
				if fi != nil {
					fileInfos = append(fileInfos, fi)
				}
			}

			halt := limitFiles && len(fileInfos) >= maxFiles
			if halt {
				return HaltErr
			}

			return nil
		},
		ErrorCallback: func(_ string, err error) godirwalk.ErrorAction {
			if err == HaltErr {
				return godirwalk.Halt
			}

			return godirwalk.SkipNode
		},
	})

	return fileInfos, nil
}

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

	for {
		nextParent := filepath.Dir(prevParent)

		rootReached := nextParent == prevParent
		if rootReached {
			return false, nil
		}

		nextParentInfo, _ := os.Stat(nextParent)
		if os.SameFile(nextParentInfo, targetInfo) {
			return true, nil
		}

		prevParent = nextParent
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
