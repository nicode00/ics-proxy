package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nicompile/ics-proxy/internal/config"
	"github.com/nicompile/ics-proxy/internal/parser"
)

var conf config.Config

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(conf.Url)
	if err != nil {
		panic(err)
	}

	body := strings.Builder{}

	p := parser.New(resp.Body, &body)
	err = p.Parse()
	if err != nil {
		panic(err)
	}

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
