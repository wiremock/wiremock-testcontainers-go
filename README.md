# WireMock Module for Testcontainers Go

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/wiremock/wiremock-testcontainers-go)](https://github.com/wiremock/wiremock-testcontainers-go/releases)
[![Slack](https://img.shields.io/badge/slack-slack.wiremock.org-brightgreen?style=flat&logo=slack)](https://slack.wiremock.org/)
[![GitHub contributors](https://img.shields.io/github/contributors/wiremock/wiremock-testcontainers-go)](https://github.com/wiremock/wiremock-testcontainers-go/graphs/contributors)

## Note

The Testcontainers module does not work with the official image at the moment,
because the Mappings and Files directories are not initialized there.
Use a custom image on the top of it, see `Dockerfile`.

## Supported features

The following features are now explicitly included in the module's API:

- Passing API Mapping files
- Passing Resource files

More features will be added over time.

## Usage

```golang

import (
 "context"
 "net/http"
 "testing"

 "github.com/pkg/errors"
 "github.com/wiremock/wiremock-testcontainers-go"
)

func TestWireMock(t *testing.T) {
 // Create Container
 ctx := context.Background()
 container, err := RunContainer(ctx,
  WithMappingFile("hello", filepath.Join("testdata", "hello-world.json")),
 )
 if err != nil {
  t.Fatal(err)
 }

 // Clean up the container after the test is complete
 t.Cleanup(func() {
  if err := container.Terminate(ctx); err != nil {
   t.Fatalf("failed to terminate container: %s", err)
  }
 })

 uri, err := GetURI(ctx, container)
 if err != nil {
  t.Fatal(err, "unable to get container endpoint")
 }

 statusCode, out, err := SendHttpGet(uri, "/hello")
 if err != nil {
  t.Fatal(err, "Failed to get a response")
 }
 if statusCode != 200 {
  t.Fatalf("expected HTTP-200 but got %d", statusCode)
 }
 if string(out) != "Hello, world!" {
  t.Fatalf("expected 'Hello, world!' but got %v", string(out))
 }
}
```

## License

The module is licensed under [Apache License v.2](./LICENSE)
