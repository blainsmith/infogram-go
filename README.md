# Infogram API Go

[![Go Reference](https://pkg.go.dev/badge/github.com/blainsmith/infogram-go.svg)](https://pkg.go.dev/github.com/blainsmith/infogram-go)

Go client for the [Infogram API](https://developers.infogr.am/rest/).

## Installation

```
$ go get github.com/blainsmith/infogram-go
```

## Usage

**Default Client**
```go
import "github.com/blainsmith/infogram-go"

func main() {
    client := infogram.NewClient("api-key", "api-secret")

    infographic, _ := client.Infographic(13)
}
```

**With Options**
```go
import "github.com/blainsmith/infogram-go"

func main() {
    httpClient := http.Client{
        Timeout: 10 * time.Second,
    }

    client := infogram.NewClient(
        "api-key",
        "api-secret",
        infogram.ClientOptHTTPClient(&httpClient),
        infogram.ClientOptEndpoint("https://example.com/infogram"),
    )

    infographic, _ := client.Infographic(13)
}
```