package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	r := gin.Default()
	r.POST("/translate", OpenAiProxy())
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}

var OPENAI_KEY = os.Getenv("OPENAI_KEY")
var OPENAI_MODEL = os.Getenv("OPENAI_MODEL")

type Data struct {
	Content string `json:content`
	Lang    string `json:lang,omitempty` //default to English
	Model   string `json:lang,omitempty` //default to gpt-3.5-turbo
}

func (d *Data) setDefault() {
	if d.Lang == "" {
		d.Lang = "English"
	}
	if d.Model == "" {
		d.Model = OPENAI_MODEL
	}

}

type OpenApiRequest struct {
	Model       string    `json:model`
	Messages    []Message `json:messages`
	Temperature uint      `json:temperature`
	MaxTokens   uint      `json:max_tokens,omitempty`
}

type Message struct {
	Role    string `json:role`
	Content string `json:content`
}

func systemMsg(lang string) Message {
	return Message{
		Role:    "system",
		Content: "You are a professional translator that can translate any language to " + lang,
	}
}

func OpenAiProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		openapi := "https://api.openai.com/v1/chat/completions"
		director := func(req *http.Request) {
			_url, err := url.Parse(openapi)
			if err != nil {
				panic(err)
			}
			req.URL = _url
			req.Host = ""
			req.Header.Set("Authorization", "Bearer "+OPENAI_KEY)
			req.Header.Set("Content-Type", "application/json")
			var data Data
			err = c.ShouldBindJSON(&data)
			if err != nil {
				panic(err)
			}
			data.setDefault()

			body := OpenApiRequest{
				Model:       OPENAI_MODEL,
				Temperature: 0,
			}
			content := fmt.Sprintf("please translate below passage to %s: %s",
				data.Lang, data.Content)
			msg := Message{
				Role:    "user",
				Content: content,
			}
			body.Messages = []Message{systemMsg(data.Lang), msg}
			var length int
			req.Body, length = NewRequestBody(body)
			req.ContentLength = int64(length)
		}
		proxy := &httputil.ReverseProxy{
			Director: director,
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

type Body io.ReadCloser

func NewRequestBody(body any) (Body, int) {
	_bytes, _ := json.Marshal(body)
	return io.NopCloser(bytes.NewBuffer(_bytes)), len(_bytes)
}
