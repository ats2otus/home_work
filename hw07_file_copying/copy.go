package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 {
		return fmt.Errorf("invalid value of limit %d", limit)
	}
	if offset < 0 {
		return fmt.Errorf("invalid value of offset %d", offset)
	}
	if limit == 0 {
		limit = math.MaxInt64
	}

	ff, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("src:open: %v", err)
	}
	defer ff.Close()

	info, err := ff.Stat()
	if err != nil {
		return fmt.Errorf("src:stat: %v", err)
	}

	fileSize := info.Size()
	switch {
	case fileSize == 0:
		return ErrUnsupportedFile
	case offset > fileSize:
		return ErrOffsetExceedsFileSize
	}
	if _, err := ff.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	ft, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("dst:create: %v", err)
	}
	defer ft.Close()

	pw := ProgressWriter{
		dst: ft, log: os.Stdout,
		total: min(limit, fileSize-offset),
		head:  "▓", tail: "░", width: 50,
	}

	if _, err := io.CopyN(&pw, ff, limit); err != nil {
		if err != io.EOF {
			return fmt.Errorf("copy: %v", err)
		}
	}

	return nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
