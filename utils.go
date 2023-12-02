package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func (s String) contains(split string, sub string) bool {
	parts := strings.Split(s.get(), split)
	for _, part := range parts {
		if strings.Contains(part, sub) {
			return true
		}
	}
	return false
}

var debugEnabled = os.Getenv("DEBUG") == "true"

func debug(v ...any) {
	if debugEnabled {
		log.Println(v...)
	}
}
