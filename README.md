# DevSecTools SDK for Go

**Experimental. Untested. Not ready for primetime. Exploring code generation.**

## Overview

A Go SDK for interacting with the [DevSecTools API].

This client provides an easy way to interact with the [DevSecTools API], which scans websites for security-related information such as HTTP version support and TLS configurations.

* ✅ Requires [Go] 1.23+.
* ✅ Uses [Goroutines] to handle HTTP requests and supports both synchronous and asynchronous (parallel) requests.
* ✅ Supports [versions of Go which receive support](https://go.dev/doc/devel/release#policy) from the core team.

## Model

* [openapi.json](https://github.com/northwood-labs/devsec-tools/raw/refs/heads/main/openapi.json)

## Usage

### Instantiating with default configuration

```go
package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "github.com/northwood-labs/devsec-tools-sdk-go/devsectools"
)

func main() {
  client := devsectools.NewClient()

  // ...
}
```

### Custom configuration

```go
package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "github.com/northwood-labs/devsec-tools-sdk-go/devsectools"
)

func main() {
  client := devsectools.NewClientWithConfig(&devsectools.Config{
    Endpoint: &devsectools.LOCALDEV, // Use the local development environment
    Timeout:  10 * time.Second,      // Set timeout to 10 seconds
  })

  // ...
}
```

### Updating configuration at runtime

```go
package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "github.com/northwood-labs/devsec-tools-sdk-go/devsectools"
)

func main() {
  client := devsectools.NewClient()

  // One or the other...
  client.SetEndpoint(&devsectools.LOCALDEV)
  client.SetBaseURL("http://localhost:8080")

  client.SetTimeout(15 * time.Second)

  // ...
}
```

### Making single requests

```go
package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "github.com/northwood-labs/devsec-tools-sdk-go/devsectools"
)

var ctx := context.Background()

func main() {
  client := devsectools.NewClient()

  httpInfo, err := client.HTTP(ctx, "example.com")
  if err != nil {
    log.Fatalf("Error fetching HTTP info: %v", err)
  }

  fmt.Printf(
    "HTTP Support: HTTP/1.1=%v, HTTP/2=%v, HTTP/3=%v\n",
    httpInfo.HTTP11,
    httpInfo.HTTP2,
    httpInfo.HTTP3,
  )
}
```

### Making parallel/batch requests

```go
package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "github.com/northwood-labs/devsec-tools-sdk-go/devsectools"
)

var ctx := context.Background()

func main() {
  client := devsectools.NewClient()

  // Define batch requests
  batchRequests := []devsectools.BatchRequest{
    {Method: "http", URL: "apple.com",  Result: &devsectools.HttpResponse{}},
    {Method: "tls",  URL: "google.com", Result: &devsectools.TlsResponse{}},
  }

  // Execute batch requests
  client.Batch(ctx, batchRequests)

  // Print results
  for _, req := range batchRequests {
    if req.Err != nil {
      log.Printf("Error fetching %s: %v\n", req.Method, req.Err)
      continue
    }

    fmt.Printf("Result for %s: %+v\n", req.Method, req.Result)
  }
}
```

[DevSecTools API]: https://devsec.tools
[Go]: https://go.dev
[Goroutines]: https://go.dev/tour/concurrency
