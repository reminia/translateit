package main

import (
	"bytes"
	"compress/gzip"
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
	log.SetPrefix("[translate-debug] ")

	r := gin.Default()
	r.Use(CORS)
	r.POST("/translate", OpenAiProxy(handleOpenAiResponse))
	r.POST("/translate/openai", OpenAiProxy(identityResponseHandler))
	r.GET("/translate/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "I am up!")
	})

	port := String(os.Getenv("HTTP_PORT")).orElse("8080").get()
	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal("server start failed with error", err)
	}
}

func CORS(c *gin.Context) {
	origin := os.Getenv("ALLOW_ORIGINS")
	debug("ALLOW_ORIGINS", origin)
	originalHeader := c.GetHeader("Origin")
	if origin == "" {
		originalHeader = "*"
	} else if !String(origin).contains(",", originalHeader) {
		debug("Request origin", originalHeader)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Header("Access-Control-Allow-Origin", originalHeader)
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
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
		} else {
			debug("Response:", response.Status)
			dump, _ := httputil.DumpResponse(response, true)
			debug(string(dump))
			return nil
		}
	}
}

func identityResponseHandler(_ *gin.Context) ResponseCallBack {
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
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
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
			req.Header.Set("Authorization", "Bearer "+OpenaiKey)
			req.Header.Set("Content-Type", "application/json")
			var ask Ask
			err = c.ShouldBindJSON(&ask)
			if err != nil {
				panic(err)
			}
			ask.setDefault()
			body := OpenAiRequest{
				Model:       ask.Model,
				Temperature: OpenaiTemperature,
			}
			msg := Message{
				Role:    "user",
				Content: ask.Content,
			}
			body.Messages = []Message{systemMsg(ask.Lang), msg}
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
