package testutils

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/traPtitech/traPortfolio/internal/pkgs/config"
)

func RunMySQLContainerOnDocker(user, pass, host string) (db *sql.DB, port int, closeFunc func() error, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, 0, nil, fmt.Errorf("create docker pool: %w", err)
	}

	pool.MaxWait = 1 * time.Minute

	if err := pool.Client.Ping(); err != nil {
		return nil, 0, nil, fmt.Errorf("ping docker: %w", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mariadb",
		Tag:        "10.6.4",
		Env:        []string{"MYSQL_ROOT_PASSWORD=" + pass},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, 0, nil, fmt.Errorf("run docker container: %w", err)
	}

	if err := resource.Expire(60); err != nil {
		return nil, 0, nil, fmt.Errorf("expire docker container: %w", err)
	}

	port, _ = strconv.Atoi(resource.GetPort("3306/tcp"))
	sqlConfig := config.SQLConfig{
		User: user,
		Pass: pass,
		Host: host,
		Name: "",
		Port: port,
	}

	if err := pool.Retry(func() error {
		_db, err := sql.Open("mysql", sqlConfig.DsnConfig().FormatDSN())
		if err != nil {
			return fmt.Errorf("open mysql: %w", err)
		}

		if err := _db.Ping(); err != nil {
			return fmt.Errorf("ping mysql: %w", err)
		}

		db = _db

		return nil
	}); err != nil {
		return nil, 0, nil, fmt.Errorf("retry: %w", err)
	}

	closeFunc = func() error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close mysql: %w", err)
		}

		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("purge docker container: %w", err)
		}

		return nil
	}

	return
}
