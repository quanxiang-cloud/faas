package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
)

// DefaultPath default file path
var DefaultPath = "./configs/config.yml"

// Config config info
type Config struct {
	Port        string        `yaml:"port"`
	Model       string        `yaml:"model"`
	InternalNet client.Config `yaml:"internalNet"`
	Log         logger.Config `yaml:"log"`
	Mysql       mysql.Config  `yaml:"mysql"`
	K8s         K8s           `yaml:"k8s"`
	Docker      Docker        `yaml:"docker"`
}

// K8s k8s
type K8s struct {
	NameSpace string `yaml:"namespace"`
}

// Docker docker
type Docker struct {
	NameSpace string `yaml:"namespace"`
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
