package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const defaultFilePermissions = 0644

// FileInfo represents the information available on a regular file.
type FileInfo struct {
	Name string // base name of the file
	Ext  string // file extension
	Dir  string // directory containing the file
	Path string // full file path
	Size int64  // file size in bytes
}

// MoveFileSafe moves the file with the given filename to the given destination.
// If the given destination already exists, MoveFileSafe prevents overwrites by
// trying other destinations in incrementing order, up to maxTries times.
// MoveFileSafe returns the destination to which the file is moved.
func MoveFileSafe(filename, destFilename string, maxTries int) (string, error) {
	if err := assertCopyable(filename, destFilename); err != nil {
		return "", err
	}

	nextFilename, err := nextFilename(destFilename, maxTries)
	if err != nil {
		return "", err
	}

	if err := moveFile(filename, nextFilename); err != nil {
		return "", err
	}

	return nextFilename, nil
}

// CopyFileSafe copies the file with the given filename to the given destination.
// If the given destination already exists, CopyFileSafe prevents overwrites by
// trying other destinations in incrementing order, up to maxTries times.
// CopyFileSafe returns the destination to which the file is copied.
func CopyFileSafe(filename, destFilename string, maxTries int) (string, error) {
	if err := assertCopyable(filename, destFilename); err != nil {
		return "", err
	}

	nextFilename, err := nextFilename(destFilename, maxTries)
	if err != nil {
		return "", err
	}

	if err := copyFile(filename, nextFilename); err != nil {
		return "", err
	}

	return nextFilename, nil
}

func nextFilename(filename string, maxTries int) (string, error) {
	nextFile, err := CreateNextFile(filename, maxTries)
	if err != nil {
		return "", err
	}
	if err := nextFile.Close(); err != nil {
		return "", err
	}

	return nextFile.Name(), nil
}

// CreateNextFile creates a file based on the given filename and returns it.
// If filename already exists, CreateNextFile inserts a counter in the filename
// and tries to create that file. The counter goes from 1 to maxTries included.
func CreateNextFile(filename string, maxTries int) (*os.File, error) {
	dir, name := filepath.Split(filepath.Clean(filename))

	for i := 0; i <= maxTries; i++ {
		nextFilename := filepath.Join(dir, insertCounter(name, i))
		nextFile, err := CreateFile(nextFilename)
		if err != nil {
			continue
		}
		return nextFile, nil
	}

	return nil, MaxTriesErr
}

// insertCounter inserts a counter in the given filename with the given value.
// The counter is inserted before the last dot, if any.
// For example, given filename "test.json" and a value of 3,
// insertCounter returns "test(3).json".
// If the value is less than 1, the original filename is returned.
func insertCounter(filename string, value int) string {
	if value < 1 {
		return filename
	}

	insertPos := strings.LastIndex(filename, ".")
	if insertPos == -1 {
		insertPos = len(filename)
	}

	counter := fmt.Sprintf("(%d)", value)
	newFilename := filename[:insertPos] + counter + filename[insertPos:]
	return newFilename
}

// CreateFile creates a file with the given filename and returns it.
// If filename already exists, CreateFile returns an error.
// CreateFile requires exclusive access to the given filename.
func CreateFile(filename string) (*os.File, error) {
	filename = filepath.Clean(filename)

	// Exclusive access to filename,
	// see https://golang.org/src/os/error_test.go
	// and https://stackoverflow.com/a/22483001
	file, err := os.OpenFile(
		filename,
		os.O_RDWR|os.O_CREATE|os.O_EXCL,
		defaultFilePermissions,
	)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// MoveFile moves the file with the given filename to the given destination.
// MoveFile overwrites existing destination files.
func MoveFile(filename, destFilename string) error {
	if err := assertCopyable(filename, destFilename); err != nil {
		return err
	}

	return moveFile(filename, destFilename)
}

func moveFile(filename, destFilename string) error {
	// Try a simple rename operation.
	moved := os.Rename(filename, destFilename) == nil
	if moved {
		return nil
	}

	// Otherwise, copy and remove.
	if err := copyFile(filename, destFilename); err != nil {
		return err
	}
	_ = RemoveFile(filename)
	return nil
}

// CopyFile copies the file with the given filename to the given destination.
// CopyFile overwrites existing destination files.
func CopyFile(filename, destFilename string) error {
	if err := assertCopyable(filename, destFilename); err != nil {
		return err
	}

	return copyFile(filename, destFilename)
}

func assertCopyable(filename, destFilename string) error {
	srcInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("%q is not a regular file", filename)
	}

	if strings.TrimSpace(destFilename) == "" {
		return DestFilenameEmptyErr
	}

	destInfo, _ := os.Stat(destFilename)
	destExistsAndNotRegular := destInfo != nil && !destInfo.Mode().IsRegular()
	if destExistsAndNotRegular {
		return fmt.Errorf("%q is not a regular file", destFilename)
	}

	sameFile := destInfo != nil && os.SameFile(srcInfo, destInfo)
	if sameFile {
		return SourceDestSameFileErr
	}

	return nil
}

func copyFile(filename, destFilename string) error {
	srcFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilename)
	if err != nil {
		return err
	}

	if _, err := io.Copy(destFile, srcFile); err != nil {
		_ = destFile.Close()
		return err
	}

	return destFile.Close()
}

// RemoveFile removes the file with the given filename.
func RemoveFile(filename string) error {
	return os.Remove(filename)
}

// ReadFileInfo returns the information available on a regular file.
func ReadFileInfo(filename string) (*FileInfo, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%q is not a regular file", filename)
	}

	name := info.Name()
	path := filepath.Clean(filename)
	fi := &FileInfo{
		Name: name,
		Ext:  filepath.Ext(name),
		Dir:  filepath.Dir(path),
		Path: path,
		Size: info.Size(),
	}
	return fi, nil
}

// AssertFile returns an error if the given filename is not a regular file.
func AssertFile(filename string) error {
	isRegular, err := IsFile(filename)
	if err != nil || !isRegular {
		return fmt.Errorf("%q is not  regular file", filename)
	}
	return nil
}

// IsFile returns true if the given filename represents a regular file.
func IsFile(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}
