package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
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
	stack := awscdk.NewStack(scope, &id, &sprops)

	//get the current account and region
	account_id := *awscdk.Aws_ACCOUNT_ID()
	// Create a VPC with two public and two private subnets
	vpc := awsec2.NewVpc(stack, jsii.String("txn-vpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})

	//create ecs cluster
	cluster := awsecs.NewCluster(stack, jsii.String("txn-cluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	executionRole := awsiam.NewRole(stack, jsii.String("txn-execution-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AmazonECSTaskExecutionRolePolicy")),
		},
	})

	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("txn-sg"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		AllowAllOutbound:  jsii.Bool(true),
		SecurityGroupName: jsii.String("txn-sg"),
	})

	securityGroup.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(8080)),
		jsii.String("Allow traffic from anywhere on port 8080"),
		nil,
	)

	// Create Task Definition
	taskDef := awsecs.NewFargateTaskDefinition(stack, jsii.String("txn-task-def"), &awsecs.FargateTaskDefinitionProps{
		Cpu:            jsii.Number(256),
		MemoryLimitMiB: jsii.Number(512),
		ExecutionRole:  executionRole,
	})

	// concatenate the account number and region to create the image name
	repo_name := "transaction-repository"
	image_name := jsii.String(fmt.Sprintf("%s.dkr.ecr.us-east-2.amazonaws.com/%s:latest", account_id, repo_name))

	container := taskDef.AddContainer(jsii.String("txn-container"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.ContainerImage_FromRegistry(image_name, &awsecs.RepositoryImageProps{}),
	})
	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(8080),
		Protocol:      awsecs.Protocol_TCP,
	})

	// create ecs service
	awsecs.NewFargateService(stack, jsii.String("transaction-service"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: taskDef,
		AssignPublicIp: jsii.Bool(true),
		DesiredCount:   jsii.Number(1),
		SecurityGroups: &[]awsec2.ISecurityGroup{securityGroup},
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewClusterStack(app, "TransactionServiceStack2", &ClusterStackProps{
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
