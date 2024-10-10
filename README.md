# Application Deployment Guide

## Overview

This guide provides instructions for building and deploying the load balancer application. The application consists of an app server and a load balancer, managed using [AWS CDK](https://aws.amazon.com/cdk/)

## Prerequisites

Before you can run the application, ensure you have the following installed:

1. **Go**: Version 1.16 or later.

   - Download from [golang.org/dl](https://golang.org/dl/)
   - [Installation guide](https://golang.org/doc/install)

2. **Node.js**: Version 14 or later.

   - Download from [nodejs.org](https://nodejs.org/)
   - [Installation guide](https://nodejs.org/en/download/package-manager/)

3. **AWS CDK**: Version 2.x.

   - Install it globally using npm:
     ```bash
     npm install -g aws-cdk
     ```
   - [AWS CDK Getting Started](https://docs.aws.amazon.com/cdk/v2/guide/getting_started.html)

4. **AWS CLI**: Version 2.x and configured with appropriate credentials.

   - Installation instructions at [aws.amazon.com/cli](https://aws.amazon.com/cli/)
   - [Configuration guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html)
   - Make sure you have configured your AWS CLI with:
     ```bash
     aws configure
     ```

5. **ZIP utility**: To create ZIP archives.

   - This is usually pre-installed on most Unix-like systems.
   - For Windows, you can use [7-Zip](https://www.7-zip.org/)

6. **Configure .env** :
   - Make sure you have the following envs setup for **cdk** to run :
     ```bash
     CDK_DEFAULT_ACCOUNT=<AWS_ACCOUNT_ID>
     CDK_DEFAULT_REGION=<REGION>
     ```

- Refer [this](https://docs.aws.amazon.com/cdk/v2/guide/environments.html) in case of issues.

7. Create config.yaml :
   - Make a config.yaml file in the **deploy_scripts** directory.
   - Example :

```bash
worker: 3                 # minimum number of worker nodes
pool: 30                  # number of incoming request
stats-dir: /tmp/stats     # location of the stats directory
avg-delay: 350            # average delay in ms
failure: 20               # failure percentage
max_workers: 6            # maximum number of worker nodes
```

## Directory Structure

```bash
GoBalance/
├── app_server/
│   ├── main.go
│   ├── other files and directories
│   └── //...
│
├── load_balancer/
│   ├── main.go
│   ├── other files and directories
│   └── //...
│
└── deploy_scripts/
│   ├── assests/
│   ├── deploy_scripts.go
|   ├── config.yaml
│   ├── other files and directories
│   └── //...
├── Makefile
├── .env
└── other files
```

## Makefile Commands

The provided Makefile contains the following commands:

- `build`: Compiles the Go applications(web server and load balancer) for the Linux OS and creates ZIP archives.
- `deploy`: Deploys the application using AWS CDK.
- `destroy`: Destroys the AWS resources created by the CDK.
- `clean`: Removes compiled binaries and ZIP files.
- `all`: Runs destroy, clean, build, and deploy in sequence, ignoring errors from destroy.

## Running the Application

### Step 1: Build the Application

To compile the applications and create ZIP files, run:

```bash
make build
```

### Step 2: Deploy the Application

To deploy the application to AWS, run:

```bash
make deploy
```

### Step 3: Clean Up (Optional)

To clean up the generated files and resources, you can run:

```bash
make clean
```

### Step 4: Destroy Resources

If you want to destroy the deployed resources, run:

```bash
make destroy
```

### Step 5: Run All Steps

To run all steps (destroy, clean, build, and deploy) in one command, run:

```bash
make all
```

**Note**: The `make all` command will continue executing even if the `make destroy` step fails.

## Troubleshooting

- Ensure you have the correct [permissions set in your AWS account](https://docs.aws.amazon.com/cdk/v2/guide/permissions.html).
- Check your Go and Node.js versions if you encounter compatibility issues.
- If you face issues with CDK, try running `cdk bootstrap` to set up the required resources. See [CDK Bootstrapping](https://docs.aws.amazon.com/cdk/v2/guide/bootstrapping.html) for more information.
