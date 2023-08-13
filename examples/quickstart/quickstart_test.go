package testcontainers_wiremock_quickstart

import (
	"context"
	"net/http"
	"testing"

	"github.com/wiremock/go-wiremock"

	. "github.com/wiremock/wiremock-testcontainers-go"
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

	// Send a simple HTTP GET request to the mocked API
	statusCode, out, err := SendHttpGet(container, "/hello", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}

	// Verify the response
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}

	if string(out) != "Hello, world!" {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestWireMockClient(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// Use the WireMock client to stub a new endpoint manually
	err = container.Client.StubFor(
		wiremock.Get(wiremock.URLEqualTo("/hello")).
			WillReturnResponse(
				wiremock.NewResponse().
					WithJSONBody(map[string]string{"result": "Hello, world!"}).
					WithHeader("Content-Type", "application/json").
					WithStatus(http.StatusOK),
			),
	)

	if err != nil {
		t.Fatal(err)
	}

	statusCode, out, err := SendHttpGet(container, "/hello", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}

	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}

	if string(out) != `{"result":"Hello, world!"}` {
		t.Fatalf("expected 'Hello, world!' but got %v", out)
	}
}
