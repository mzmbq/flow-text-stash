package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single line",
			input: "\n\nHello World, \r\n",
			want:  "Hello World,",
		},
		{
			name:  "multiple lines",
			input: "Hello World, \r\n Second line\r\n",
			want:  "Hello World,",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, Wrap(tt.input))
	}
}
