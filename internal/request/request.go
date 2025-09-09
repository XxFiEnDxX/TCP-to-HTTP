package request

import (
	"bytes"
	"fmt"
	"io"
)

type parseState string

const (
	StateInit  parseState = "init"
	StateDone  parseState = "done"
	StateError parseState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parseState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("You just encounter malformed Request line!🙈")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("Unsupported HTTP version!🙈")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("Request in error state!")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, 0, nil
	}

	line := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(line, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	HttpParts := bytes.Split(parts[2], []byte("/"))

	if len(HttpParts) != 2 || string(HttpParts[0]) != "HTTP" || string(HttpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(HttpParts[1]),
	}, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, nil
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}
