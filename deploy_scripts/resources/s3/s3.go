package s3

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/jsii-runtime-go"
)

func CreateS3BucketAndUploadAssets(stack awscdk.Stack) awss3.Bucket {
	bucket := awss3.NewBucket(stack, jsii.String("AssetsBucket"), &awss3.BucketProps{
		BucketName:        jsii.String(fmt.Sprintf("gobalance-assets-bucket-%s", *stack.Account())),
		PublicReadAccess:  jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
		BlockPublicAccess: awss3.NewBlockPublicAccess(&awss3.BlockPublicAccessOptions{
			BlockPublicAcls:       jsii.Bool(false),
			RestrictPublicBuckets: jsii.Bool(false),
			IgnorePublicAcls:      jsii.Bool(false),
			BlockPublicPolicy:     jsii.Bool(false),
		}),
	})

	awss3deployment.NewBucketDeployment(stack, jsii.String("AppServerDeployment"), &awss3deployment.BucketDeploymentProps{
		Sources:              &[]awss3deployment.ISource{awss3deployment.Source_Asset(jsii.String("./assets/app_server.zip"), &awss3assets.AssetOptions{})},
		DestinationBucket:    bucket,
		DestinationKeyPrefix: jsii.String("app_server"),
	})

	awss3deployment.NewBucketDeployment(stack, jsii.String("LoadBalancerDeployment"), &awss3deployment.BucketDeploymentProps{
		Sources:              &[]awss3deployment.ISource{awss3deployment.Source_Asset(jsii.String("./assets/load_balancer.zip"), &awss3assets.AssetOptions{})},
		DestinationBucket:    bucket,
		DestinationKeyPrefix: jsii.String("load_balancer"),
	})

	bucketPolicy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:    jsii.Strings("s3:GetObject"),
		Resources:  jsii.Strings(fmt.Sprintf("%s/*", *bucket.BucketArn())),
		Principals: &[]awsiam.IPrincipal{awsiam.NewAnyPrincipal()},
	})
	bucket.AddToResourcePolicy(bucketPolicy)

	return bucket
}