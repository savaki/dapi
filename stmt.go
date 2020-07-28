package dapi

import (
	"context"
	"database/sql/driver"
	"fmt"
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
	input := &rdsdataservice.ExecuteStatementInput{
		Database:    aws.String(s.config.database),
		ResourceArn: aws.String(s.config.resourceARN),
		SecretArn:   aws.String(s.config.secretARN),
		Sql:         aws.String(s.query),
	}

	for _, arg := range args {
		param := rdsdataservice.SqlParameter{
			Name:     aws.String(arg.Name),
			TypeHint: nil,
			Value:    valueOf(arg.Value),
		}
		input.Parameters = append(input.Parameters, &param)
	}

	output, err := s.config.api.ExecuteStatementWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return &Result{
		output: output,
	}, nil
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

func valueOf(value driver.Value) *rdsdataservice.Field {
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
		s := v.Format(time.RFC3339)
		return &rdsdataservice.Field{StringValue: aws.String(s)}
	default:
		return &rdsdataservice.Field{IsNull: aws.Bool(true)}
	}
}
