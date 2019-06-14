package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		dirname string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"invalid dir",
			args{
				"",
			},
			false,
			true,
		},
		{
			"file is not a dir",
			args{
				file1.Name(),
			},
			false,
			false,
		},
		{
			"valid dir",
			args{
				dir1,
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := IsDir(tt.args.dirname)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestAssertDir(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		dirname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"invalid dir",
			args{
				"",
			},
			true,
		},
		{
			"file is not a dir",
			args{
				file1.Name(),
			},
			true,
		},
		{
			"valid dir",
			args{
				dir1,
			},
			false,
		},
	}
	for _, tt := range tests {
		gotErr := AssertDir(tt.args.dirname)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func TestSameDir(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir2)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		dirname1 string
		dirname2 string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"invalid dirname1",
			args{
				"",
				dir2,
			},
			false,
			true,
		},
		{
			"invalid dirname2",
			args{
				dir1,
				"",
			},
			false,
			true,
		},
		{
			"file is not a dirname1",
			args{
				file1.Name(),
				dir2,
			},
			false,
			true,
		},
		{
			"file is not a dirname2",
			args{
				dir1,
				file1.Name(),
			},
			false,
			true,
		},
		{
			"different dirs",
			args{
				dir1,
				dir2,
			},
			false,
			false,
		},
		{
			"same dir1",
			args{
				dir1,
				dir1,
			},
			true,
			false,
		},
		{
			"same dir2",
			args{
				dir2,
				dir2,
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := SameDir(tt.args.dirname1, tt.args.dirname2)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestSubdirOf(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir2)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	type args struct {
		dirname    string
		targetname string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"invalid dirname",
			args{
				"",
				dir2,
			},
			false,
			true,
		},
		{
			"invalid targetname",
			args{
				dir1,
				"",
			},
			false,
			true,
		},
		{
			"file is not a dirname",
			args{
				file1.Name(),
				dir2,
			},
			false,
			true,
		},
		{
			"file is not a targetname",
			args{
				dir1,
				file1.Name(),
			},
			false,
			true,
		},
		{
			"dir1 is not a subdir of itself",
			args{
				dir1,
				dir1,
			},
			false,
			false,
		},
		{
			"dir2 is not a subdir of itself",
			args{
				dir2,
				dir2,
			},
			false,
			false,
		},
		{
			"dir1 is not a subdir of dir2",
			args{
				dir1,
				dir2,
			},
			false,
			false,
		},
		{
			"dir2 is not a subdir of dir1",
			args{
				dir2,
				dir1,
			},
			false,
			false,
		},
		{
			"dir1 is a subdir of its direct parent",
			args{
				dir1,
				filepath.Dir(dir1),
			},
			true,
			false,
		},
		{
			"dir2 is a subdir of its direct parent",
			args{
				dir2,
				filepath.Dir(dir2),
			},
			true,
			false,
		},
		{
			"dir1 is a subdir of its parent's parent",
			args{
				dir1,
				filepath.Dir(filepath.Dir(dir1)),
			},
			true,
			false,
		},
		{
			"dir2 is a subdir of its parent's parent",
			args{
				dir2,
				filepath.Dir(filepath.Dir(dir2)),
			},
			true,
			false,
		},
		{
			"dir1 is a subdir of its root parent",
			args{
				dir1,
				rootParent(dir1),
			},
			true,
			false,
		},
		{
			"dir2 is a subdir of its root parent",
			args{
				dir2,
				rootParent(dir2),
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := SubdirOf(tt.args.dirname, tt.args.targetname)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func rootParent(dirname string) string {
	prevParent := filepath.Clean(dirname)
	nextParent := filepath.Dir(dirname)

	for {
		if nextParent == prevParent {
			return nextParent
		}

		prevParent = nextParent
		nextParent = filepath.Dir(nextParent)
	}
}