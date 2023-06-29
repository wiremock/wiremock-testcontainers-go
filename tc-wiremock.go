package testcontainers_wiremock

import (
	"context"
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

const defaultWireMockImage = "docker.io/wiremock/wiremock"
const defaultWireMockVersion = "2.35"
const defaultPort = 8080

type WireMockContainer struct {
	testcontainers.Container
	version string
}

type WireMockExtension struct {
	testcontainers.Container
	id        string
	classname string
	jarPath   string
}

// RunContainer creates an instance of the postgres container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*WireMockContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        defaultWireMockImage + ":" + defaultWireMockVersion,
		ExposedPorts: []string{"8080/tcp"},
		Cmd:          []string{""},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	req.Cmd = append(req.Cmd, "--disable-banner")

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &WireMockContainer{Container: container}, nil
}

func WithMappingFile(id string, filePath string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		cfgFile := testcontainers.ContainerFile{
			HostFilePath:      filePath,
			ContainerFilePath: "/home/wiremock/mappings/" + id + ".json",
			FileMode:          0755,
		}

		req.Files = append(req.Files, cfgFile)
	}

}

func WithFile(name string, filePath string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		cfgFile := testcontainers.ContainerFile{
			HostFilePath:      filePath,
			ContainerFilePath: "/home/wiremock/__files/",
			FileMode:          0755,
		}

		req.Files = append(req.Files, cfgFile)
	}

}

func GetURI(ctx context.Context, container testcontainers.Container) (string, error) {
	hostIP, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(strconv.Itoa(defaultPort)))
	if err != nil {
		return "", err
	}

	return hostIP + ":" + mappedPort.Port(), nil
}
