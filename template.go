package cloudformationdeploy

import(
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/sns"
	"github.com/awslabs/goformation/v4/cloudformation/sqs"
	"strconv"
	"time"
)

// CreateTemplate - build the Cloudformation template
func CreateTemplate(name string)(*cloudformation.Template){
	template := cloudformation.NewTemplate()
	
	// Create an Amazon SNS topic, with a unique name based off the current timestamp
	template.Resources["GoFormationCompareTopic"] = &sns.Topic{
		TopicName: "CfdCdkCompareTopic" + strconv.FormatInt(time.Now().Unix(), 10),
		};
		
	template.Resources["GoFormationCompareQueue"] = &sqs.Queue{
		VisibilityTimeout: 300,
	};
		
	return template
}