# translateit ![ci](https://github.com/reminia/translateit/actions/workflows/go-build.yml/badge.svg)

A simple translate api that proxies to openai written by Golang.

## Run

1. Setup environment variables:

* OPENAI_KEY,  key of openai, required.
* OPENAI_MODEL, model of openai, optional, default to gpt-3.5-turbo.
* HTTP_PORT, port of the server, optional, default to 8080.

2. Build by `make build`.
3. Start server by `./translate`. Try it with:

```bash
curl -X POST \
-d '{"content": "Hi, how are you?", "lang": "Chinese", "model": "gpt-3.5-turbo"}' \
http://localhost:8080/translate
```

4. translate-cli is a way to test translate endpoints, start translate server before using the cli tool.
   Try it like `./translate-cli -c "content" -l "Chinese" -m "gpt-3.5-turbo"`, -l and -m are optional.
   The endpoint is http://localhost:8080 by default, customize it by setting TRANSLATE_ENDPOINT env var.

CORS is enabled by default, set up ALLOW_ORIGIN env var to control the origin otherwise all is allowed.
## Endpoints

There are 3 endpoints for now:

1. POST /translate, translate the content and return the simplified response with reply and reason only.
   Ask is the request and Answer is the response.
   ```golang
   // Ask the request of /translate
   type Ask struct {
   	Content string `json:"content"`
   	Lang    string `json:"lang,omitempty"`  //optional, default to English
   	Model   string `json:"model,omitempty"` //optional, default to gpt-3.5-turbo
   }

   // Answer the response of Ask
   type Answer struct {
   	Reply  string `json:"reply"`
   	Reason string `json:"reason"`
   }
   ```
2. POST /translate/openai, translate the content and return the original openai response for debugging purpose.
   The request is still Ask. The response is what described in the [openai doc](https://platform.openai.com/docs/api-reference/completions).
3. GET /translate/ping, a server healthy check api.

## Docker

Use the [image](https://github.com/reminia/translateit/pkgs/container/translateit) build by github action or build the image yourself.

* build: `docker build -t translate .`
* setup OPENAI_KEY, OPENAI_MODEL in an env file
* start: `docker run --env-file env -p 8081:8080 -d translate`

## Additional

Use the companion [translateit-btn](https://github.com/reminia/translateit-btn) with the api.
