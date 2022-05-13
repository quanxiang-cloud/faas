package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/cabin/tailormade/db/redis"
)

// DefaultPath default file path
var DefaultPath = "./configs/config.yml"

// Config config info
type Config struct {
	Port        string         `yaml:"port"`
	Model       string         `yaml:"model"`
	InternalNet client.Config  `yaml:"internalNet"`
	Log         logger.Config  `yaml:"log"`
	Mysql       mysql.Config   `yaml:"mysql"`
	Docker      Docker         `yaml:"docker"`
	Redis       redis.Config   `yaml:"redis"`
	Elastic     elastic.Config `yaml:"elastic"`
	Graph       struct {
		Runs  []string   `yaml:"runs,omitempty" json:"runs"`
		Steps [][]string `yaml:"steps,omitempty" json:"steps"`
	} `yaml:"graph,omitempty"`

	BuildImages map[string]string `yaml:"build-images"`
	Templates   []*Template       `yaml:"templates"`
}

// Docker docker
type Docker struct {
	NameSpace string `yaml:"namespace"`
}

type Template struct {
	FullName string `yaml:"full_name"`
	Content  string `yaml:"content"`
}

// NewConfig new
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	config := &Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
