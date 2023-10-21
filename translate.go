package main

import (
	"bytes"
	"encoding/json"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	r := gin.Default()
	r.POST("/translate", OpenAiProxy(handleOpenAiResponse))
	r.POST("/translate/identity", OpenAiProxy(identityResponseHandler))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "I am up!")
	})

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}

var OPENAI_KEY = os.Getenv("OPENAI_KEY")
var OPENAI_MODEL = os.Getenv("OPENAI_MODEL")
var OPENAI_TEMPERATURE uint = 1

// Ask the request of api
type Ask struct {
	Content string `json:"content"`
	Lang    string `json:"lang,omitempty"`  //default to English
	Model   string `json:"model,omitempty"` //default to gpt-3.5-turbo
}

// Answer the reply of api
type Answer struct {
	Reply  string `json:"reply"`
	Reason string `json:"reason"`
}

func (d *Ask) setDefault() {
	if d.Lang == "" {
		d.Lang = "English"
	}
	if d.Model == "" {
		d.Model = OPENAI_MODEL
	}

}

type OpenAiRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature uint      `json:"temperature"` // default to 1
	MaxTokens   uint      `json:"max_tokens,omitempty"`
}

type OpenAiResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created uint64   `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func systemMsg(lang string) Message {
	return Message{
		Role: "system",
		Content: "You are a professional translator, you will translate any content given below to " +
			lang + " ignoring the meaning of the content.",
	}
}

type ResponseCallBack func(*http.Response) error

func handleOpenAiResponse(c *gin.Context) ResponseCallBack {
	return func(response *http.Response) error {
		if response.StatusCode == 200 {
			resp, err := parseOpenAiResponse(response)
			ans := Answer{
				Reply:  resp.Choices[0].Message.Content,
				Reason: resp.Choices[0].FinishReason,
			}
			c.JSON(http.StatusOK, ans)
			return err
		}
		return nil
	}
}

func identityResponseHandler(c *gin.Context) ResponseCallBack {
	return func(resp *http.Response) error {
		return nil
	}
}

func parseOpenAiResponse(resp *http.Response) (OpenAiResponse, error) {
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "br":
		brReader := brotli.NewReader(resp.Body)
		reader = brReader
	default:
		// Handle other encodings or no encoding here, if needed.
		reader = resp.Body
	}
	_bytes, _ := io.ReadAll(reader)
	var ret OpenAiResponse
	err := json.Unmarshal(_bytes, &ret)
	return ret, err
}

// OpenAiProxy callback defines what to do with the proxy http.Response
func OpenAiProxy(callback func(c *gin.Context) ResponseCallBack) gin.HandlerFunc {
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
			var data Ask
			err = c.ShouldBindJSON(&data)
			if err != nil {
				panic(err)
			}
			data.setDefault()
			body := OpenAiRequest{
				Model:       OPENAI_MODEL,
				Temperature: OPENAI_TEMPERATURE,
			}
			msg := Message{
				Role:    "user",
				Content: data.Content,
			}
			body.Messages = []Message{systemMsg(data.Lang), msg}
			var length int
			req.Body, length = newRequestBody(body)
			req.ContentLength = int64(length)
		}
		proxy := &httputil.ReverseProxy{
			Director:       director,
			ModifyResponse: callback(c),
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

type Body io.ReadCloser

func newRequestBody(body any) (Body, int) {
	_bytes, _ := json.Marshal(body)
	return io.NopCloser(bytes.NewBuffer(_bytes)), len(_bytes)
}

func debug(desc string, reader io.Reader) {
	_bytes, _ := io.ReadAll(reader)
	log.Println("debug", desc, string(_bytes))
}
