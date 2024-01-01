package conf

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
)

var iBookConfig *Config
var RELOAD = 1

type Config struct {
	ServerConf *Server `yaml:"server"`
	DataConf   *Data   `yaml:"data"`
	SecretConf *Secret `yaml:"secret"`
	Sms        *SMS    `yaml:"sms"`
}

type Server struct {
	Port   string `yaml:"port"`
	Domain string `yaml:"domain"`
}

type Data struct {
	MysqlConf *MySQL `yaml:"mysql"`
	RedisConf *Redis `yaml:"redis"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Addr string `yaml:"addr"`
}

type Secret struct {
	JwtConf *Jwt `yaml:"jwt"`
}

type Jwt struct {
	Key              string `yaml:"key"`
	LifeDurationTime int64  `yaml:"life_duration_time"`
}

type SMS struct {
	Login *SMSLogin `yaml:"login"`
}

type SMSLogin struct {
	Aliyun     *AliyunSMS     `yaml:"aliyun"`
	Tencentyun *TencentyunSMS `yaml:"tencentyun"`
}

type AliyunSMS struct {
	SignName        string `yaml:"sign_name"`
	TplCode         string `yaml:"tpl_code"`
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

type TencentyunSMS struct {
	SignName  string `yaml:"sign_name"`
	TplCode   string `yaml:"tpl_code"`
	AppId     string `yaml:"app_id"`
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

func GetConf(flags ...int) *Config {
	if iBookConfig == nil || (len(flags) != 0 && flags[0] == RELOAD) {
		dir, err := filepath.Abs(filepath.Dir("./"))
		if err != nil {
			log.Fatal(err)
		}
		configPath := filepath.Join(dir, "configs", "config.yaml")
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
