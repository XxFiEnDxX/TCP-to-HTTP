package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_REQUEST_LINE = fmt.Errorf("You just encounter an ERROR!🙈")
var SEPARATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, b, nil
	}

	line := b[:idx]
	restOfMsg := b[idx+len(SEPARATOR):]

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_REQUEST_LINE
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error)
