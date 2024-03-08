package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ClusterStackProps struct {
	awscdk.StackProps
}

func NewClusterStack(scope constructs.Construct, id string, props *ClusterStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	app := awscdk.NewStack(scope, &id, &sprops)

	// Create a new S3 bucket
	bucket := awss3.NewBucket(app, jsii.String("Reports"), &awss3.BucketProps{
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Create a new Lambda function from the Go program
	goLambda := awslambda.NewFunction(app, jsii.String("GoLambdaFunction"), &awslambda.FunctionProps{
		Code:       awslambda.Code_FromAsset(jsii.String("C:/Users/mrvil/projects/stori-challenge/report-lambda"), nil),
		Handler:    jsii.String("main"),
		Runtime:    awslambda.Runtime_PROVIDED_AL2023(),
		MemorySize: jsii.Number(128),
	})

	// Grant the Lambda function permissions to access the S3 bucket
	bucket.GrantRead(goLambda, nil)

	// Create an S3 notification to trigger the Lambda function on object creation events
	bucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3notifications.NewLambdaDestination(goLambda), &awss3.NotificationKeyFilter{
		Prefix: jsii.String("report/*"),
	})

	return app
}

func main() {
	app := awscdk.NewApp(nil)

	NewClusterStack(app, "ReportLambdaStack", &ClusterStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String("533267442883"),
		Region:  jsii.String("us-east-2"),
	}

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
