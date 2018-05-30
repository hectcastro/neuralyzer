# neuralyzer

A Go based Amazon Lambda function to periodically delete stale tweets, complete with:

- Deployment configuration via Serverless Application Model (SAM)
- Secure Twitter API credential lookup via Simple Systems Manager (SSM) Parameter Store
- Replica execution environment testing via [aws-sam-cli](https://github.com/awslabs/aws-sam-cli)

*Heavily* inspired by [ephemeral](https://github.com/vickylai/ephemeral), by **@vickylai**. All of the functionality was redone from scratch to get familiar with using Go and the AWS SDK for Go. Also, as an excuse to take a closer look at SAM.

