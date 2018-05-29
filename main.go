package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	twitterClient     *anaconda.TwitterApi
	tweetAgeThreshold = getEnv("NEURALYZER_TWEET_AGE_THRESHOLD", "2190h")
)

// getEnv does an environment variable lookup, but falls back
// to the provided default.
func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}

// getUserTimeline retrieves a list of tweets from the
// authenticated user's Twitter timeline.
func getUserTimeline() ([]anaconda.Tweet, error) {
	timeline, err := twitterClient.GetUserTimeline(
		url.Values{
			"count":       {"200"},
			"include_rts": {"true"},
		},
	)

	if err != nil {
		return make([]anaconda.Tweet, 0), err
	}

	return timeline, nil
}

// getParameter retrieves the value from from a parameter stored
// in the AWS SSM Parameter Store.
func getParameter(svc *ssm.SSM, path string) string {
	p, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name:           aws.String(path),
			WithDecryption: aws.Bool(true),
		},
	)

	if err != nil {
		log.Fatalf("Unable to read [%s] SSM parameter", path)
	}

	return aws.StringValue(p.Parameter.Value)
}

// init initializes the AWS SSM and Twitter client connections.
func init() {
	svc := ssm.New(session.New())

	twitterClient = anaconda.NewTwitterApiWithCredentials(
		getParameter(svc, "/neuralyzer/twitter/accessToken"),
		getParameter(svc, "/neuralyzer/twitter/accessTokenSecret"),
		getParameter(svc, "/neuralyzer/twitter/consumerKey"),
		getParameter(svc, "/neuralyzer/twitter/consumerKeySecret"),
	)
}

// HandleRequest is the Amazon Lambda handler function responsible
// for processing a CloudWatchEvent.
func HandleRequest(ctx context.Context, e events.CloudWatchEvent) {
	timeline, _ := getUserTimeline()
	parsedTweetAgeThreshold, _ := time.ParseDuration(tweetAgeThreshold)

	for _, tweet := range timeline {
		createdTime, err := tweet.CreatedAtTime()

		if err == nil {
			if time.Since(createdTime) > parsedTweetAgeThreshold {
				// deletedTweet, _ := twitterClient.DeleteTweet(tweet.Id, true)
				// log.Printf("ID [%d] %s", deletedTweet.Id, deletedTweet.Text)
				log.Printf("ID [%d] %s", tweet.Id, tweet.Text)
			}
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
