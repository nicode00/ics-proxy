package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nicode00/ics-proxy/internal/parser"
)

type config struct {
	Url string `json:"url"`
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(url())
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/calendar")

	parser.Trim(resp.Body, w)
}

func url() string {
	content, err := ioutil.ReadFile("res/config.json")
	if err != nil {
		panic(err)
	}

	var c config
	err = json.Unmarshal(content, &c)
	if err != nil {
		panic(err)
	}

	return c.Url
}
