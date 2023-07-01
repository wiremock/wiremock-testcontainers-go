package testcontainers_wiremock

import (
	"context"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const defaultWireMockImage = "docker.io/wiremock/wiremock"
const defaultWireMockVersion = "2.35.0-for-tc"
const defaultPort = "8080"

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
		ExposedPorts: []string{defaultPort + "/tcp"},
		Cmd:          []string{""},
		WaitingFor:   wait.ForHTTP("/__admin").WithPort(nat.Port(defaultPort)),
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
			ContainerFilePath: filepath.Join("/home/wiremock/mappings", id+".json"),
			FileMode:          0755,
		}

		req.Files = append(req.Files, cfgFile)
	}

}

func WithFile(name string, filePath string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		cfgFile := testcontainers.ContainerFile{
			HostFilePath:      filePath,
			ContainerFilePath: "/home/wiremock/__files/" + name,
			FileMode:          0755,
		}

		req.Files = append(req.Files, cfgFile)
	}
}

func WithImage(image string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Image = image
	}
}

func GetURI(ctx context.Context, container testcontainers.Container) (string, error) {
	hostIP, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(defaultPort))
	if err != nil {
		return "", err
	}

	return "http://" + hostIP + ":" + mappedPort.Port(), nil
}
