# neuralyzer

## Table of Contents

- [Configuration](#configuration)
  - [SSM parameters](#ssm-parameters)
  - [Environment variables](#environment-variables)
- [Testing](#testing)
- [Deployment](#deployment)

A Go based Amazon Lambda function to periodically delete stale tweets, complete with:

- Deployment configuration via [Serverless Application Model (SAM)](https://github.com/awslabs/serverless-application-model)
- Secure Twitter API credential lookup via [Simple Systems Manager (SSM) Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
- Local Lambda execution environment testing via [aws-sam-cli](https://github.com/awslabs/aws-sam-cli)

*Heavily* inspired by [ephemeral](https://github.com/vickylai/ephemeral), by **@vickylai**. All of the tweet deleting functionality was redone from scratch as an exercise to get familiar with Go and the AWS SDK for Go. Also, it was a good excuse to take a closer look at SAM.

## Configuration

All of the application configuration is provided through a mixture SSM parameters and environment variables. The sensitive ones come through the former, and the rest through the latter.

### SSM parameters

- `/neuralyzer/twitter/accessToken`: Twitter API access token
- `/neuralyzer/twitter/accessTokenSecret`: Twitter API access token secret
- `/neuralyzer/twitter/consumerKey`: Twitter API consumer key
- `/neuralyzer/twitter/consumerKeySecret`: Twitter API consumer key secret

### Environment variables

- `NEURALYZER_TWEET_AGE_THRESHOLD`: Tweet age threshold for deleting tweets (default: `2190h`)

## Testing

Once the SSM parameters above exist, [install](https://github.com/awslabs/aws-sam-cli#installation) `sam` and its prerequisites (Docker, Python 2.7).

**Note**: Be careful not to accidentally delete all of your tweets by invoking the function locally!

Next, use the `test` target of the `Makefile` to build a Linux compatible binary and execute it within a container image that mimics the Amazon Lambda execution environment for Go.

```bash
$ make test -e AWS_PROFILE=personal
GOOS=linux go build -o main main.go
sam local generate-event schedule | sam local invoke NeuralyzerFunction
2018-05-29 23:05:59 Reading invoke payload from stdin (you can also pass it from file with --event)
2018-05-29 23:06:00 Invoking main (go1.x)
2018-05-29 23:06:00 Found credentials in shared credentials file: ~/.aws/credentials

Fetching lambci/lambda:go1.x Docker container image......
2018-05-29 23:06:00 Mounting /Users/hector/.go/src/github.com/hectcastro/neuralyzer as /var/task:ro inside runtime container
START RequestId: 747fc7fe-5471-13ca-a7ae-8cbeb0335031 Version: $LATEST
2018/05/30 03:06:17 ID [1001653560956858368] Hi, there.
END RequestId: 747fc7fe-5471-13ca-a7ae-8cbeb0335031
REPORT RequestId: 747fc7fe-5471-13ca-a7ae-8cbeb0335031  Duration: 479.12 ms     Billed Duration: 500 ms Memory Size: 128 MB     Max Memory Used: 14 MB
null
```

## Deployment

The deployment process works in two phases:

1. Package the function and upload it to S3
2. Deploy the function and its supporting infrastructure (referencing the artifact above)

Both of these steps are handled by the `sam` CLI, which seems to mostly wrap the `aws cloudformation` CLI.

First, use the `package` target of the `Makefile` to upload the binary to S3 and reference it in a newly created `packaged.yaml` CloudFormation configuration.

```bash
$ make package -e AWS_PROFILE=personal
GOOS=linux go build -o main main.go
sam package --s3-bucket neuralyzer-global-config-us-east-1 --template-file template.yaml --output-template-file packaged.yaml                                                    
Uploading to 7001c68762c2fcda61de373e0a30563d  29187040 / 29187040.0  (100.00%)
Successfully packaged artifacts and wrote output template to file packaged.yaml.
Execute the following command to deploy the packaged template
aws cloudformation deploy --template-file packaged.yaml --stack-name <YOUR STACK NAME>
```

Lastly, use the `deploy` target of the `Makefile` to apply the `packaged.yaml` CloudFormation configuration.

```bash
$ make deploy -e AWS_PROFILE=personal 
sam deploy --template-file packaged.yaml --stack-name Neuralyzer --capabilities CAPABILITY_IAM

Waiting for changeset to be created..
Waiting for stack create/update to complete
Successfully created/updated stack - Neuralyzer
```
