package main

import (
	"fmt"
	"io"
)

func debug(desc string, reader io.Reader) {
	_bytes, _ := io.ReadAll(reader)
	fmt.Println("debug", desc, string(_bytes))
}
