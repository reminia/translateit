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
	r := gin.Default()
	r.Use(CORS)
	r.POST("/translate", OpenAiProxy(handleOpenAiResponse))
	r.POST("/translate/openai", OpenAiProxy(identityResponseHandler))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "I am up!")
	})

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal("server start failed with error", err)
	}
}

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
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
		}
		return nil
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
