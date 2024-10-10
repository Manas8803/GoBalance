package stack

import (
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"GoBalance/deploy_scripts/config"
	"GoBalance/deploy_scripts/resources/ec2"
	"GoBalance/deploy_scripts/resources/s3"
)

type DeployScriptsStackProps struct {
	awscdk.StackProps
}

func NewDeployScriptsStack(scope constructs.Construct, id string, props *DeployScriptsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := createVPC(stack)
	worker_sg := createWorkerSecurityGroup(stack, vpc)
	lb_sg := createLBSecurityGroup(stack, vpc)

	bucket := s3.CreateS3BucketAndUploadAssets(stack)
	bucketName := *bucket.BucketName()
	region := os.Getenv("CDK_DEFAULT_REGION")

	worker_assets := []string{
		fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, "app_server/app_server/app_server"),
	}

	var worker_instances []awsec2.Instance
	for i := 0; i < config.VMConfigs.MaxWorkers; i++ {
		worker_instance := ec2.CreateEC2Instance(stack, vpc, worker_sg, fmt.Sprintf("WorkerInstance%d", i+1), true, i+1, worker_assets, []awsec2.Instance{})
		worker_instances = append(worker_instances, worker_instance)
	}

	lb_assets := []string{
		fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, "load_balancer/load_balancer/load_balancer"),
		fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, "load_balancer/load_balancer/nodes.txt"),
	}

	ec2.CreateEC2Instance(stack, vpc, lb_sg, "LoadBalancerInstance", false, 0, lb_assets, worker_instances)

	return stack
}

func createVPC(stack awscdk.Stack) awsec2.IVpc {
	return awsec2.NewVpc(stack, jsii.String("GoBalanceVPC"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})
}

func createWorkerSecurityGroup(stack awscdk.Stack, vpc awsec2.IVpc) awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(stack, jsii.String("GoBalance_Worker-SG"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		AllowAllOutbound:  jsii.Bool(true),
		SecurityGroupName: jsii.String("GoBalance_Worker-SG"),
	})

	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("Allow SSH"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(8080)), jsii.String("Allow HTTP on PORT 8080"), jsii.Bool(false))

	return sg
}

func createLBSecurityGroup(stack awscdk.Stack, vpc awsec2.IVpc) awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(stack, jsii.String("GoBalance_LB-SG"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		AllowAllOutbound:  jsii.Bool(true),
		SecurityGroupName: jsii.String("GoBalance_LB-SG"),
	})

	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("Allow SSH"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(2000)), jsii.String("Allow HTTP on PORT 2000"), jsii.Bool(false))

	return sg
}