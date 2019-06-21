package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"non-existing file",
			args{
				"",
			},
			false,
			true,
		},
		{
			"directory is not a regular file",
			args{
				dir1,
			},
			false,
			false,
		},
		{
			"regular file",
			args{
				file1.Name(),
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := IsFile(tt.args.filename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestAssertFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"non-existing file",
			args{
				"",
			},
			true,
		},
		{
			"directory is not a regular file",
			args{
				dir1,
			},
			true,
		},
		{
			"regular file",
			args{
				file1.Name(),
			},
			false,
		},
	}
	for _, tt := range tests {
		gotErr := AssertFile(tt.args.filename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func TestReadFileInfo(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file*.txt")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *FileInfo
		wantErr bool
	}{
		{
			"non-existing file",
			args{
				"",
			},
			nil,
			true,
		},
		{
			"directory is not a regular file",
			args{
				dir1,
			},
			nil,
			true,
		},
		{
			"regular file",
			args{
				file1.Name(),
			},
			&FileInfo{
				Name: filepath.Base(file1.Name()),
				Ext:  filepath.Ext(file1.Name()),
				Dir:  filepath.Dir(file1.Name()),
				Path: file1.Name(),
				Size: 0,
			},
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := ReadFileInfo(tt.args.filename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestCreateFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file*.txt")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		wantFile bool
		wantErr  bool
	}{
		{
			"directory already exists",
			args{
				dir1,
			},
			false,
			true,
		},
		{
			"file already exists",
			args{
				file1.Name(),
			},
			false,
			true,
		},
		{
			"filename is available",
			args{
				filepath.Join(dir1, "file"),
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := CreateFile(tt.args.filename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.wantFile, got != nil, tt.name)

		if gotErr != nil {
			got.Close()
		}
	}
}

func TestRemoveFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1, err := ioutil.TempFile(dir1, "file*.txt")
	assert.Nil(err)
	file1.Close()
	defer os.RemoveAll(dir1)

	file2, err := ioutil.TempFile("", "file*.txt")
	assert.Nil(err)
	file2.Close()
	defer os.Remove(file2.Name())

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid file",
			args{
				"",
			},
			true,
		},
		{
			"non-empty directory",
			args{
				dir1,
			},
			true,
		},
		{
			"removable file",
			args{
				file2.Name(),
			},
			false,
		},
	}
	for _, tt := range tests {
		gotErr := RemoveFile(tt.args.filename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func TestCreateNextFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1, err := ioutil.TempFile(dir1, "file*.txt")
	assert.Nil(err)
	file1.Close()
	file2, err := ioutil.TempFile(dir1, "file*")
	assert.Nil(err)
	file2.Close()
	file3, err := ioutil.TempFile(dir1, "file*.json.txt")
	assert.Nil(err)
	file3.Close()
	defer os.RemoveAll(dir1)

	type args struct {
		filename string
		maxTries int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"existing file with extension, zero tries",
			args{
				file1.Name(),
				0,
			},
			"",
			true,
		},
		{
			"existing file with extension, one try",
			args{
				file1.Name(),
				1,
			},
			strings.TrimSuffix(file1.Name(), ".txt") + "(1)" + ".txt",
			false,
		},
		{
			"existing file without extension, zero tries",
			args{
				file2.Name(),
				0,
			},
			"",
			true,
		},
		{
			"existing file without extension, one try",
			args{
				file2.Name(),
				1,
			},
			file2.Name() + "(1)",
			false,
		},
		{
			"existing file with extension and other dot, zero tries",
			args{
				file3.Name(),
				0,
			},
			"",
			true,
		},
		{
			"existing file with extension and other dot, one try",
			args{
				file3.Name(),
				1,
			},
			strings.TrimSuffix(file3.Name(), ".txt") + "(1)" + ".txt",
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := CreateNextFile(tt.args.filename, tt.args.maxTries)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		if got != nil {
			assert.Equal(tt.want, got.Name(), tt.name)
		}
	}
}

func TestCopyFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1Name := filepath.Join(dir1, "file1.txt")
	err = ioutil.WriteFile(
		file1Name,
		[]byte("hello world"),
		defaultFilePermissions,
	)
	assert.Nil(err)
	file2, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	file2.Close()
	file3, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	defer file3.Close()
	file4, err := CreateFile(filepath.Join(dir1, "file4"))
	assert.Nil(err)
	defer file4.Close()
	defer os.RemoveAll(dir1)

	type args struct {
		filename     string
		destFilename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid filename",
			args{
				"",
				file2.Name(),
			},
			true,
		},
		{
			"invalid copy from directory",
			args{
				dir1,
				file2.Name(),
			},
			true,
		},
		{
			"invalid copy to directory",
			args{
				file1Name,
				dir1,
			},
			true,
		},
		{
			"invalid copy to same file",
			args{
				file1Name,
				file1Name,
			},
			true,
		},
		{
			"invalid copy to invalid file",
			args{
				file1Name,
				"",
			},
			true,
		},
		{
			"copy to another file",
			args{
				file1Name,
				file2.Name(),
			},
			false,
		},
		{
			"copy to open file",
			args{
				file1Name,
				file3.Name(),
			},
			false,
		},
		{
			"copy to exclusive file",
			args{
				file1Name,
				file4.Name(),
			},
			false,
		},
	}
	for _, tt := range tests {
		gotErr := CopyFile(tt.args.filename, tt.args.destFilename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func TestMoveFile(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1Name := filepath.Join(dir1, "file1.txt")
	err = ioutil.WriteFile(
		file1Name,
		[]byte("hello world"),
		defaultFilePermissions,
	)
	assert.Nil(err)
	file2, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	file2.Close()
	file3, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	defer file3.Close()
	file4, err := CreateFile(filepath.Join(dir1, "file4"))
	assert.Nil(err)
	defer file4.Close()
	defer os.RemoveAll(dir1)

	type args struct {
		filename     string
		destFilename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid filename",
			args{
				"",
				file2.Name(),
			},
			true,
		},
		{
			"invalid move from directory",
			args{
				dir1,
				file2.Name(),
			},
			true,
		},
		{
			"invalid move to directory",
			args{
				file1Name,
				dir1,
			},
			true,
		},
		{
			"invalid move to same file",
			args{
				file1Name,
				file1Name,
			},
			true,
		},
		{
			"invalid move to invalid file",
			args{
				file1Name,
				"",
			},
			true,
		},
		{
			"move to another file",
			args{
				file1Name,
				file2.Name(),
			},
			false,
		},
		{
			"move to open file",
			args{
				file2.Name(),
				file3.Name(),
			},
			false,
		},
		{
			"move to exclusive file",
			args{
				file3.Name(),
				file4.Name(),
			},
			false,
		},
	}
	for _, tt := range tests {
		gotErr := MoveFile(tt.args.filename, tt.args.destFilename)
		if tt.wantErr != (gotErr != nil) {
			fmt.Println()
			fmt.Println(tt.args.filename)
			fmt.Println(tt.args.destFilename)
			fmt.Println(gotErr)
			fmt.Println()
		}
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func TestCopyFileSafe(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1Name := filepath.Join(dir1, "file1.txt")
	err = ioutil.WriteFile(
		file1Name,
		[]byte("hello world"),
		defaultFilePermissions,
	)
	assert.Nil(err)
	file2, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	file2.Close()
	file3, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	defer file3.Close()
	file4, err := CreateFile(filepath.Join(dir1, "file4"))
	assert.Nil(err)
	defer file4.Close()
	defer os.RemoveAll(dir1)

	type args struct {
		filename     string
		destFilename string
		maxTries     int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"invalid filename",
			args{
				"",
				file1Name,
				1,
			},
			"",
			true,
		},
		{
			"invalid copy from directory",
			args{
				dir1,
				file2.Name(),
				1,
			},
			"",
			true,
		},
		{
			"invalid copy to dir",
			args{
				file1Name,
				dir1,
				1,
			},
			"",
			true,
		},
		{
			"invalid copy to same file",
			args{
				file1Name,
				file1Name,
				1,
			},
			"",
			true,
		},
		{
			"invalid copy to invalid file",
			args{
				file1Name,
				"",
				1,
			},
			"",
			true,
		},
		{
			"existing destination, zero maxTries",
			args{
				file1Name,
				file2.Name(),
				0,
			},
			"",
			true,
		},
		{
			"existing destination, one maxTries",
			args{
				file1Name,
				file2.Name(),
				1,
			},
			file2.Name() + "(1)",
			false,
		},
		{
			"copy to open file",
			args{
				file1Name,
				file3.Name(),
				1,
			},
			file3.Name() + "(1)",
			false,
		},
		{
			"copy to exclusive file",
			args{
				file1Name,
				file4.Name(),
				1,
			},
			file4.Name() + "(1)",
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := CopyFileSafe(tt.args.filename, tt.args.destFilename, tt.args.maxTries)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestMoveFileSafe(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1Name := filepath.Join(dir1, "file1.txt")
	err = ioutil.WriteFile(
		file1Name,
		[]byte("hello world"),
		defaultFilePermissions,
	)
	assert.Nil(err)
	file2, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	file2.Close()
	file3, err := ioutil.TempFile(dir1, "file")
	assert.Nil(err)
	defer file3.Close()
	file4, err := CreateFile(filepath.Join(dir1, "file4"))
	assert.Nil(err)
	defer file4.Close()
	defer os.RemoveAll(dir1)

	type args struct {
		filename     string
		destFilename string
		maxTries     int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"invalid filename",
			args{
				"",
				file1Name,
				1,
			},
			"",
			true,
		},
		{
			"invalid move from directory",
			args{
				dir1,
				file2.Name(),
				1,
			},
			"",
			true,
		},
		{
			"invalid move to dir",
			args{
				file1Name,
				dir1,
				1,
			},
			"",
			true,
		},
		{
			"invalid move to same file",
			args{
				file1Name,
				file1Name,
				1,
			},
			"",
			true,
		},
		{
			"invalid move to invalid file (1)",
			args{
				file1Name,
				"",
				1,
			},
			"",
			true,
		},
		{
			"invalid move to invalid file (2)",
			args{
				file1Name,
				"   ",
				1,
			},
			"",
			true,
		},
		{
			"existing destination, zero maxTries",
			args{
				file1Name,
				file2.Name(),
				0,
			},
			"",
			true,
		},
		{
			"existing destination, one maxTries",
			args{
				file1Name,
				file2.Name(),
				1,
			},
			file2.Name() + "(1)",
			false,
		},
		{
			"move to open file",
			args{
				file2.Name(),
				file3.Name(),
				1,
			},
			file3.Name() + "(1)",
			false,
		},
		{
			"move to exclusive file",
			args{
				file3.Name(),
				file4.Name(),
				1,
			},
			file4.Name() + "(1)",
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := MoveFileSafe(tt.args.filename, tt.args.destFilename, tt.args.maxTries)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func Test_copyFile(t *testing.T) {
	assert := assert.New(t)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		filename     string
		destFilename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid filename",
			args{
				"",
				file1.Name(),
			},
			true,
		},
		{
			"invalid destFilename",
			args{
				file1.Name(),
				"",
			},
			true,
		},
	}
	for _, tt := range tests {
		gotErr := copyFile(tt.args.filename, tt.args.destFilename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func Test_moveFile(t *testing.T) {
	assert := assert.New(t)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		filename     string
		destFilename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid filename",
			args{
				"",
				file1.Name(),
			},
			true,
		},
		{
			"invalid destFilename",
			args{
				file1.Name(),
				"",
			},
			true,
		},
	}
	for _, tt := range tests {
		gotErr := moveFile(tt.args.filename, tt.args.destFilename)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}
