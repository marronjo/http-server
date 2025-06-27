package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"
const colon = ":"
const space = " "
const empty = ""

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlf_idx := bytes.Index(data, []byte(crlf))
	if crlf_idx == -1 {
		return 0, false, nil
	}
	if crlf_idx == 0 {
		return 2, true, nil
	}

	colon_idx := bytes.Index(data, []byte(colon))
	if colon_idx == -1 {
		return 0, false, fmt.Errorf("invalid header format : missing colon")
	}

	d := string(data)
	key := d[:colon_idx]
	val := d[colon_idx+1:]

	trimKey := strings.TrimLeftFunc(key, unicode.IsSpace)
	if strings.Contains(trimKey, space) {
		return 0, false, fmt.Errorf("invalid header format : whitespace between key and colon")
	}

	trimVal := strings.TrimSpace(val)
	if trimVal == empty {
		return 0, false, fmt.Errorf("invalid header format : empty value")
	}

	h.Set(trimKey, trimVal)

	return crlf_idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
