package cloudformationdeploy_test

import (
	cfd "cloudformationdeploy"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"gotest.tools/assert"
)

// Mon Jan 2 15:04:05 MST 2006
const layoutAWS = "2006-01-02T15:04:05.000000-07:00"


func TestPopulateData(t *testing.T) {
	var countState = 0;
	// make and configure a mocked DeployInterface
	mockedDeployInterface := &cfd.DeployInterfaceMock{
		CreateStackFunc: func(ctx context.Context, params *cloudformation.CreateStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateStackOutput, error) {
				panic("mock out the CreateStack method")
		},
		DeleteStackFunc: func(ctx context.Context, params *cloudformation.DeleteStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error) {
				panic("mock out the DeleteStack method")
		},
		DescribeStackEventsFunc: func(ctx context.Context, params *cloudformation.DescribeStackEventsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackEventsOutput, error) {
			var Events cloudformation.DescribeStackEventsOutput;
			var data []byte;
			var err error;
			if( countState == 0){
				data, err = ioutil.ReadFile("test/events1.json")
				countState++
			}else if( countState == 1){
				data, err = ioutil.ReadFile("test/events2.json")
				countState++
			}
			if err != nil {
				fmt.Println("File reading error", err)
			}
			json.Unmarshal(data, &Events);

			return &Events,nil;

		},
	}

	dataPre := map[string]cfd.CloudFormationResource{
			"testcfn" : {
				LogicalResourceId: "testfncn",
				Type: "AWS::CloudFormation::Stack",
			},
			"MyTopic" : {
				LogicalResourceId: "MyTopic",
				Type: "AWS::SNS::Topic",
			},
			"NotMyTopic" : {
				LogicalResourceId: "NotMyTopic",
				Type: "AWS::SNS::Topic",
			},
	}

	// Timestamps from events1.json
	t1, _ := time.Parse(layoutAWS, "2020-11-06T10:55:46.074000+00:00");
	t2, _ := time.Parse(layoutAWS, "2020-11-06T10:55:49.190000+00:00");
	t3, _ := time.Parse(layoutAWS, "2020-11-06T10:55:49.187000+00:00")
	dataTarget1 := map[string]cfd.CloudFormationResource{
			"testcfn" : {
				LogicalResourceId: "testfncn",
				PhysicalResourceId: "",
				Status: "CREATE_IN_PROGRESS",
				Type: "AWS::CloudFormation::Stack",
				Timestamp: t1,
			},
			"MyTopic" : {
				LogicalResourceId: "MyTopic",
				Status: "CREATE_IN_PROGRESS",
				Type: "AWS::SNS::Topic",
				Timestamp: t2,
			},
			"NotMyTopic" : {
				LogicalResourceId: "NotMyTopic",
				Status: "CREATE_IN_PROGRESS",
				Type: "AWS::SNS::Topic",
				Timestamp: t3,
			},
	}

	// dataTarget2 := map[string]cfd.CloudFormationResource{
	// 		"testcfn" : {
	// 			LogicalResourceId: "testfncn",
	// 			PhysicalResourceId: "",
	// 			Status: "CREATE_COMPLETE",
	// 			Type: "AWS::CloudFormation::Stack",
	// 		},
	// 		"MyTopic" : {
	// 			LogicalResourceId: "MyTopic",
	// 			Status: "CREATE_COMPLETE",
	// 			Type: "AWS::SNS::Topic",
	// 		},
	// 		"NotMyTopic" : {
	// 			LogicalResourceId: "NotMyTopic",
	// 			Status:"CREATE_COMPLETE",
	// 			Type: "AWS::SNS::Topic",
	// 		},
	// }

	data1 := cfd.PopulateData(mockedDeployInterface, "TestStack", dataPre);
	assert.DeepEqual(t,dataTarget1, data1)

	// data2 := cfd.PopulateData(mockedDeployInterface, "TestStack", data1);
	// assert.DeepEqual(t,dataTarget2, data2)
}