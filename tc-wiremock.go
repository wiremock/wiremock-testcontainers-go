package testcontainers_wiremock

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

const defaultWireMockImage = "docker.io/wiremock/wiremock"
const defaultWireMockVersion = "2.35"

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

//func WithVersion(version string) testcontainers.CustomizeRequestOption {
//	return func(req *testcontainers.GenericContainerRequest) {
//		req
//	}
//
//}
