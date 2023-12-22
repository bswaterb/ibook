package conf

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
)

var iBookConfig *Config

type Config struct {
	ServerConf Server `yaml:"server"`
	DataConf   Data   `yaml:"data"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Data struct {
	MysqlConf MySQL `yaml:"mysql"`
	RedisConf Redis `yaml:"redis"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Addr string `yaml:"addr"`
}

func GetConf() *Config {
	if iBookConfig == nil || iBookConfig.ServerConf.Port == "" {
		dir, err := filepath.Abs(filepath.Dir("./"))
		if err != nil {
			log.Fatal(err)
		}
		configPath := filepath.Join(dir, "backend", "configs", "config.yaml")
		log.Printf("Loading config from: %s", configPath)
		// confFile := "/Users/bswaterb/Coding/go/ibook/backend/configs/config.yaml"
		yamlFile, err := os.Open(configPath)
		if err != nil {
			panic("配置文件未知错误: " + err.Error())
		}
		defer yamlFile.Close()
		yamlData, err := io.ReadAll(yamlFile)
		config := &Config{}
		err = yaml.Unmarshal(yamlData, config)
		iBookConfig = config
		if err != nil {
			panic("配置文件解析错误" + err.Error())
		}
	}
	return iBookConfig
}
