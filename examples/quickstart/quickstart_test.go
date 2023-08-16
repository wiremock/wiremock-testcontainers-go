package testcontainers_wiremock_quickstart

import (
	"context"
	"testing"

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
