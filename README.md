[![Build Status](https://travis-ci.org/ekomobile/grpc-consul-resolver.svg)](https://travis-ci.org/ekomobile/grpc-consul-resolver)
[![GitHub release](https://img.shields.io/github/release/nvx/grpc-consul-resolver.svg)](https://github.com/ekomobile/nvx/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/nvx/grpc-consul-resolver)](https://goreportcard.com/report/github.com/nvx/grpc-consul-resolver)
[![GoDoc](https://godoc.org/github.com/nvx/grpc-consul-resolver?status.svg)](https://godoc.org/github.com/nvx/grpc-consul-resolver)

# gRPC Consul resolver
This lib resolves Consul services by name. 

# Usage

Somewhere in your `init` code: 
```go
import (
    "github.com/nvx/grpc-consul-resolver"
)

// Will query consul every 5 seconds.
resolver.RegisterDefault(time.Second * 5)
```

Getting connection:
```go
conn, err := grpc.DialContext(ctx, "consul://service/my-awesome-service")
```

Using a tag
```go
conn, err := grpc.DialContext(ctx, "consul://service/my-awesome-service/http")
```

Using a prepared query
```go
conn, err := grpc.DialContext(ctx, "consul://query/my-awesome-query")
```

With round-robin balancer:
```go
import (
    "google.golang.org/grpc/balancer/roundrobin"
)

conn, err := grpc.DialContext(ctx, "consul://service/my-awesome-service", grpc.WithBalancerName(roundrobin.Name))
```
