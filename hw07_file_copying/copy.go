package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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
		return fmt.Errorf("src:stat: is empty")
	case offset > fileSize:
		return fmt.Errorf("offset[%d] is bigger than file size[%d]", offset, fileSize)
	}
	if _, err := ff.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	ft, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer ft.Close()

	pw := ProgressWriter{
		dst: ft, log: os.Stdout,
		total: min(limit, fileSize-offset),
		head:  "▓", tail: "░", width: 50,
	}

	if _, err := io.CopyN(&pw, ff, limit); err != nil {
		if err != io.EOF {
			return err
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

type ProgressWriter struct {
	size       int
	total      int64
	width      int
	head, tail string
	dst, log   io.Writer
}

func (pw *ProgressWriter) Write(data []byte) (int, error) {
	write, err := pw.dst.Write(data)
	pw.size += write

	percent := float64(pw.size*100) / float64(pw.total)
	middlePoint := int(percent) * pw.width / 100

	head := strings.Repeat(pw.head, middlePoint)
	tail := strings.Repeat(pw.tail, pw.width-middlePoint)
	fmt.Fprintf(pw.log, "\r%.2f%%\t[%s%s]", percent, head, tail)

	return write, err
}
