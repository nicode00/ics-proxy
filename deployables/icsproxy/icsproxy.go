package icsproxy

import (
	"log"
	"net/http"
	"strings"

	"github.com/nicompile/ics-proxy/internal/config"
	"github.com/nicompile/ics-proxy/internal/parser"
	"github.com/nicompile/infra-library-go/pkg/serverlessfunction"
)

func Function(request serverlessfunction.Request) serverlessfunction.Response {
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

	return serverlessfunction.Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/calendar"},
		Body:       body.String(),
	}
}
