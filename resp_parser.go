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
//I use unsafe pointer for string conversion because i need to lower memory allocation.
func (c *RespReader) parse() error {
	if len(c.payload) < 1 {
		return errors.New("empty data")
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

		if len(argsindex) >= 1 {
			c.command = *(*string)(unsafe.Pointer(&pp[argsindex[0]]))
		}

		if len(argsindex) > 1 {
			for i := 1; i < len(argsindex); i++ {
				c.args += *(*string)(unsafe.Pointer(&pp[argsindex[i]]))
				if i < len(argsindex)-1 {
					c.args += " "
				}
			}
		}
	case TypeString, TypeError, TypeBulkString:
		strpayload := pp[0][1:len(pp[0])]
		c.str = *(*string)(unsafe.Pointer(&strpayload))
	case TypeInteger:
		c.integer = 0 //oğuzhan: does not require integer value my case
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
func (c *RespReader) Size() int {
	return len(c.payload)
}

func (c *RespReader) String() string {
	return c.str
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
