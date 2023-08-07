# WireMock Module for Testcontainers Go

[![GoDoc](https://godoc.org/github.com/wiremock/wiremock-testcontainers-go?status.svg)](http://godoc.org/github.com/wiremock/wiremock-testcontainers-go)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/wiremock/wiremock-testcontainers-go)](https://github.com/wiremock/wiremock-testcontainers-go/releases)
[![Slack](https://img.shields.io/badge/slack-slack.wiremock.org-brightgreen?style=flat&logo=slack)](https://slack.wiremock.org/)
[![GitHub contributors](https://img.shields.io/github/contributors/wiremock/wiremock-testcontainers-go)](https://github.com/wiremock/wiremock-testcontainers-go/graphs/contributors)

This module allows provisioning the [WireMock API mock server](https://wiremock.org/) as a standalone container within your unit tests,
based on the official [WireMock Docker](https://github.com/wiremock/wiremock-docker) images (`2.35.0-1` or above) or compatible custom images.

You can learn more about WireMock and Golang on this [WireMock solutions page](https://wiremock.org/docs/solutions/golang/).

## Supported features

The following features are now explicitly included in the module's API:

- Passing API Mapping files
- Passing Resource files
- Sending HTTP requests to the mocked container

More features will be added over time.

## Quick Start

See the [Quick Start Guide](./docs/quickstart.md).
Just a teaser of how it feels at the real speed!

![Quickstart demo GIF](./docs/images/quickstart.gif)

## Requirements

- Golang version 1.17 or above, so all modern Golang projects should be compatible with it.
- The module supports the official [WireMock Docker](https://github.com/wiremock/wiremock-docker) images 2.35.0-1 or above.
- Custom images are supported too as long as they follow the same CLI and API structure.

## Usage

```golang
import (
  "context"
  . "github.com/wiremock/wiremock-testcontainers-go"
  "testing"
)

func TestWireMock(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithMappingFile("hello", "hello-world.json"),
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

	// Send the HTTP GET request to the mocked API
	statusCode, out, err := SendHttpGet(container, "/hello", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}

	// Verify the response
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

## References

- [WireMock Website](https://wiremock.org/)
- [WireMock and GOlang Solutions page](https://wiremock.org/docs/solutions/golang/)
- [Testcontainers for Go](https://golang.testcontainers.org/)
- [WireMock Module page on the Testcontainers marketplace](https://testcontainers.com/modules/wiremock/)
