# WireMock Module for Testcontainers Go

## Note

The Testcontainers module does not work with the official image at the moment,
because the Mappings and Files directories are not initialized there.
Use a custom image on the top of it, see `Dockerfile`.

## Usage

```golang

import (
 "context"
 "io/ioutil"
 "net/http"
 "testing"

 _ "github.com/lib/pq"
 "github.com/pkg/errors"
 "github.com/oleg-nenashev/wiremock-testcontainers-go"
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
```
