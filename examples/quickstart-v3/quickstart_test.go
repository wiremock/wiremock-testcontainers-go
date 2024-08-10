package testcontainers_wiremock_quickstart

import (
	"context"
	"strings"
	"testing"

	. "github.com/wiremock/wiremock-testcontainers-go"
)

func TestWireMock(t *testing.T) {
	ctx := context.Background()

	container, err := RunContainerAndStopOnCleanup(ctx, t,
		WithImage("docker.io/wiremock/wiremock:3.9.1"),
		WithMappingFile("path", "path-template.json"),
	)
	if err != nil {
		t.Fatal(err)
	}

	statusCode, out, err := SendHttpGet(container, "/v1/contacts/12345/addresses/99876", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}

	// Verify the response
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}

	if !strings.Contains(out, "12345") {
		t.Fatalf("expected '12345' but got %s", out)
	}
	if !strings.Contains(out, "99876") {
		t.Fatalf("expected '99876' but got %s", out)
	}
}
