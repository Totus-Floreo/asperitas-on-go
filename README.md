[![Go Report Card](https://goreportcard.com/badge/github.com/Totus-Floreo/asperitas-on-go)](https://goreportcard.com/report/github.com/Totus-Floreo/asperitas-on-go)
[![Go](https://github.com/Totus-Floreo/asperitas-on-go/actions/workflows/go.yml/badge.svg)](https://github.com/Totus-Floreo/asperitas/blob/main/.github/workflows/go.yml)
[![Status](https://badgen.net/badge/status/indevelopment/blue?icon=github)](https://github.com/Totus-Floreo/asperitas-on-go)
[![MIT](https://badgen.net/badge/license/MIT/blue)](https://github.com/Totus-Floreo/asperitas-on-go/blob/main/LICENSE)

# asperitas-go
Simple reddit-clone based on [Asperitas](https://github.com/d11z/asperitas) js-front and my golang-backend

## Quickstart

### Setup
```go
package <your_package>

import "os"

func <your_func>() {
    Os.Setenv("signature", "<Your_signature>")
    Os.Setenv("port", "<Your_port>")
    Os.Setenv("redis", "<Your_redis_port>")
	Os.Serenv("pg_url", "<<username>:<password>@<host>:<port>/<database>>")
}
```
if you use a non local redis db, needs setup this block in ./cmd/asperitas/main.go
```go
rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost" + os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
    
```
### Run the app
```sh
go run ./cmd/asperitas/main.go
```
### Alternative method

Create bash file and run

```sh
#Linux
#!/bin/bash

signature=<signature>
port=:<port>
redis=:<redis-port>
pg_url=<<username>:<password>@<host>:<port>/<database>>

export signature port redis pg_url

go run ./cmd/asperitas/main.go
```

## Thanks
Thank [d11z](https://github.com/d11z/asperitas) for the idea and frontend.

## License
This project is made available under the **MIT License**.

