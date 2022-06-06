package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProgressWriter(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		size     int
		total    int
	}{
		{
			name:     "0%",
			expected: "\r0.00%\t[----------]",
			size:     0,
			total:    100,
		},
		{
			name:     "9%",
			expected: "\r9.00%\t[----------]",
			size:     9,
			total:    100,
		},
		{
			name:     "50%",
			expected: "\r50.00%\t[+++++-----]",
			size:     50,
			total:    100,
		},
		{
			name:     "100%",
			expected: "\r100.00%\t[++++++++++]",
			size:     100,
			total:    100,
		},
		{
			name:     "99%",
			expected: "\r99.00%\t[+++++++++-]",
			size:     99,
			total:    100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buff bytes.Buffer
			pw := ProgressWriter{
				head: "+", tail: "-",
				dst: io.Discard, log: &buff,
				width: 10, total: int64(tc.total),
			}
			written, err := pw.Write(make([]byte, tc.size))
			assert.NoError(t, err)
			assert.Equal(t, tc.size, written)
			assert.Equal(t, tc.expected, buff.String())
		})
	}
}
