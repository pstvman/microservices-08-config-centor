package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"

	"net/http"
)

var (
	Resume resumeConfig
)

type resumeConfig struct {
	Name string
	Age int
	Sex string
}

const (
	kAppName 		= "APP_NAME"
	kConfigServer 	= "CONFIG_SERVER"
	kConfigLabel 	= "CONFIG_LABLE"
	kConfigProfile 	= "CONFIG_PROFILE"
	kConfigType 	= "CONFIG_TYPE"
	kAmqpURI       = "AmqpURI"
)

func init() {
	viper.AutomaticEnv()	// 通过环境变量修改任意配置
	initDefault()			// 初始化 viper 配置

	// 开启监听
	go StartListener(viper.GetString(kAppName), viper.GetString(kAmqpURI), "springCloudBus")

	// 读取yaml
	if err := loadRemoteConfig(); err != nil {
		log.Fatal("Fail to load config: ", err)
	}

	// 反序列化struct
	if err := sub("resume", &Resume); err != nil {
		log.Fatal("Fail to parse config: ", err)
	}
}

func initDefault() {
	viper.SetDefault(kAppName, "client-demo")					// 应用名
	viper.SetDefault(kConfigServer, "http://localhost:8888")	// 服务地址
	viper.SetDefault(kConfigLabel, "master")						// 分支
	viper.SetDefault(kConfigProfile, "dev")						// 环境
	viper.SetDefault(kConfigType, "yaml")						// 配置文件格式

	viper.SetDefault(kAmqpURI, "amqp://guest:guest@localhost:5672")	// RabbitMQ 地址
}

func loadRemoteConfig() (err error) {
	confAddr := fmt.Sprintf("%v/%v/%v-%v.yml", 
		viper.Get(kConfigServer), 
		viper.Get(kConfigLabel),
		viper.Get(kAppName),
		viper.Get(kConfigProfile))
	
	// 发送http请求
	resp, err := http.Get(confAddr)
	if err != nil {
		fmt.Println("Send request err: ", err)
		return
	}
	defer resp.Body.Close()

	viper.SetConfigType(viper.GetString(kConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		fmt.Println("Load config err: ", err)
		return
	}
	log.Println("Load config from: ", confAddr)
	return
}

func sub(key string, value interface{}) error {
	sub := viper.Sub(key)
	if sub == nil {
		return fmt.Errorf("configuration for %s not found", key)
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

