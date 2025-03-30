package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var Resume ResumeInformation

type ResumeInformation struct {
	Name string
	Sex string
	Age int
	Habits []interface{}
}

type ResumeSetting struct {
	TimeStamp string
	Address string
	ResumeInformation ResumeInformation
}

// 解析yaml
func parseYaml(v *viper.Viper) {
	var resumeConfig ResumeSetting
	err := v.Unmarshal(&resumeConfig)
	if err != nil {
		fmt.Println("parse yaml err: ", err)
	}
	fmt.Println("resume config:\n", resumeConfig)
}

// 获取sub-tree
func sub(key string, value interface{}) error {
	log.Printf("配置文件的前缀为： %v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

func init() {
	viper.AutomaticEnv()	// 通过环境变量修改任意配置
	initDefault()			// 初始化 viper 配置

	// 读取yaml
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err: %s\n", err)
	}

	// 反序列化struct
	if err := sub("ResumeInformation", &Resume); err != nil {
		log.Fatal("Fail to parse config", err)
	}
}

func initDefault() {
	viper.SetConfigName("resume-config")	// 读取的配置文件名
	viper.AddConfigPath("./config/")		// 配置文件路径
	viper.AddConfigPath("GOPATH/src/")		// 可添加多个可选路径
	viper.SetConfigType("yaml")				// 配置文件类型
}

func main() {
	fmt.Printf(
		"姓名：%s\n 爱好：%s\n 性别：%s\n 年龄：%d\n",
		Resume.Name,
		Resume.Habits,
		Resume.Sex,
		Resume.Age,
	)

	// 反序列化并输出ResumeSetting
	parseYaml(viper.GetViper())
}


