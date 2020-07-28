package dapi

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type Tx struct {
	context context.Context
	config  *config
	conn    *Conn
}

func (t *Tx) Commit() error {
	input := rdsdataservice.CommitTransactionInput{
		ResourceArn:   aws.String(t.config.resourceARN),
		SecretArn:     aws.String(t.config.secretARN),
		TransactionId: aws.String(t.conn.transactionID),
	}
	if _, err := t.config.api.CommitTransactionWithContext(t.context, &input); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	t.conn.transactionID = ""

	return nil
}

func (t *Tx) Rollback() error {
	input := rdsdataservice.RollbackTransactionInput{
		ResourceArn:   aws.String(t.config.resourceARN),
		SecretArn:     aws.String(t.config.secretARN),
		TransactionId: aws.String(t.conn.transactionID),
	}
	if _, err := t.config.api.RollbackTransactionWithContext(t.context, &input); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	t.conn.transactionID = ""

	return nil
}
