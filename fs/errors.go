package fs

// Err represents an error.
type Err string

// Error implements the error interface.
func (e Err) Error() string {
	return string(e)
}

// HaltErr is the error used to halt walks when reading directories.
const HaltErr = Err("fs-halt")
