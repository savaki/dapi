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

import "github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"

type config struct {
	api         rdsdataserviceiface.RDSDataServiceAPI
	database    string
	secretARN   string
	resourceARN string
}

const timeFormat = "2006-01-02 15:04:05.999999999"
