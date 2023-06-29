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

	endpoint, err := container.PortEndpoint(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get container endpoint")
	}

	c := Client(endpoint)
	res, err := http.Get(c.url + "/upper?word=" + word)
	if err != nil {
		return "", errors.Wrap(err, "unable to complete Get request")
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read response data")
	}

	if string(out) != "Hello, world!" {
		t.Errorf("expected 'Hello, world!' but got %v", string(out))
	}

}
