package dapi

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type Conn struct {
	ctx           context.Context
	config        *config
	transactionID string // if in middle of tx
}

func newConn(ctx context.Context, config *config) *Conn {
	return &Conn{
		ctx:    ctx,
		config: config,
	}
}

func (c *Conn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	input := rdsdataservice.BeginTransactionInput{
		Database:    aws.String(c.config.database),
		ResourceArn: aws.String(c.config.resourceARN),
		SecretArn:   aws.String(c.config.secretARN),
	}
	output, err := c.config.api.BeginTransactionWithContext(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}

	c.transactionID = aws.StringValue(output.TransactionId)

	return &Tx{
		context: ctx,
		config:  c.config,
		conn:    c,
	}, nil
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return executeStatement(ctx, c.config, query, c.transactionID, args...)
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return newStmt(context.Background(), c.config, query), nil
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return newStmt(ctx, c.config, query), nil
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return executeStatement(ctx, c.config, query, c.transactionID, args...)
}
