package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var endpoint = String(os.Getenv("TRANSLATE_ENDPOINT")).
	orElse("http://localhost:8080/translate").
	get()

var original = false

func parseFlags() Ask {
	var (
		content string
		lang    string
		model   string
	)
	flag.StringVar(&content, "c", "", "The content to be translated")
	flag.StringVar(&lang, "l", "English", "The language to be translated")
	flag.StringVar(&model, "m", "gpt-3.5-turbo", "The chatGPT model to be chose")
	flag.BoolVar(&original, "o", false, "Set to get original OpenAI response")
	flag.Parse()

	payload := Ask{
		content,
		lang,
		model,
	}
	return payload
}

func translate(ask Ask) Answer {
	_bytes, _ := json.Marshal(ask)
	fmt.Println("original", original)
	if original {
		endpoint = "http://localhost:8080/translate/openai"
	}
	resp, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(_bytes))
	defer resp.Body.Close()

	var ans Answer
	json.NewDecoder(resp.Body).Decode(&ans)
	return ans
}

func main() {
	ask := parseFlags()
	fmt.Println("ask", ask)
	ans := translate(ask)
	fmt.Println(ans.Reply)
}
