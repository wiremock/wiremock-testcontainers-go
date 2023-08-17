# Using the REST API client

The WireMock module includes [Go WireMock](https://github.com/wiremock/go-wiremock/)
client that you can retrieve from the `WireMockContainer` instance
and use for runtime configuration or observability data access.

```golang
	// Create Container
	ctx := context.Background()
	container, err := RunDefaultContainerAndStopOnCleanup(ctx)
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
```

See the example [here](https://github.com/wiremock/wiremock-testcontainers-go/tree/main/examples/using_api_client).
