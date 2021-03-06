package hw02unpackstring

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(src string) (string, error) {
	if isInvalidStr(src) {
		return "", ErrInvalidString
	}

	var prev rune
	var buf bytes.Buffer

	for i := 0; i < len(src); {
		ch, size := utf8.DecodeRuneInString(src[i:])
		switch {
		case unicode.IsDigit(ch):
			if prev == 0 {
				return "", ErrInvalidString
			}
			if amount := int(ch - '0'); amount != 0 {
				buf.WriteString(strings.Repeat(string(prev), amount))
			}
			prev = 0
		case ch == '\\':
			buf.WriteRune(prev)
			i++
			prev = rune(src[i])
		default:
			if prev != 0 {
				buf.WriteRune(prev)
			}
			prev = ch
		}
		i += size
	}
	if prev != 0 {
		buf.WriteRune(prev)
	}
	return buf.String(), nil
}

func isInvalidStr(src string) bool {
	if len(src) == 0 {
		return false
	}
	return unicode.IsDigit(rune(src[0])) || strings.HasSuffix(src, "\\")
}
