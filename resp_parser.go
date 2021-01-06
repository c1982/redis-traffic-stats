package main

import (
	"bytes"
	"errors"
	"unsafe"
)

const (
	TypeString = iota
	TypeBulkString
	TypeInteger
	TypeArray
	TypeError
	TypeUnknown

	MaxCommandArgsSize = 50
)

//RespReader presents RESP packges from ethernet.
type RespReader struct {
	payload []byte
	command string
	args    string
	size    int
	str     string
	integer int64
}

//NewRespReader create a RESP object on the app
func NewRespReader(payload []byte) (*RespReader, error) {
	r := &RespReader{
		payload: payload,
	}
	err := r.parse()
	return r, err
}

//parse basic RESP parser for my case.
//I use unsafe pointer for string conversion because I need to lower memory allocation.
func (c *RespReader) parse() error {
	if c.payload == nil {
		return errors.New("payload is nil")
	}

	if len(c.payload) < 1 {
		return errors.New("payload is empty")
	}

	pp := bytes.Split(c.payload, []byte{'\r', '\n'})
	argsindex := []int{}

	switch c.Type() {
	case TypeArray:
		for i := 0; i < len(pp); i++ {
			if bytes.HasPrefix(pp[i], []byte{'$'}) {
				argsindex = append(argsindex, i+1)
			}
		}

		if len(argsindex) > 0 {
			if argsindex[0] < len(pp) {
				c.command = *(*string)(unsafe.Pointer(&pp[argsindex[0]]))
			}
		}

		if len(argsindex) >= 1 {
			if len(pp[argsindex[1]]) > MaxCommandArgsSize {
				capturefirst50 := pp[argsindex[1]][0 : MaxCommandArgsSize-1]
				c.args = *(*string)(unsafe.Pointer(&capturefirst50))
			} else {
				c.args = *(*string)(unsafe.Pointer(&pp[argsindex[1]]))
			}
		}
	case TypeString, TypeError, TypeBulkString, TypeInteger: //does not require
	}

	return nil
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
