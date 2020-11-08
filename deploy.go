package cloudformationdeploy

import (
	"context"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"fmt"
	"github.com/awslabs/goformation/v4/cloudformation"
	"time"
	"github.com/alexeyco/simpletable"
	// tm "github.com/buger/goterm"

)

// CloudFormationResource holder for status
type CloudFormationResource struct {
	LogicalResourceId string
	PhysicalResourceId string
	Status string
	Type string
	Timestamp time.Time
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
func ShowStatus(client DeployInterface, name string, template *cloudformation.Template){
	
	//maxRows := len(template.Resources)+1;
	
    data := map[string]CloudFormationResource{}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "#"},
			{Align: simpletable.AlignLeft, Text: "STATE"},
			{Align: simpletable.AlignLeft, Text: "Logical"},
			{Align: simpletable.AlignLeft, Text: "Type"},
		},
		
	}


	i := 1
	for k, v := range template.Resources {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", i)},
			{Align: simpletable.AlignLeft, Text: "state"},
			{Align: simpletable.AlignLeft, Text: k},
			{Align: simpletable.AlignLeft, Text: v.AWSCloudFormationType()},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		i = i+1
		item := &CloudFormationResource{
			LogicalResourceId: k,
			PhysicalResourceId: "",
			Status: "-",
			Type: v.AWSCloudFormationType(),
		}
		data[k] = *item;
	}
	
	// table.SetStyle(simpletable.StyleCompactClassic)
	// tm.Clear()
	// tm.MoveCursor(1,1)

	// tm.Println(table.String())
	// tm.Flush()

	i = 0;
	for id, v := range template.Resources {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", i)},
			{Align: simpletable.AlignLeft, Text: data[id].Status},
			{Align: simpletable.AlignLeft, Text: id},
			{Align: simpletable.AlignLeft, Text: v.AWSCloudFormationType()},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		i = i+1;
	}
	// tm.Clear()
	// tm.MoveCursor(1,1)
	// tm.Println(table.String())
	// tm.Flush()
	
	doEvery(2000*time.Millisecond,client,name,template, table,data)
	
}

func doEvery(d time.Duration, client DeployInterface, name string,template  *cloudformation.Template, table *simpletable.Table, data map[string]CloudFormationResource) {
	for range time.Tick(d) {
		PopulateData(client, name, data)
	}

}

// PopulateData update status from describe call
func PopulateData(client DeployInterface, name string,data map[string]CloudFormationResource)( map[string]CloudFormationResource){
	params := &cfn.DescribeStackEventsInput{
		StackName: &name,
	}
	output, error := client.DescribeStackEvents(context.TODO(), params)
	if( error != nil){
		panic(error)
	}

	// Update Status and Timestamp if newer
	for i := 0; i < len(output.StackEvents); i++ {
		
		event := *output.StackEvents[i];		
		item := data[*event.LogicalResourceId]

		if( event.Timestamp.After(item.Timestamp) ){
			item.Status = string(event.ResourceStatus);
			item.Timestamp = *event.Timestamp;
			data[*event.LogicalResourceId] = item;
			
		}
		
	}
	return data;

}
