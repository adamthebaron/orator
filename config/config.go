package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// site configuration from config.yaml
type siteconfig struct {
	title       string
	subtitle    string
	description string
	keywords    string
	author      string
	extra       map[string]interface{}
}

func (sc *SiteConfig) ReadConfig(fpath string) {
	contents, err := ioutil.ReadFile(fpath)

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(contents, sc)
}
