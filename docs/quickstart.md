# Quick Start - WireMock on Testcontainers Go

If you have the Golang development environment ready to go and do not want a step by step guide,
you can just clone the project’s repository in
[`wiremock/wiremock-testcontainers-go`](https://github.com/wiremock/wiremock-testcontainers-go),
go to the `examples/quickstart` directory
and run `go build` and then `go test` in the root or try the [examples](https://github.com/wiremock/wiremock-testcontainers-go/blob/main/examples/quickstart/quickstart_test.go).
Any pull requests will be welcome!

## Pre-requisites

- Go 1.20 or above
- Docker-API compatible container runtime ([more info](https://golang.testcontainers.org/system_requirements/docker/))

## Create test project

Create the `go.mod` file with the following content:

```go
module wiremock.org/testcontainers-go-quickstart

go 1.20

require (
    github.com/pkg/errors v0.9.1
    github.com/wiremock/wiremock-testcontainers-go v1.0.0-alpha-7
)
```

Then, run `go mod install` to install the dependencies and prepare the environment

## Create the test file

Create a `quickstart_test.go` file with the package name.
Add dependencies we will need for this demo, and also create the test stub:

```go
package testcontainers_wiremock_quickstart

import (
  "context"
  . "github.com/wiremock/wiremock-testcontainers-go"
  "testing"
)

func TestWireMock(t *testing.T) {
    // Our future work will be here
}
```

## Create the test resource

For our demo, we will need to expose a test WireMock Mapping.
Create the `hello-world.json` file with the following content:

```json
{
  "request": {
    "method": "GET",
    "url": "/hello"
  },

  "response": {
    "status": 200,
    "body": "Hello, world!"
  }
}
```

## Add Testcontainers initialization

In `func TestWireMock(t *testing.T)`, add the following code:

### wiremock v2
```go
// Create Container
ctx := context.Background()
container, err := RunContainerAndStopOnCleanup(ctx, t,
    WithMappingFile("hello", "hello-world.json"),
)
if err != nil {
    t.Fatal(err)
}
```

### wiremock v3
```go
// Create Container
ctx := context.Background()
container, err := RunContainerAndStopOnCleanup(ctx, t,
    WithImage("docker.io/wiremock/wiremock:3.9.1"),
    WithMappingFile("hello", "hello-world.json"),
)
if err != nil {
    t.Fatal(err)
}
```

## Add logic to send a request

Now, we will need to send an HTTP request to our test API.
To do so, we will use a build-in method:

```go
func TestWireMock(t *testing.T) {
    // ... Previous initialization code

    // Send a simple HTTP GET request to the mocked API
    statusCode, out, err := SendHttpGet(container, "/hello", nil)
    if err != nil {
        t.Fatal(err, "Failed to get a response")
    }

    // ... Validation will be here
}
```

In the code above, we used the `SendHttpGet` method to send a HTTP GET request. The library also offers methods to send requests with other HTTP methods, i.e. `SendHttpPost`, `SendHttpDelete`,`SendHttpPatch`, `SendHttpPut`.

## Verify the response

Now, add the verification logic that will check correctness of the WireMock response:

```go
func TestWireMock(t *testing.T) {
    // ... Previous initialization code

    // ... Previous HTTP request send code

    // Verify the response
    if statusCode != 200 {
        t.Fatalf("expected HTTP-200 but got %d", statusCode)
    }

    if string(out) != "Hello, world!" {
        t.Fatalf("expected 'Hello, world!' but got %v", string(out))
    }
}
```

## Run the test

We are finally ready to run the test!
Do the following:

```bash
go test
```

If everything goes right, you will see the following console output:

![Quick Start Demo](./images/quickstart.gif)

## Read more

See the [documentation root](../README.md) for the references to more features, examples and demos.
