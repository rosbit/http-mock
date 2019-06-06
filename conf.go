package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"fmt"
)

type StaticConf struct {
	Root  string
	Alias string
}

type CorsConf struct {
	AllowOrigin      string `yaml:"allow-origin"`
	AllowHeaders     string `yaml:"allow-headers"`
	AllowMethods     string `yaml:"allow-methods"`
	ExposeHeaders    string `yaml:"expose-headers"`
	AllowCredentials string `yaml:"allow-credentials"`
}

type Response struct {
	Status  int
	Headers map[string]string
	Cookies map[string]string
	// 1 of body types
	Body         string                        // output as it is
	TmplBody     string `yaml:"tmpl-body"`     // execute text/template
	FileBody     string `yaml:"file-body"`     // output file content
	RedirectBody string `yaml:"redirect-body"` // 302 redirect
}

type Action struct {
	Uri    string
	Method string
	Response
}

type ServiceConfT struct {
	Port               int    `yaml:"port"`
	StaticConf                `yaml:"static-home"`
	CorsConf                  `yaml:"cors"`
	DefaultContentType string `yaml:"default-content-type"`
	Actions          []Action
}

var ServiceConf ServiceConfT

func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var confFile string
	if err := getEnv("CONF_FILE", &confFile, true); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &ServiceConf); err != nil {
		return err
	}

	return nil
}

func DumpConf() {
	// fmt.Printf("conf: %#v\n", ServiceConf)
}
