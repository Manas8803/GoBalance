package ec2

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"

	"GoBalance/deploy_scripts/config"
)

func CreateEC2Instance(stack awscdk.Stack, vpc awsec2.IVpc, sg awsec2.ISecurityGroup, name string, isWorker bool, workerID int, assets []string, worker_instances []awsec2.Instance) awsec2.Instance {
	var userDataScript string

	commonUserDataScript := `
		#!/bin/bash
		apt-get update
		apt-get install -y golang-go unzip
	`

	if isWorker {
		userDataScript = createWorkerUserDataScript(commonUserDataScript, assets[0], workerID)
	} else {
		userDataScript = createLBUserDataScript(commonUserDataScript, assets[0], worker_instances)
	}

	userData := awsec2.UserData_ForLinux(&awsec2.LinuxUserDataOptions{Shebang: jsii.String("#!/bin/bash")})
	userData.AddCommands(jsii.String(userDataScript))

	ubuntuImage := awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
		Name:   jsii.String("ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"),
		Owners: jsii.Strings("099720109477"),
	})

	instance := awsec2.NewInstance(stack, jsii.String(name), &awsec2.InstanceProps{
		InstanceType:  awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		MachineImage:  ubuntuImage,
		Vpc:           vpc,
		VpcSubnets:    &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
		SecurityGroup: sg,
		UserData:      userData,
		InstanceName:  jsii.String(name),
	})

	return instance
}

func createWorkerUserDataScript(commonUserDataScript, assetURL string, workerID int) string {
	return fmt.Sprintf(`
    %s
    mkdir -p /home/ubuntu/app
    
    # curl file from the given url
    curl -O %s

    # copy the file to /home/ubuntu/app
    cp app_server /home/ubuntu/app
    
    chmod +x /home/ubuntu/app/app_server

    # Create .env file with environment variables
    cat << EOF > /home/ubuntu/app/.env
WORKER_DIR="%s"
FAIL_PERCENT=%f
WORKER_ID=%d
EOF

    # Create stats directory
    mkdir -p /home/ubuntu/app/%s

    # Start the worker application
    cd /home/ubuntu/app
    ./app_server >> /home/ubuntu/app/app_server.log 2>&1 &
`, commonUserDataScript, assetURL, config.VMConfigs.StatsDir, config.VMConfigs.Failure, workerID, config.VMConfigs.StatsDir)
}

func createLBUserDataScript(commonUserDataScript, assetURL string, worker_instances []awsec2.Instance) string {
	var available_nodes strings.Builder
	var standby_nodes strings.Builder
	var all_nodes strings.Builder
	idx := 0
	for _, workerInstance := range worker_instances {
		idx++
		if idx <= config.VMConfigs.Worker {
			available_nodes.WriteString(fmt.Sprintf("%s\n", *workerInstance.InstancePublicIp()))
		} else {
			standby_nodes.WriteString(fmt.Sprintf("%s\n", *workerInstance.InstancePublicIp()))
		}
		all_nodes.WriteString(fmt.Sprintf("%s\n", *workerInstance.InstancePublicIp()))
	}

	return fmt.Sprintf(`
		%s
		mkdir -p /home/ubuntu/app
		
		# curl file from the given url for load_balancer
		curl -O %s
		
		# copy the load_balancer file to /home/ubuntu/app
		cp load_balancer /home/ubuntu/app

		# Create nodes.txt file with worker instance public IPs
		cat << EOF > /home/ubuntu/app/available_nodes.txt
%s
EOF
		cat << EOF > /home/ubuntu/app/standby_nodes.txt
%s
EOF
		cat << EOF > /home/ubuntu/app/all_nodes.txt
%s
EOF
		cat << EOF > /home/ubuntu/app/.env
POOL=%d
MAX_WORKER=%d
WORKER=%d
EOF
		# Provide appropriate permissions
		chmod +x /home/ubuntu/app/load_balancer
		chmod +rw /home/ubuntu/app/available_nodes.txt
		chmod +rw /home/ubuntu/app/standby_nodes.txt
		chmod +rw /home/ubuntu/app/.env
		
		# Start the load balancer application
		cd /home/ubuntu/app
		./load_balancer >> /home/ubuntu/app/load_balancer.log 2>&1 &
		curl -X POST %s -H 'Content-Type: application/json' -d '{"email": "%s","message": {"subject": "GoBalance Deployed","body": "<p style=\"color: black;\">The deployment process has been <span style=\"color: green;\">completed successfully</span>.</p><p style=\"color: black;\">Check your aws console and get the</p>"}}'
	`, commonUserDataScript, assetURL, available_nodes.String(), standby_nodes.String(), all_nodes.String(), config.VMConfigs.Pool, config.VMConfigs.MaxWorkers, config.VMConfigs.Worker, os.Getenv("MAIL_API"), os.Getenv("ADMIN_MAIL"))
}
