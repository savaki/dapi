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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"time"
)

type Stmt struct {
	ctx    context.Context
	config *config
	query  string
}

func (s *Stmt) Close() error {
	return nil
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (s *Stmt) NumInput() int {
	return -1
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("implement me: Exec")
}

func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return executeStatement(ctx, s.config, s.query, "", args...)
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	panic("implement me: Query")
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	panic("implement me: QueryContext (Stmt)")
}

func newStmt(ctx context.Context, config *config, query string) *Stmt {
	return &Stmt{
		ctx:    ctx,
		config: config,
		query:  query,
	}
}

func asField(value driver.Value) *rdsdataservice.Field {
	switch v := value.(type) {
	case int64:
		return &rdsdataservice.Field{LongValue: aws.Int64(v)}
	case float64:
		return &rdsdataservice.Field{DoubleValue: aws.Float64(v)}
	case bool:
		return &rdsdataservice.Field{BooleanValue: aws.Bool(v)}
	case []byte:
		return &rdsdataservice.Field{BlobValue: v}
	case string:
		return &rdsdataservice.Field{StringValue: aws.String(v)}
	case time.Time:
		s := v.Format("2006-01-02 15:04:05")
		return &rdsdataservice.Field{StringValue: aws.String(s)}
	default:
		return &rdsdataservice.Field{IsNull: aws.Bool(true)}
	}
}

func valueOf(field *rdsdataservice.Field) driver.Value {
	switch {
	case field.ArrayValue != nil:
		switch {
		case field.ArrayValue.BooleanValues != nil:
			return aws.BoolValueSlice(field.ArrayValue.BooleanValues)
		case field.ArrayValue.DoubleValues != nil:
			return aws.Float64ValueSlice(field.ArrayValue.DoubleValues)
		case field.ArrayValue.LongValues != nil:
			return aws.Int64ValueSlice(field.ArrayValue.LongValues)
		case field.ArrayValue.StringValues != nil:
			return aws.StringValueSlice(field.ArrayValue.StringValues)
		default:
			return nil
		}
	case field.BlobValue != nil:
		return field.BlobValue
	case field.BooleanValue != nil:
		return *field.BooleanValue
	case field.DoubleValue != nil:
		return *field.DoubleValue
	case field.LongValue != nil:
		return *field.LongValue
	case field.StringValue != nil:
		return *field.StringValue
	default:
		return nil
	}
}
