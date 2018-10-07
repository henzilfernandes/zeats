package config

import (
	"flag"
	"os"
	"io/ioutil"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var (
	cfg        *Config
	configPath string
)

// method will be called automatically when the package is loaded
func init() {
	flag.StringVar(&configPath, "config", os.Getenv("GOPATH")+"/src/zeats/config.yaml",
		"config file path")
	flag.Parse()

	if cfg != nil {
		return
	}

	cfg = Vals()
}

func Vals() *Config {
	if cfg != nil {
		return cfg
	}

	if fileExists(configPath) {
		return initConfig(configPath)
	}

	glog.Exitf("Incorrect file path [%s]", configPath)
	return &Config{}
}

// Init method reads config from yaml file
func initConfig(path string) *Config {
	cin := &Config{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Fatalf("Error while reading yaml file %s from path %s\n", err, yamlFile)
	}

	err = yaml.Unmarshal(yamlFile, &cin)
	if err != nil {
		glog.Fatal("Unable to decode config file", err)
	}

	cfg = cin
	return cfg
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
