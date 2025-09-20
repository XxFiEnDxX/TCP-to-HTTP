package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"tcp.to.http/internal/headers"
)

type parseState string

const (
	StateInit   parseState = "init"
	StateHeader parseState = "headers"
	StateBody   parseState = "body"
	StateDone   parseState = "done"
	StateError  parseState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
	Body          string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	state       parseState
}

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exist := headers.Get(name)
	if !exist {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
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

func (r *Request) hasBody() bool {
	length := getInt(r.Headers, "content-length", 0)
	return length > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentRead := data[read:]
		if len(currentRead) == 0 {
			break outer
		}
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE

		case StateInit:
			rl, n, err := parseRequestLine(currentRead)
			if err != nil {
				r.state = StateError
				return 0, nil
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = StateHeader

		case StateHeader:
			n, done, err := r.Headers.Parse(currentRead)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}
			read += n

			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}
		case StateBody:
			length := getInt(r.Headers, "content-length", 0)
			if length == 0 {
				panic("Chuncked not implemented")
			}
			remaining := min(length-len(r.Body), len(currentRead))
			r.Body += string(currentRead[:remaining])
			read += remaining

			if len(r.Body) == length {
				r.state = StateDone
			}
		case StateDone:
			break outer
		default:
			panic("something somethings")
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
