package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nicode00/ics-proxy/internal/config"
	"github.com/nicode00/ics-proxy/internal/parser"
)

var conf config.Config

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(conf.Url)
	if err != nil {
		panic(err)
	}

	body := strings.Builder{}

	parser.Trim(resp.Body, &body)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/calendar"},
		Body:       body.String(),
	}, nil
}

func main() {
	conf = config.Get()
	lambda.Start(handler)
}
