package dapi

import (
	"context"
	"fmt"
	"io"

	"database/sql/driver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type Result struct {
	output *rdsdataservice.ExecuteStatementOutput
}

func (r *Result) Columns() []string {
	var columns []string
	for _, meta := range r.output.ColumnMetadata {
		columns = append(columns, aws.StringValue(meta.Name))
	}
	return columns
}

func (r *Result) Close() error {
	return nil
}

func (r *Result) Next(dest []driver.Value) error {
	return io.EOF
}

func (r *Result) LastInsertId() (int64, error) {
	panic("implement me: LastInsertId")
}

func (r *Result) RowsAffected() (int64, error) {
	return aws.Int64Value(r.output.NumberOfRecordsUpdated), nil
}

func executeStatement(ctx context.Context, config *config, query string, args ...driver.NamedValue) (driver.Rows, error) {
	input := &rdsdataservice.ExecuteStatementInput{
		Database:              aws.String(config.database),
		IncludeResultMetadata: aws.Bool(true),
		ResourceArn:           aws.String(config.resourceARN),
		SecretArn:             aws.String(config.secretARN),
		Sql:                   aws.String(query),
	}

	for _, arg := range args {
		param := rdsdataservice.SqlParameter{
			Name:     aws.String(arg.Name),
			TypeHint: nil,
			Value:    valueOf(arg.Value),
		}
		input.Parameters = append(input.Parameters, &param)
	}

	output, err := config.api.ExecuteStatementWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return &Result{
		output: output,
	}, nil
}
