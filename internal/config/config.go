package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"soul/global"
	"strings"
)

func LoadConfig() *viper.Viper {
	// 初始化viper
	v := viper.New()

	// 命令行参数获取
	pflag.StringP("env", "e", "dev", `运行环境, 可选项目: dev or test prod`)
	pflag.StringP("config", "c", "app-dev.yaml", `配置文件路径`)
	pflag.Lookup("config").DefValue = "[./app-dev.yaml, ./config/app-dev.yaml]"
	pflag.BoolP("migrate", "m", false, `迁移数据库`)
	pflag.StringP("kubeconfig", "k", "", `kubeconfig path, default in k8s cluster`)

	pflag.Parse()

	// 命令行参数绑定
	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Println("[Init] 命令行参数绑定失败")
		panic(err.Error())
	}

	// 环境变量参数绑定
	err := v.BindEnv("env", "RUN_ENV")
	if err != nil {
		fmt.Println("[Init] 环境变量参数绑定失败")
		panic(err.Error())
	}

	// 获取当前运行环境
	env := strings.TrimSpace(v.GetString("env"))

	// 设置默认参数
	if env == "test" {
		setTestDefaultParams(v)
	} else if env == "prod" {
		setProdDefaultParams(v)
	} else {
		setDevDefaultParams(v)
	}

	// 不同环境读取不同配置文件
	switch env {
	// 只能是下面三种环境, 如果为其他的就设置为dev
	case "dev", "test", "prod":
		filePath := v.GetString("config")
		if filePath != "" && v.IsSet("config") {
			// 如果指定了配置文件路径，就读取指定的配置文件
			v.SetConfigFile(filePath)
		} else {
			// 没有指定配置文件，设置对应环境的默认配置文件路径
			v.SetConfigName(fmt.Sprintf("app-%s", env))
		}
		fmt.Printf("[Init] 当前运行环境: %s\n", env)
	default:
		v.Set("env", "dev")
		v.SetConfigName("app-dev.yaml") // 默认dev环境
		fmt.Printf("[Init] 未知环境,使用默认运行环境, 默认: %s\n", "dev")
	}

	// 设置viper, 加载配置文件
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("[Init] 配置文件读取失败: %s, \n使用默认配置\n", err.Error())
	} else {
		fmt.Printf("[Init] 使用配置文件: %s\n", v.ConfigFileUsed())
	}

	// 设置环境变量前缀为appName
	//v.SetEnvPrefix(v.GetString("appName"))

	// 动态加载配置
	//v.WatchConfig()
	//
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	if err = v.UnmarshalExact(&global.Config); err != nil {
	//		fmt.Println("动态加载配置失败" + err.Error())
	//	}
	//})

	err = v.Unmarshal(&global.Config)
	if err != nil {
		panic("加载配置失败" + err.Error())
	}

	return v
}
