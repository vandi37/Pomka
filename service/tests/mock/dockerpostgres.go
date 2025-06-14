package mock

import (
	"fmt"

	"github.com/ory/dockertest/v3"
)

type DockerPool struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

type Config struct {
	User     string
	Name     string
	Password string
}

func PostgresUp(cfg Config) (*DockerPool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	resource, err := pool.Run("postgres", "latest", []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Password),
		fmt.Sprintf("POSTGRES_USER=%s", cfg.User),
		fmt.Sprintf("POSTGRES_DB=%s", cfg.Name),
	})
	if err != nil {
		return nil, err
	}

	return &DockerPool{pool: pool, resource: resource}, nil
}

func (m *DockerPool) PostgresDown() error {
	if err := m.pool.Purge(m.resource); err != nil {
		return err
	}

	return nil
}
