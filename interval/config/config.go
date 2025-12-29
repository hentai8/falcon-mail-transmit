package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

const (
	defaultServerPort         = "8080"
	defaultJWTExpirationHours = 72
)

type Config struct {
	ServerPort string `json:"server_port"`
	LogDir     string `json:"log_dir"`
	Lmstfy     struct {
		Host      string `json:"host"`
		Port      int    `json:"port"`
		Namespace string `json:"namespace"`
		Token     string `json:"token"`
	} `json:"lmstfy"`
	Redis     RedisStruct  `json:"redis" validate:"required"`
	Feishu    FeishuStruct `json:"feishu" validate:"required"`
	MailTypes []string     `json:"mail_types"`
}

type RedisStruct struct {
	Address string `json:"address"`
	DB      int    `json:"db"`
}

type FeishuStruct struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func (c Config) Validate() error {
	validate := validator.New()
	return validate.Struct(&c)
}

func Load() (*Config, error) {
	configFileName, _ := filepath.Abs("../config/config.json")
	c := Config{
		ServerPort: defaultServerPort,
		//JWTExpiration: defaultJWTExpirationHours,
	}

	bytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
