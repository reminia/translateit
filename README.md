A simple translate api that proxies to openai written by Golang.

## Run

1. Setup environment variables:

* OPENAI_KEY,  key of openai, required.
* OPENAI_MODEL, model of openai, optional, default to gpt-3.5-turbo.
* HTTP_PORT, port of the server, optional, default to 8080.

2. Start with `go run .`

## Endpoints

There are 3 endpoints for now:

1. POST /translate, translate the content and return the simplified response with reply and reason only.
   Ask is the request and Answer is the response.
   ```golang
   type Ask struct {
   	Content string `json:"content"`
   	Lang    string `json:"lang,omitempty"`  //optional, default to English
   	Model   string `json:"model,omitempty"` //optional, default to gpt-3.5-turbo
   }

   // Answer the reply of api
   type Answer struct {
   	Reply  string `json:"reply"`
   	Reason string `json:"reason"`
   }
   ```
2. POST /translate/openai, translate the content and return the original openai response for debugging purpose.
   The request is still Ask. The response is what described in the [openai doc](https://platform.openai.com/docs/api-reference/completions).
3. GET /health, a server healthy check api.
