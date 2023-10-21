A simple translate api that proxies to openai written by Golang.

## Run

1. Setup environment variables:

* OPENAI_KEY,  key of openai, required.
* OPENAI_MODEL, model of openai, optional, default to gpt-3.5-turbo.
* HTTP_PORT, port of the server, optional, default to 8080.

2. Start with `go run .`