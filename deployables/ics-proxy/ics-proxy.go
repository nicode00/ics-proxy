package icsproxy

import (
	"log"
	"net/http"
	"strings"

	"github.com/nicompile/ics-proxy/internal/config"
	"github.com/nicompile/ics-proxy/internal/parser"
	"github.com/nicompile/infra-library-go/pkg/serverlessapi"
)

func proxy(request serverlessapi.Request) serverlessapi.Response {
	conf := config.Get()
	resp, err := http.Get(conf.Url)
	if err != nil {
		log.Fatal(err)
	}

	body := strings.Builder{}
	p := parser.New(resp.Body, &body)
	err = p.Parse()
	if err != nil {
		panic(err)
	}

	return serverlessapi.Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/calendar"},
		Body:       body.String(),
	}
}

func GetEndpoints() []serverlessapi.Endpoint {
	return []serverlessapi.Endpoint{
		{
			Method: "GET",
			Path:   "/",
			Target: proxy,
			Timeout: serverlessapi.Timeout{
				Minutes: 0,
				Seconds: 30,
			},
		},
	}
}
