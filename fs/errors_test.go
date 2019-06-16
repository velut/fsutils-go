package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr_Error(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		e    Err
		want string
	}{
		{
			"halt error",
			HaltErr,
			HaltErr.Error(),
		},
	}
	for _, tt := range tests {
		got := tt.e.Error()
		assert.Equal(tt.want, got, tt.name)
	}
}
