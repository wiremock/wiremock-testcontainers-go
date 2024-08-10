package testcontainers_wiremock_using_api_client

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
	container, err := RunDefaultContainerAndStopOnCleanup(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

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
