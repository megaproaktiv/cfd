package cloudformationdeploy

import (
	"context"
	"fmt"
	"time"

	"github.com/alexeyco/simpletable"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation"
	tm "github.com/buger/goterm"
)

const CREATE_COMPLETE="CREATE_COMPLETE"
const CREATE_IN_PROGRESS = "CREATE_IN_PROGRESS"

const (
	ColorDefault = "\x1b[39m"

	ColorRed   = "\x1b[91m"
	ColorGreen = "\x1b[32m"
	ColorBlue  = "\x1b[94m"
	ColorGray  = "\x1b[90m"
)

// CloudFormationResource holder for status
type CloudFormationResource struct {
	LogicalResourceID string
	PhysicalResourceID string
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
	
	// Prepopulate
	
    data := map[string]CloudFormationResource{}
	i := 1
	for k, v := range template.Resources {
		i = i+1
		item := &CloudFormationResource{
			LogicalResourceID: k,
			PhysicalResourceID: "",
			Status: "-",
			Type: v.AWSCloudFormationType(),
		}
		data[k] = *item;
	}
	
	// Draw
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "ID"},
			{Align: simpletable.AlignLeft, Text: "State"},
			{Align: simpletable.AlignLeft, Text: "Type"},
		},
		
	}
	table.SetStyle(simpletable.StyleCompactLite)
	
	first := true
	for !IsStackCompleted(data){
		tm.Clear()
		tm.MoveCursor(1,1)
		data = PopulateData(client, name, data);
		i = 0;
		var statustext string
		for id, v := range data {
			if( v.Status == CREATE_COMPLETE){
				statustext = green(CREATE_COMPLETE)
			}else {		
				statustext = gray(v.Status)
			}

			r := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: id},
				{Align: simpletable.AlignLeft, Text: statustext},
				{Align: simpletable.AlignLeft, Text: v.Type},
			}
			if !first {
				table.Body.Cells[i]=r
			}else{
				table.Body.Cells = append(table.Body.Cells, r)
			}
			i = i+1;
		}
		first = false
		tm.Println(table.String())
		tm.Flush()
		time.Sleep(1 * time.Second) 
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

// IsStackCompleted check for everything "completed"
func IsStackCompleted(data map[string]CloudFormationResource) bool {
	for _, value := range data {
		if(value.Status != "CREATE_COMPLETE"){
			return false
		}
	}
	return true;
}

func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

func gray(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGray, s, ColorDefault)
}