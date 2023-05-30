[![Go Report Card](https://goreportcard.com/badge/github.com/Totus-Floreo/asperitas-on-go)](https://goreportcard.com/report/github.com/Totus-Floreo/asperitas-on-go)
[![Go](https://github.com/Totus-Floreo/asperitas-on-go/actions/workflows/go.yml/badge.svg)](https://github.com/Totus-Floreo/asperitas/blob/main/.github/workflows/go.yml)
[![Test](https://github.com/Totus-Floreo/asperitas-on-go/actions/workflows/test.yml/badge.svg)](https://github.com/Totus-Floreo/asperitas/blob/main/.github/workflows/test.yml)
[![codecov](https://codecov.io/gh/Totus-Floreo/asperitas-on-go/branch/main/graph/badge.svg?token=X9I4VJAFRC)](https://codecov.io/gh/Totus-Floreo/asperitas-on-go)
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

### Tests
```sh
#To Test with coverage
go test $(go list ./... | grep -v /test/) -coverprofile=overallcoverage ./...
# rn test in same folder with code, it will moved back to ./test/...

#To generate http
go test -coverpkg=./internal/<funcional>/<type>/ -coverprofile=<package> ./test/<funcional>/<package>_test.go

#To watch in http
go tool cover -html=test/<funcional>/<package>
```

#### Mockgen
```sh
#If module have interface for needed struct
mockgen --destination mocks/<interface>.go --package=mocks  --build_flags=--mod=mod <moduleURL> <interface>
#If module haven't interface for needed struct
#You need create interface that will implement struct ?? rewrite this later
mockgen -source <your_hand_made_interface>_interface.go -destination mocks/<your_hand_made_interface>.go -package=mocks 
```

## Thanks
Thank [d11z](https://github.com/d11z/asperitas) for the idea and frontend.

## License
This project is made available under the **MIT License**.

