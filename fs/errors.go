package fs

// Err represents an error.
type Err string

// Error implements the error interface.
func (e Err) Error() string {
	return string(e)
}

// HaltErr is the error used to halt walks when reading directories.
const HaltErr = Err("fs: halt")

// NoReadDirOptionsErr is the error returned when no options are given to ReadDir.
const NoReadDirOptionsErr = Err("fs: no options specified for ReadDir")

// MaxTriesErr is the error returned when an operation fails
// after exceeding the maximum number of tries.
const MaxTriesErr = Err("fs: exceeded maximum number of tries")

// DestFilenameEmptyErr is the error returned when the filename of a destination file is empty.
const DestFilenameEmptyErr = Err("fs: destination filename cannot be empty")

// SourceDestSameFileErr is the error returned when the source and destination files coincide.
const SourceDestSameFileErr = Err("fs: source and destination are the same file")
