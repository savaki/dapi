package dapi

import (
	"context"
	"database/sql/driver"
)

type Connector struct {
	config *config
	driver driver.Driver
}

func newConnector(config *config, driver driver.Driver) *Connector {
	return &Connector{
		config: config,
		driver: driver,
	}
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return newConn(ctx, c.config), nil
}

func (c *Connector) Driver() driver.Driver {
	return c.driver
}
