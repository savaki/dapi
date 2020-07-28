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
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

func TestDriver_Open(t *testing.T) {
	var (
		accessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
		secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		sessionToken    = os.Getenv("AWS_SESSION_TOKEN")
		database        = os.Getenv("DATABASE")
		secretARN       = os.Getenv("SECRET_ARN")
		resourceARN     = os.Getenv("RESOURCE_ARN")
	)

	if accessKeyID == "" || secretAccessKey == "" || database == "" || secretARN == "" || resourceARN == "" {
		t.SkipNow()
	}

	var (
		ctx = context.Background()
		s   = session.Must(session.NewSession(aws.NewConfig().
			WithCredentials(credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken)).
			WithRegion("us-west-2")))
		api    = rdsdataservice.New(s)
		driver = New(api)
		dsn    = fmt.Sprintf("secret=%v resource=%v database=%v", secretARN, resourceARN, database)
	)

	sql.Register("mysql", driver)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	defer db.Close()

	t.Run("prepare", func(t *testing.T) {
		stmt, err := db.Prepare("select * from information_schema.tables")
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		defer stmt.Close()

		result, err := stmt.Exec()
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := rowsAffected, int64(0); got != want {
			t.Fatalf("got %v; want %v", got, want)
		}
	})

	t.Run("query", func(t *testing.T) {
		rows, err := db.QueryContext(ctx, "select * from information_schema.tables")
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		defer rows.Close()

		fmt.Println(rows.Columns())
		for rows.Next() {
		}
	})
}
