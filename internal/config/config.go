package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Url string `json:"url"`
}

func Get() Config {
	content, err := ioutil.ReadFile("res/config.json")
	if err != nil {
		panic(err)
	}

	var c Config
	err = json.Unmarshal(content, &c)
	if err != nil {
		panic(err)
	}
	return c
}
