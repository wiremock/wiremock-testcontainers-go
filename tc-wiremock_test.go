package testcontainers_wiremock

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
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

func TestWireMockWithResource(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithMappingFile("hello", filepath.Join("testdata", "hello-world-resource.json")),
		WithFile("hello-world-resource-response.xml", filepath.Join("testdata", "hello-world-resource-response.xml")),
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

	statusCode, out, err := SendHttpGet(uri, "/hello-from-file")
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if !strings.Contains(out, "Hello, world!") {
		t.Fatalf("expected 'Hello, world!' but got %v", string(out))
	}
}

func SendHttpGet(url string, endpoint string) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, url+endpoint, nil)
	if err != nil {
		return -1, "", errors.Wrap(err, "unable to complete Get request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, "", errors.Wrap(err, "unable to complete Get request")
	}

	out, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, "", errors.Wrap(err, "unable to read response data")
	}

	return res.StatusCode, string(out), nil
}
