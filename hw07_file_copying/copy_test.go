package main

import (
	"math"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	root, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(root) })

	testCases := []struct {
		name     string
		limit    int64
		offset   int64
		hasErr   bool
		src, dst string
	}{
		{
			name:   "negative offset",
			limit:  -1,
			hasErr: true,
		},
		{
			name:   "negative offset",
			offset: -1,
			hasErr: true,
		},
		{
			name:   "too big offset",
			offset: math.MaxInt,
			hasErr: true,
			src:    "testdata/input.txt",
		},
		{
			name:   "absence src file",
			hasErr: true,
			src:    "testdata/not_found.txt",
		},
		{
			name:   "sizeless src file",
			hasErr: true,
			src:    "/dev/urandom",
		},
		{
			name:   "unreachable dst file path",
			hasErr: true,
			src:    "testdata/input.txt",
			dst:    "/tmp/hw07/file_copying/dst.txt",
		},
		{
			name:  "with limit",
			limit: 10,
			src:   "testdata/input.txt",
			dst:   path.Join(root, "out.limit.txt"),
		},
		{
			name:   "with offset",
			offset: 10,
			src:    "testdata/input.txt",
			dst:    path.Join(root, "out.offset.txt"),
		},
		{
			name:  "with limit and offset",
			limit: 10, offset: 10,
			src: "testdata/input.txt",
			dst: path.Join(root, "out.limit.offset.txt"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.src, tc.dst, tc.offset, tc.limit)
			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.Equal(
					t,
					mustReadFileData(tc.src, tc.offset, tc.limit),
					mustReadFileData(tc.dst, 0, 0),
				)
			}
		})
	}
}

func mustReadFileData(path string, from, size int64) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return string(data[from:])
	}
	return string(data[from : from+size])
}
