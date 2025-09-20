package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	for _, char := range str {
		switch {
		case char >= 'a' && char <= 'z':
		case char >= 'A' && char <= 'Z':
		case char >= '0' && char <= '9':
		case char == '!' ||
			char == '#' ||
			char == '$' ||
			char == '%' ||
			char == '&' ||
			char == '\'' ||
			char == '*' ||
			char == '+' ||
			char == '-' ||
			char == '.' ||
			char == '^' ||
			char == '_' ||
			char == '`' ||
			char == '|' ||
			char == '~':
		default:
			return false
		}
	}
	return true
}

var rn = []byte("\r\n")

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header line!🤨")
	}

	fieldName := parts[0]
	fieldValue := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(fieldName, []byte(" ")) {
		return "", "", fmt.Errorf("malformed header field name!🤨")
	}

	return string(fieldName), string(fieldValue), nil
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) (string, bool) {
	str, ok := h.headers[strings.ToLower(name)]
	return str, ok
}

func (h *Headers) Replace(name, value string) {
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		// EMPTY HEADER
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		fieldName, fieldValue, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(fieldName)) {
			return 0, false, fmt.Errorf("malformed header name")
		}
		read += (idx + len(rn))
		h.Set(fieldName, fieldValue)
	}

	return read, done, nil
}
