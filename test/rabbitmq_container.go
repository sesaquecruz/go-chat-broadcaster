package test

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RabbitMQContainer struct {
	container testcontainers.Container
	url       string
}

func NewRabbitMQContainer(ctx context.Context, configDir string) *RabbitMQContainer {
	rabbitmq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.12.3-management-alpine",
		ExposedPorts: []string{"5672/tcp"},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{
					HostPath: fmt.Sprintf("%s/rabbitmq.config", configDir),
				},
				Target: "/etc/rabbitmq/rabbitmq.config",
			},
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{
					HostPath: fmt.Sprintf("%s/rabbitmq.json", configDir),
				},
				Target: "/etc/rabbitmq/definitions.json",
			},
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: rabbitmq,
			Started:          true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "5672/tcp")
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQContainer{
		container: container,
		url:       fmt.Sprintf("amqp://guest:guest@%s:%s", host, port.Port()),
	}
}

func (c *RabbitMQContainer) Url() string {
	return c.url
}
