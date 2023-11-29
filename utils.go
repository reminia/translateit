package main

import (
	"io"
	"net/http"
)

func dumpRespBody(resp *http.Response) string {
	bytes, _ := io.ReadAll(resp.Body)
	return string(bytes)
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
