package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var ServerConfigs *serviceConfig

func InitConfigs() error {
	var err error
	bs, err := os.ReadFile("./etc/config.yaml")
	if err != nil {
		return err
	}
	ServerConfigs, err = parse(bs)
	return nil
}

type postgresqlConfig struct {
	HostIP   string `yaml:"host_ip"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	DbName   string `yaml:"db_name"`
	Password string `yaml:"password"`
}

type webConfig struct {
	JWTKey string `yaml:"jwt_key"`
}

type serviceConfig struct {
	Postgresql postgresqlConfig `yaml:"postgresql"`
	Redis      redisConfig      `yaml:"redis"`
	Server     webConfig        `yaml:"web"`
}

type redisConfig struct {
	HostIP   string `yaml:"host_ip"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

func parse(data []byte) (*serviceConfig, error) {
	c := &serviceConfig{}
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
