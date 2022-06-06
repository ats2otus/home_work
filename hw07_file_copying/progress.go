package main

import (
	"fmt"
	"io"
	"strings"
)

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
