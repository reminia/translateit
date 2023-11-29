package main

import (
	"fmt"
	"io"
)

func debug(desc string, reader io.Reader) {
	_bytes, _ := io.ReadAll(reader)
	fmt.Println("debug", desc, string(_bytes))
}

type String string

func (s String) orElse(that string) String {
	if s == "" {
		return String(that)
	}
	return s
}

func (s String) get() string {
	return string(s)
}
