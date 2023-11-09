package main

import "os"

// Ask the request of /translate
type Ask struct {
	Content string `json:"content"`
	Lang    string `json:"lang,omitempty"`  //optional, default to English
	Model   string `json:"model,omitempty"` //optional, default to gpt-3.5-turbo
}

// Answer the reply of Ask
type Answer struct {
	Reply  string `json:"reply"`
	Reason string `json:"reason"`
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
			lang + " ignoring the meaning of the content. And you should keep the original format without change and" +
			"also please don't translate any code blocks inside.",
	}
}

var OPENAI_KEY = os.Getenv("OPENAI_KEY")
var OPENAI_MODEL = os.Getenv("OPENAI_MODEL")
var OPENAI_TEMPERATURE uint = 1

func (d *Ask) setDefault() {
	if d.Lang == "" {
		d.Lang = "English"
	}
	if d.Model == "" {
		d.Model = OPENAI_MODEL
		if d.Model == "" {
			d.Model = "gpt-3.5-turbo"
		}
	}
}
