package testcontainers_wiremock

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/wiremock/go-wiremock"
)

const defaultV3WireMockImage = "docker.io/wiremock/wiremock:3.9.1"

func TestWireMockV3(t *testing.T) {
	// Create Container, by deliberately using the RunContainer API
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
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

	statusCode, out, err := SendHttpGet(container, "/hello", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if string(out) != "Hello, world!" {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3SendHttpGetWorksWithQueryParamsPassedInArgument(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainerAndStopOnCleanup(ctx, t,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("get", filepath.Join("testdata", "url-with-query-params.json")),
	)
	if err != nil {
		t.Fatal(err)
	}

	statusCode, out, err := SendHttpGet(container, "/get", map[string]string{"query": "test"})
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if out != "" {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3SendHttpGetWorksWithQueryParamsProvidedInURL(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainerAndStopOnCleanup(ctx, t,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("get", filepath.Join("testdata", "url-with-query-params.json")),
	)
	if err != nil {
		t.Fatal(err)
	}

	statusCode, out, err := SendHttpGet(container, "/get?query=test", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if out != "" {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3SendHttpDelete(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainerAndStopOnCleanup(ctx, t,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("delete", filepath.Join("testdata", "204-no-content.json")),
	)
	if err != nil {
		t.Fatal(err)
	}

	statusCode, out, err := SendHttpDelete(container, "/delete")
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 204 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if out != "" {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3SendHttpPatch(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainerAndStopOnCleanup(ctx, t,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("patch", filepath.Join("testdata", "200-patch.json")),
		WithFile("sample-model.json", filepath.Join("testdata", "sample-model.json")),
	)
	if err != nil {
		t.Fatal(err)
	}
	var jsonBody = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	statusCode, out, err := SendHttpPatch(container, "/patch", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if !strings.Contains(out, "sampleField1") {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3SendHttpPut(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainerAndStopOnCleanup(ctx, t,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("put", filepath.Join("testdata", "200-put.json")),
		WithFile("sample-model.json", filepath.Join("testdata", "sample-model.json")),
	)
	if err != nil {
		t.Fatal(err)
	}

	var jsonBody = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	statusCode, out, err := SendHttpPut(container, "/put", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if !strings.Contains(out, "sampleField1") {
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3WireMockClient(t *testing.T) {
	// Create Container
	ctx := context.Background()
	container, err := RunContainerAndStopOnCleanup(ctx, t, WithImage(defaultV3WireMockImage))
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
		t.Fatalf("expected 'Hello, world!' but got %s", out)
	}
}

func TestV3Health(t *testing.T) {
	// Create Container, by deliberately using the RunContainer API
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
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

	statusCode, out, err := SendHttpGet(container, "/__admin/health", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	if !strings.Contains(out, "Wiremock is ok") {
		t.Fatalf("expected 'Wiremock is ok' in response body but got %s", out)
	}
}

func TestV3FormParameters(t *testing.T) {
	// Create Container, by deliberately using the RunContainer API
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("v3", filepath.Join("testdata", "200-v3-form-parameters.json")),
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

	endpoint, err := GetURI(ctx, container)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("tool", "WireMock")
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint+"/install-tool", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	statusCode, _, err := sendTestRequest(t, req)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
}

func TestV3JsonSchema(t *testing.T) {
	// Create Container, by deliberately using the RunContainer API
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("v3", filepath.Join("testdata", "200-v3-json-schema.json")),
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

	body := []byte(`{"name": "WireMock","tag": "v3"}`)
	statusCode, _, err := SendHttpPost(container, "/schema-match", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}

	invalidBody := `{"tag": "v2"}`
	statusCode, _, err = SendHttpPost(container, "/schema-match", strings.NewReader(invalidBody))
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 404 {
		t.Fatalf("expected HTTP-404 but got %d", statusCode)
	}
}

func TestV3MultipQueryValues(t *testing.T) {
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("v3", filepath.Join("testdata", "v3-url-with-multi-query-values.json")),
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

	statusCode, _, err := SendHttpGet(container, "/things?id=1&id=2&id=3", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
}

func TestV3PathTemplate(t *testing.T) {
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("v3", filepath.Join("testdata", "v3-path-template.json")),
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

	statusCode, out, err := SendHttpGet(container, "/v1/contacts/12345/addresses/99876", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
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

func TestV3Faker(t *testing.T) {
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("v3", filepath.Join("testdata", "200-v3-random.json")),
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

	statusCode, out, err := SendHttpGet(container, "/random", nil)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}

	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
	randomNumber, err := strconv.Atoi(out)
	if err != nil {
		t.Fatalf("expected a number but got %s", out)
	}

	if randomNumber < 1 || randomNumber > 9 {
		t.Fatalf("expected number between 1-9 but got %s", out)
	}
}

func TestV3Auth(t *testing.T) {
	ctx := context.Background()
	container, err := RunContainer(ctx,
		WithImage(defaultV3WireMockImage),
		WithMappingFile("v3", filepath.Join("testdata", "200-v3-basic-auth.json")),
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

	endpoint, err := GetURI(ctx, container)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"/basic-auth", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth("user", "pass")

	statusCode, _, err := sendTestRequest(t, req)
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}
	if statusCode != 200 {
		t.Fatalf("expected HTTP-200 but got %d", statusCode)
	}
}

func sendTestRequest(t *testing.T, req *http.Request) (int, string, error) {
	t.Helper()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(out), nil
}
