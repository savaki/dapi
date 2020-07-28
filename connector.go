// Copyright 2020 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
