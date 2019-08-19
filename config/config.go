package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type rabbitConfiguration struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	Pwd string `yaml:"pwd"`
	Queues []*queueConfiguration `yaml:"queues"`
}
type queueConfiguration struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	Qos int `yaml:"qos"`
}
type cacheConfiguration struct {
	Enable string `yaml:"enable"`
	Redis redisConfiguration `yaml:"redis"`
}
type redisConfiguration struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pwd string `yaml:"pwd"`
}
type AppConfiguration struct {
	CacheType *cacheConfiguration `yaml:"cache"`
	RabbitConfig *rabbitConfiguration `yaml:"rabbit"`
}
func(a *AppConfiguration) setConfigs() {
	yamlFile, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		log.Println(err)
		panic("read app.yaml file failed")
	}
	err = yaml.Unmarshal(yamlFile, a)
	if err != nil {
		log.Println(err)
		panic("unmarshal file app.yaml failed")
	}
}
var appConfig *AppConfiguration
var once sync.Once
func GetSystemConfig() *AppConfiguration{
	once.Do(func(){
		appConfig = &AppConfiguration{}
		appConfig.setConfigs()
	})
	return appConfig
}



