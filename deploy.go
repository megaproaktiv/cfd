package cloudformationdeploy

import (
	"context"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"fmt"
	// "github.com/jroimartin/gocui"
	"github.com/awslabs/goformation/v4/cloudformation"
	"time"
)

// CloudFormationResource holder for status
type CloudFormationResource struct {
	Name string
	Type string
	Status string
}

//go:generate moq -out deploy_moq_test.go . DeployInterface

// DeployInterface all deployment functions
type DeployInterface interface {
	CreateStack(ctx context.Context, params *cfn.CreateStackInput, optFns ...func(*cfn.Options)) (*cfn.CreateStackOutput, error)
	DescribeStackEvents(ctx context.Context, params *cfn.DescribeStackEventsInput, optFns ...func(*cfn.Options)) (*cfn.DescribeStackEventsOutput, error)
	DeleteStack(ctx context.Context, params *cfn.DeleteStackInput, optFns ...func(*cfn.Options)) (*cfn.DeleteStackOutput, error)


}

// CreateStack first time stack creation
func CreateStack(client DeployInterface,name string, template *cloudformation.Template){
	stack, _ := template.YAML()
	templateBody := string(stack)

	params := &cfn.CreateStackInput{
		StackName: &name,
		TemplateBody: &templateBody,
	}

	client.CreateStack(context.TODO(),params)
}

// DeleteStack first time stack creation
func DeleteStack(client DeployInterface,name string){

	params := &cfn.DeleteStackInput{
		StackName: &name,
	}

	client.DeleteStack(context.TODO(),params)
}


// ShowStatus status of stack
func ShowStatus(client DeployInterface, name string){
	
	doEvery(2000*time.Millisecond,client,name)
	
}

func showSingleStatus(client DeployInterface, name string){
	params := &cfn.DescribeStackEventsInput{
		StackName: &name,
	}
	output, error := client.DescribeStackEvents(context.TODO(), params)
	if( error != nil){
		panic(error)
	}

	for i := 0; i < len(output.StackEvents); i++ {
		event := *output.StackEvents[i];
		
		fmt.Println(*event.LogicalResourceId, " : ", event.ResourceStatus, ": ", *event.ResourceType);
		
	}
	fmt.Println("==========");

}

func doEvery(d time.Duration, client DeployInterface, name string) {
	for range time.Tick(d) {
		showSingleStatus(client, name)
	}
}