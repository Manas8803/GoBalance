package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"

	"GoBalance/deploy_scripts/config"
	"GoBalance/deploy_scripts/stack"
)

var StackName = "GoBalanceStack"

func main() {
	config.LoadConfig()
	defer jsii.Close()
	app := awscdk.NewApp(nil)

	stack.NewDeployScriptsStack(app, StackName, &stack.DeployScriptsStackProps{
		StackProps: awscdk.StackProps{
			StackName: jsii.String(StackName),
			Env:       env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading env : ", err)
		log.Println("Continuing..")
		return nil
	}
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}