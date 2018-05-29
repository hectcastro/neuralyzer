build:
	GOOS=linux go build -o main main.go

deploy:
	sam deploy --template-file packaged.yaml --stack-name Neuralyzer --capabilities CAPABILITY_IAM

package: build
	sam package --s3-bucket neuralyzer-global-config-us-east-1 --template-file template.yaml --output-template-file packaged.yaml

test:
	sam local generate-event schedule | sam local invoke NeuralyzerFunction

.PHONY: build deploy package test
