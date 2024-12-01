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

6. **Environment Setup**
   - Create `.env.local` file with placeholder environment variables
   - Copy `.env.local` to `.env` and fill in actual values
     ```bash
     cp .env.local .env
     ```
   - Modify `.env` with your specific configuration values

- Refer [this](https://docs.aws.amazon.com/cdk/v2/guide/environments.html) in case of issues.

7. Create config.yaml :
   - Make a config.yaml file in the **deploy_scripts** directory.
   - Example :

```bash
worker: 3                 # minimum number of worker nodes
max_workers: 6            # maximum number of worker nodes
pool: 30                  # number of allowed concurrent incoming requests
stats-dir: /tmp/stats     # location of the stats directory
failure: 20               # failure percentage
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

The Makefile provides the following commands:

- `build`: Compiles Go applications for Linux (amd64) and creates ZIP archives
  - Builds app_server and load_balancer binaries
  - Creates ZIP archives for both components
- `deploy`: Deploys the application using AWS CDK without requiring manual approval
- `destroy`: Destroys AWS resources created by CDK
- `clean`: Removes compiled binaries and ZIP files
- `test`: Runs the test script
- `all`: Performs a complete deployment cycle:
  1. Attempts to destroy existing resources
  2. Cleans up files
  3. Builds the application
  4. Deploys to AWS

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

### Step 5: Running Tests

To execute the test script:

```bash
make test
```

## Step 6: Run All Steps

To run all steps (destroy, clean, build, and deploy) in one command, run:

```bash
make all
```

**Note**: The `make all` command will continue executing even if the `make destroy` step fails.

## Troubleshooting

- Ensure you have the correct [permissions set in your AWS account](https://docs.aws.amazon.com/cdk/v2/guide/permissions.html).
- Check your Go and Node.js versions if you encounter compatibility issues.
- If you face issues with CDK, try running `cdk bootstrap` to set up the required resources. See [CDK Bootstrapping](https://docs.aws.amazon.com/cdk/v2/guide/bootstrapping.html) for more information.
