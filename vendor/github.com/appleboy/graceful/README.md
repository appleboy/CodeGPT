# graceful

[![Run Tests](https://github.com/appleboy/graceful/actions/workflows/go.yml/badge.svg)](https://github.com/appleboy/graceful/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/appleboy/graceful/branch/master/graph/badge.svg?token=zPqtcz0Rum)](https://codecov.io/gh/appleboy/graceful)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/graceful)](https://goreportcard.com/report/github.com/appleboy/graceful)
[![Go Reference](https://pkg.go.dev/badge/github.com/gin-contrib/graceful.svg)](https://pkg.go.dev/github.com/gin-contrib/graceful)

Graceful shutdown package when a service is turned off by software function.

## Example

Add running job

```go
package main

import (
  "context"
  "log"
  "time"

  "github.com/appleboy/graceful"
)

func main() {
  m := graceful.NewManager()

  // Add job 01
  m.AddRunningJob(func(ctx context.Context) error {
    for {
      select {
      case <-ctx.Done():
        return nil
      default:
        log.Println("working job 01")
        time.Sleep(1 * time.Second)
      }
    }
  })

  // Add job 02
  m.AddRunningJob(func(ctx context.Context) error {
    for {
      select {
      case <-ctx.Done():
        return nil
      default:
        log.Println("working job 02")
        time.Sleep(500 * time.Millisecond)
      }
    }
  })

  <-m.Done()
}
```

You can also add shutdown jobs.

```go
package main

import (
  "context"
  "log"
  "time"

  "github.com/appleboy/graceful"
)

func main() {
  m := graceful.NewManager()

  // Add job 01
  m.AddRunningJob(func(ctx context.Context) error {
    for {
      select {
      case <-ctx.Done():
        return nil
      default:
        log.Println("working job 01")
        time.Sleep(1 * time.Second)
      }
    }
  })

  // Add job 02
  m.AddRunningJob(func(ctx context.Context) error {
    for {
      select {
      case <-ctx.Done():
        return nil
      default:
        log.Println("working job 02")
        time.Sleep(500 * time.Millisecond)
      }
    }
  })

  // Add shutdown 01
  m.AddShutdownJob(func() error {
    log.Println("shutdown job 01 and wait 1 second")
    time.Sleep(1 * time.Second)
    return nil
  })

  // Add shutdown 02
  m.AddShutdownJob(func() error {
    log.Println("shutdown job 02 and wait 2 second")
    time.Sleep(2 * time.Second)
    return nil
  })

  <-m.Done()
}
```

Using custom logger, see the [zerolog example](./_example/example03/logger.go)

```go
m := graceful.NewManager(
  graceful.WithLogger(logger{}),
)
```

get [more information](./_example/example03/main.go)
