package testcontainers_wiremock

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func TestWireMock(t *testing.T) {

	ctx := context.Background()

	// Create Container
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

	uri, err := GetURI(ctx, container)
	if err != nil {
		t.Error(err, "unable to get container endpoint")
	}

	out, err := SendHttpGet(uri, "/hello")
	if err != nil {
		t.Error(err, "Failed to get a response")
	}

	if string(out) != "Hello, world!" {
		t.Errorf("expected 'Hello, world!' but got %v", string(out))
	}
}

func SendHttpGet(url string, endpoint string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url+endpoint, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to complete Get request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to complete Get request")
	}

	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read response data")
	}

	return string(out), nil
}
