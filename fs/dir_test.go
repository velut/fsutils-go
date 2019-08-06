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
			"dir1, followed by a separator, is not a subdir of itself",
			args{
				dir1 + string(filepath.Separator),
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
			"dir2, followed by a separator, is not a subdir of itself",
			args{
				dir2 + string(filepath.Separator),
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
			"dir1, followed by a separator, is not a subdir of dir2",
			args{
				dir1 + string(filepath.Separator),
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
			"dir2, followed by a separator, is not a subdir of dir1",
			args{
				dir2 + string(filepath.Separator),
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
			"dir1, follwed by a separator, is a subdir of its direct parent",
			args{
				dir1 + string(filepath.Separator),
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
			"dir2, follwed by a separator, is a subdir of its direct parent",
			args{
				dir2 + string(filepath.Separator),
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
			"dir1, followed by a separator, is a subdir of its parent's parent",
			args{
				dir1 + string(filepath.Separator),
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
			"dir2, followed by a separator, is a subdir of its parent's parent",
			args{
				dir2 + string(filepath.Separator),
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
			"dir1, follwed by a separator, is a subdir of its root parent",
			args{
				dir1 + string(filepath.Separator),
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
		{
			"dir2, follwed by a separator, is a subdir of its root parent",
			args{
				dir2 + string(filepath.Separator),
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

func TestReadDir(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	file1, err := ioutil.TempFile("", "file")
	assert.Nil(err)
	file1.Close()
	defer os.Remove(file1.Name())

	dir2, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	dir3, err := ioutil.TempDir(dir2, "dir")
	assert.Nil(err)
	file3, err := ioutil.TempFile(dir3, "file")
	assert.Nil(err)
	file3.Close()
	defer os.RemoveAll(dir2)

	wd, err := os.Getwd()
	assert.Nil(err)
	testdir1 := filepath.Join(filepath.Dir(wd), "testdata", "read_dir_test")
	testdir1Contents := []*FileInfo{
		{
			Name: "10.gif",
			Ext:  ".gif",
			Dir:  testdir1,
			Path: filepath.Join(testdir1, "10.gif"),
			Size: 799,
		},
		{
			Name: "20.gif",
			Ext:  ".gif",
			Dir:  testdir1,
			Path: filepath.Join(testdir1, "20.gif"),
			Size: 799,
		},
		{
			Name: "30.gif",
			Ext:  ".gif",
			Dir:  filepath.Join(testdir1, "dir1"),
			Path: filepath.Join(testdir1, "dir1", "30.gif"),
			Size: 799,
		},
		{
			Name: "40.gif",
			Ext:  ".gif",
			Dir:  filepath.Join(testdir1, "dir1"),
			Path: filepath.Join(testdir1, "dir1", "40.gif"),
			Size: 799,
		},
		{
			Name: "50.gif",
			Ext:  ".gif",
			Dir:  filepath.Join(testdir1, "dir1", "subdir1"),
			Path: filepath.Join(testdir1, "dir1", "subdir1", "50.gif"),
			Size: 799,
		},
		{
			Name: "60.gif",
			Ext:  ".gif",
			Dir:  filepath.Join(testdir1, "dir1", "subdir1"),
			Path: filepath.Join(testdir1, "dir1", "subdir1", "60.gif"),
			Size: 799,
		},
		{
			Name: "70.gif",
			Ext:  ".gif",
			Dir:  filepath.Join(testdir1, "dir2"),
			Path: filepath.Join(testdir1, "dir2", "70.gif"),
			Size: 799,
		},
		{
			Name: "80.gif",
			Ext:  ".gif",
			Dir:  filepath.Join(testdir1, "dir2"),
			Path: filepath.Join(testdir1, "dir2", "80.gif"),
			Size: 799,
		},
	}

	type args struct {
		dirname string
		options *ReadDirOptions
	}
	tests := []struct {
		name    string
		args    args
		want    []*FileInfo
		wantErr bool
	}{
		{
			"invalid options",
			args{
				dir1,
				nil,
			},
			nil,
			true,
		},
		{
			"invalid dir",
			args{
				"",
				&ReadDirOptions{},
			},
			nil,
			true,
		},
		{
			"file is not a dir",
			args{
				file1.Name(),
				&ReadDirOptions{},
			},
			nil,
			true,
		},
		{
			"empty dir, exclude subdirs",
			args{
				dir1,
				&ReadDirOptions{
					IncludeSubdirs: false,
				},
			},
			[]*FileInfo{},
			false,
		},
		{
			"empty dir, include subdirs",
			args{
				dir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
				},
			},
			[]*FileInfo{},
			false,
		},
		{
			"empty dir, include subdirs, limit 100 files",
			args{
				dir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       100,
				},
			},
			[]*FileInfo{},
			false,
		},
		{
			"empty dir, non-empty subdirs, exclude subdirs",
			args{
				dir2,
				&ReadDirOptions{
					IncludeSubdirs: false,
				},
			},
			[]*FileInfo{},
			false,
		},
		{
			"empty dir, non-empty subdirs, include subdirs",
			args{
				dir2,
				&ReadDirOptions{
					IncludeSubdirs: true,
				},
			},
			[]*FileInfo{
				{
					Name: filepath.Base(file3.Name()),
					Ext:  filepath.Ext(file3.Name()),
					Dir:  dir3,
					Path: file3.Name(),
					Size: 0,
				},
			},
			false,
		},
		{
			"empty dir, non-empty subdirs, exclude subdirs, limit 1 file",
			args{
				dir2,
				&ReadDirOptions{
					IncludeSubdirs: false,
					MaxFiles:       1,
				},
			},
			[]*FileInfo{},
			false,
		},
		{
			"empty dir, non-empty subdirs, include subdirs, limit 1 file",
			args{
				dir2,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       1,
				},
			},
			[]*FileInfo{
				{
					Name: filepath.Base(file3.Name()),
					Ext:  filepath.Ext(file3.Name()),
					Dir:  dir3,
					Path: file3.Name(),
					Size: 0,
				},
			},
			false,
		},
		{
			"empty dir, non-empty subdirs, exclude subdirs, limit 100 files",
			args{
				dir2,
				&ReadDirOptions{
					IncludeSubdirs: false,
					MaxFiles:       100,
				},
			},
			[]*FileInfo{},
			false,
		},
		{
			"empty dir, non-empty subdirs, include subdirs, limit 100 files",
			args{
				dir2,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       100,
				},
			},
			[]*FileInfo{
				{
					Name: filepath.Base(file3.Name()),
					Ext:  filepath.Ext(file3.Name()),
					Dir:  dir3,
					Path: file3.Name(),
					Size: 0,
				},
			},
			false,
		},
		{
			"non-empty dir, exclude subdirs",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: false,
				},
			},
			testdir1Contents[:2],
			false,
		},
		{
			"non-empty dir, exclude subdirs, limit 1 file",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: false,
					MaxFiles:       1,
				},
			},
			testdir1Contents[:1],
			false,
		},
		{
			"non-empty dir, exclude subdirs, limit 2 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: false,
					MaxFiles:       2,
				},
			},
			testdir1Contents[:2],
			false,
		},
		{
			"non-empty dir, exclude subdirs, limit 100 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: false,
					MaxFiles:       100,
				},
			},
			testdir1Contents[:2],
			false,
		},
		{
			"non-empty dir, include subdirs",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
				},
			},
			testdir1Contents,
			false,
		},
		{
			"non-empty dir, include subdirs, limit 1 file",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       1,
				},
			},
			testdir1Contents[:1],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 2 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       2,
				},
			},
			testdir1Contents[:2],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 3 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       3,
				},
			},
			testdir1Contents[:3],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 4 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       4,
				},
			},
			testdir1Contents[:4],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 5 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       5,
				},
			},
			testdir1Contents[:5],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 6 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       6,
				},
			},
			testdir1Contents[:6],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 7 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       7,
				},
			},
			testdir1Contents[:7],
			false,
		},
		{
			"non-empty dir, include subdirs, limit 8 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       8,
				},
			},
			testdir1Contents,
			false,
		},
		{
			"non-empty dir, include subdirs, limit 100 files",
			args{
				testdir1,
				&ReadDirOptions{
					IncludeSubdirs: true,
					MaxFiles:       100,
				},
			},
			testdir1Contents,
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := ReadDir(tt.args.dirname, tt.args.options)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}
