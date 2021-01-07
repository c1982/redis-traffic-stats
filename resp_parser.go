package main

import (
	"bytes"
	"errors"
	"regexp"
	"unsafe"
)

const (
	TypeString = iota
	TypeBulkString
	TypeInteger
	TypeArray
	TypeError
	TypeUnknown
)

//RespReader presents RESP packges from ethernet.
type RespReader struct {
	payload []byte
	command string
	args    string
	size    int
}

//NewRespReader create a RESP object on the app
func NewRespReader(payload, sep []byte, cleaner *regexp.Regexp, size int) (*RespReader, error) {
	r := &RespReader{
		payload: payload,
	}
	err := r.parse(sep, cleaner, size)
	return r, err
}

//parse basic RESP parser for my case.
//I use unsafe pointer for string conversion because I need to lower memory allocation.
func (c *RespReader) parse(sep []byte, cleaner *regexp.Regexp, maxsize int) error {
	if c.payload == nil {
		return errors.New("payload is nil")
	}

	if len(c.payload) < 1 {
		return errors.New("payload is empty")
	}

	switch c.Type() {
	case TypeArray:
		pp := bytes.Split(c.payload, []byte{'\r', '\n'})
		argsindex := []int{}
		for i := 0; i < len(pp); i++ {
			if bytes.HasPrefix(pp[i], []byte{'$'}) {
				argsindex = append(argsindex, i+1)
				if len(argsindex) > 2 {
					break
				}
			}
		}

		if len(argsindex) > 0 {
			if argsindex[0] < len(pp) {
				c.command = *(*string)(unsafe.Pointer(&pp[argsindex[0]]))
			}
		}

		if len(argsindex) >= 1 {
			first := pp[argsindex[1]]
			if len(first) > maxsize {
				if maxsize > 0 {
					first = pp[argsindex[1]][0 : maxsize-1]
				}
			}
			if len(sep) > 0 {
				first = c.removeLast(first, sep)
			}
			if cleaner != nil {
				first = c.cleanMatched(first, cleaner, sep)
			}
			c.args = *(*string)(unsafe.Pointer(&first))
		}
	default:
		return errors.New("unsuported type")
	}

	return nil
}

func (c *RespReader) removeLast(payload []byte, sep []byte) []byte {
	explode := bytes.Split(payload, sep)
	if len(explode) < 2 {
		return payload
	}
	explode = explode[0 : len(explode)-1]
	return bytes.Join(explode, sep)
}

func (c *RespReader) cleanMatched(payload []byte, pattern *regexp.Regexp, trim []byte) []byte {
	//TODO: probably ReplaceAll is slow. refactor it.
	return bytes.TrimSuffix(pattern.ReplaceAll(payload, []byte{}), trim)
}

//Command RESP command name
func (c *RespReader) Command() string {
	return c.command
}

//Args RESP command arguments
func (c *RespReader) Args() string {
	return c.args
}

//Size RESP size of the command
func (c *RespReader) Size() float64 {
	return float64(len(c.payload))
}

//Type returns the RESP command or response type
func (c *RespReader) Type() int {
	switch c.payload[0] {
	case '+':
		return TypeString
	case '$':
		return TypeBulkString
	case '*':
		return TypeArray
	case ':':
		return TypeInteger
	case '-':
		return TypeError
	}

	return TypeUnknown
}
