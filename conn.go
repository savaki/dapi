package dapi

import (
	"context"
	"database/sql/driver"
	"fmt"
)

type Conn struct {
	ctx    context.Context
	config *config
}

func newConn(ctx context.Context, config *config) *Conn {
	return &Conn{
		ctx:    ctx,
		config: config,
	}
}

func (c *Conn) Begin() (driver.Tx, error) {
	panic("implement me: Begin")
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	panic("implement me: ExecContext")
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	fmt.Println("Prepare")
	return newStmt(context.Background(), c.config, query), nil
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	fmt.Println("PrepareContext")
	return newStmt(ctx, c.config, query), nil
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return executeStatement(ctx, c.config, query, args...)
}
