package main

import (
	// ~"fmt"
	"strconv"
	"time"
	"cloudformationdeploy"
	"github.com/aws/aws-sdk-go-v2/config"
	cfnservice "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/sns"


)

func main() {

	// Create a new CloudFormation template
	template := cloudformation.NewTemplate()

	// Create an Amazon SNS topic, with a unique name based off the current timestamp
	template.Resources["MyTopic"] = &sns.Topic{
		TopicName: "my-topic-" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	template.Resources["NotMyTopic"] = &sns.Topic{
		TopicName: "my-topic2-" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	
	template.Resources["Topic4"] = &sns.Topic{
		TopicName: "my-topic4-" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	
	template.Resources["Topic5"] = &sns.Topic{
		TopicName: "my-topic5" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	


	cfg, err := config.LoadDefaultConfig(config.WithRegion("eu-central-1"))
    if err != nil {
        panic("unable to load SDK config, " + err.Error())
	}	
	client := cfnservice.NewFromConfig(cfg);

	const stackname = "testcfn"
	cloudformationdeploy.CreateStack(client,stackname, template)

	cloudformationdeploy.ShowStatus(client,stackname,template);

	cloudformationdeploy.DeleteStack(client,stackname);
}
