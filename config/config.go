package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func InitConfig() {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./etc") // 配置文件路径设置为当前路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	viper.AutomaticEnv()                                   // 自动绑定环境变量
	viper.BindEnv("config_name", "VIPER_CONFIG_NAME")      // 绑定环境变量到变量APP_NAME,明确指定环境变量名字 VIPER_CONFIG_NAME
	viper.BindEnv("config_version")                        // 绑定环境变量到变量APP_VERSION，自动组装环境变量名字 VIPER_CONFIG_VERSION
	viper.BindEnv("app.port", "VIPER_APP_PORT")            // 绑定环境变量到变量APP_PORT,明确指定环境变量名字 VIPER_APP_PORT
	viper.SetEnvPrefix("VIPER")                            // 设置环境变量前缀为VIPER
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 环境变量名字使用下划线分隔
	viper.BindEnv("app.name")                              // 绑定环境变量到变量APP_NAME,自动组装环境变量名字 VIPER_APP_NAME,这里依赖了SetEnvKeyReplacer参数把VIPER_APP.NAME转换成VIPER_APP_NAME

	// 监听配置文件变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	// 读取etcd的远程配置文件
	viper.AddRemoteProvider("etcd3", "http://10.130.34.85:2379", "/viper/config")
	viper.SetConfigType("json")
	err := viper.ReadRemoteConfig()
	if err != nil {
		log.Printf("读取远程配置文件失败: %v", err)
	} else {
		fmt.Println("读取远程配置文件成功：" + viper.GetString("etcd_version"))
	}

}
