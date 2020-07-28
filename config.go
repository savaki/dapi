package dapi

import "github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"

type config struct {
	api         rdsdataserviceiface.RDSDataServiceAPI
	database    string
	secretARN   string
	resourceARN string
}
