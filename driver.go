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
	"fmt"
	"strings"

	"database/sql/driver"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
)

type Driver struct {
	api rdsdataserviceiface.RDSDataServiceAPI
}

func New(api rdsdataserviceiface.RDSDataServiceAPI) driver.Driver {
	return &Driver{
		api: api,
	}
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	database, resourceARN, secretARN, ok := parseName(dsn)
	if !ok {
		return nil, fmt.Errorf("dsn must be of the form `secret={secret arn} resource={resource arn} database={database name}")
	}

	c := &config{
		api:         d.api,
		database:    database,
		resourceARN: resourceARN,
		secretARN:   secretARN,
	}

	return newConn(context.Background(), c), nil
}

func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	database, resourceARN, secretARN, ok := parseName(dsn)
	if !ok {
		return nil, fmt.Errorf("dsn must be of the form `secret={secret arn} resource={resource arn} database={database name}")
	}

	c := &config{
		api:         d.api,
		database:    database,
		resourceARN: resourceARN,
		secretARN:   secretARN,
	}

	return newConnector(c, d), nil
}

func parseName(name string) (database, resourceARN, secretARN string, ok bool) {
	for _, kv := range strings.Split(name, " ") {
		if parts := strings.SplitN(kv, "=", 2); len(parts) == 2 {
			switch k, v := parts[0], parts[1]; k {
			case "database":
				database = v
			case "resource":
				resourceARN = v
			case "secret":
				secretARN = v
			}
		}
	}

	return database, resourceARN, secretARN, database != "" && secretARN != "" && resourceARN != ""
}
