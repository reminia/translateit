A simple translate api that proxies to openapi written by Golang.

## Run

1. Setup environment variables:

* OPENAPI_KEY,  key of openapi, required.
* OPENAPI_MODEL, model of openapi, optional, default to gpt-3.5-turbo.
* HTTP_PORT, port of the server, optional, default to 8080.

2. Start with `go run .`