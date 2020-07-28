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
