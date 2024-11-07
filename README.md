介绍
通常我们在编写程序的时候，一些依赖的元数据信息，我们可以通过命令参数或者配置文件进行传递，如果元数据信息比较多，使用配置文件肯定是更加合理的，更合理的场景是使用命令行参数+配置文件的组合方式来提供这些元数据。

如果我们自己实现一个配置文件管理也是比较简单的，只需要声明一个配置文件对应的结构体，然后针对配置文件的格式类型解析成结构体即可。

这也是viper最基本的功能，但是viper在这个基础上还实现了很多其他的能力，比如：
1. 支持给配置项设置默认值
2. 支持多种格式的配置文件，比如json，toml，yaml，hcl，envfile，java properties配置文件
3. 支持实时监视配置文件变化，并且动态读取新的配置
4. 支持从环境变量读取参数值
5. 支持从远程系统获取配置值，并实时感知变化，比如consul和etcd
6. 从命令行参数读取(配合cobra)

viper参数不同设置方法的优先级，越前面优先级越高
1. 代码中显示调用set方法
2. 通过命令行的flag参数设置
3. 从env环境变量获取
4. 配置文件本身的内容
5. k/v键值存储获取的内容
6. 参数的默认值
7. 支持将内存中的配置回写配置文件

注意：viper配置项的key是不区分大小的。


相关术语


最佳实践
1. 安装viper
go get github.com/spf13/viper
mkdir viper_example
cd viper_example
go mod init viper_example
2. 一个简单的使用viper管理配置文件的用例
  1. 目录结构设计
├── config
│   └── config.go
├── etc
│   └── config.yaml
├── go.mod
├── go.sum
└── main.go
  2. 要解析的配置文件
# etc/config.yaml
config_name: config.yaml
config_version: 1.0.0

# list结构 + 嵌套结构
repository:
  - name: viper_example
    dialector: mysql
    url: mysql://root:root@localhost:3306/viper_example
  - name: viper_example_2
    url: sqlite:///viper_example_2.db
    dialector: sqlite

# map结构,key内容不确定
databases:
  mysql: true
  sqlite: true
  redis: true

log:
  level: info
  format: text
  output: stdout

app:
  name: viper_example
  port: 8080
  3. 相关代码
    1. 提供了以下的get方法来获取数据
    2. viper.Get()
    3. viper.GetStringMap()
    4. viper.GetString()
    5. viper.GetInt()
    6. viper.GetInt32()
    7. viper.GetInt64()
    8. viper.GetIntSlice()
    9. viper.GetStringSlice()等方法，你大概根据方法名字可以猜出含义和其他方法名字了。
#main.go
package main

import (
    "fmt"
    "viper_example/config"

    "github.com/spf13/viper"
)

func main() {
    // 使用viper初始化配置
    config.InitConfig()

    // 读取基本配置信息
    fmt.Println("Config Name:", viper.GetString("config_name"))
    fmt.Println("Config Version:", viper.GetString("config_version"))

    // 读取 repository 列表
    repos := viper.Get("repository").([]interface{})
    fmt.Println("Repositories:")
    for _, repo := range repos {
        repoMap := repo.(map[string]interface{})
        fmt.Printf("  - Name: %s\n", repoMap["name"])
        fmt.Printf("    Dialector: %s\n", repoMap["dialector"])
        fmt.Printf("    URL: %s\n", repoMap["url"])
    }

    // 读取 databases Map 结构
    fmt.Println("Databases:")
    databases := viper.GetStringMap("databases")
    for db, enabled := range databases {
        fmt.Printf("  %s: %v\n", db, enabled)
    }

    // 读取 log 配置
    fmt.Println("Log Configuration:")
    fmt.Printf("  Level: %s\n", viper.GetString("log.level"))
    fmt.Printf("  Format: %s\n", viper.GetString("log.format"))
    fmt.Printf("  Output: %s\n", viper.GetString("log.output"))

    // 读取 app 配置
    fmt.Println("App Configuration:")
    fmt.Printf("  Name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  Port: %d\n", viper.GetInt("app.port"))
}


# config/config.go
package config

import (
    "log"

    "github.com/spf13/viper"
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

}
  4. 重要：
    1. 从这个示例可以看出，viper解析后的配置文件，是通过key名字去取的，并不是直接解析到了一个struct里面。
    2. viper解析后的配置，是在一个全局的单例对象viper中。
3. 在viper管理配置文件的同时，引入环境变量
  1. 引入VIPER_开始的环境变量
  2. viper.BindEnv() 有两种用法，
    1. 第一种指制定配置文件中的参数名字，那么viper会根据viper.SetEnvPrefix和viper.SetEnvKeyReplacer的设置值自动绑定对应的环境变量名字
    2. 第二种就是在第二个参数中明确环境变量的名字，如果环境变量不多，建议使用第二种
  3. viper.AutomaticEnv() 会自动加载与配置文件中key相关的变量名字，会根据viper.SetEnvPrefix和viper.SetEnvKeyReplacer的设置值组装对应的环境变量名字
  4. 建议： 如果依赖环境变量不多的话，最好使用viper.BindEnv()指定两个参数，避免不确定性。
# 调整 config.go
package config

import (
    "log"
    "strings"

    "github.com/spf13/viper"
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
}

# 调整main.go
package main

import (
    "fmt"
    "viper_example/config"

    "github.com/spf13/viper"
)

func main() {
    config.InitConfig()

    // 读取基本配置信息
    fmt.Println("Config Name:", viper.GetString("config_name"))
    fmt.Println("Config Version:", viper.GetString("config_version"))

    // 读取 repository 列表
    repos := viper.Get("repository").([]interface{})
    fmt.Println("Repositories:")
    for _, repo := range repos {
        repoMap := repo.(map[string]interface{})
        fmt.Printf("  - Name: %s\n", repoMap["name"])
        fmt.Printf("    Dialector: %s\n", repoMap["dialector"])
        fmt.Printf("    URL: %s\n", repoMap["url"])
    }

    // 读取 databases Map 结构
    fmt.Println("Databases:")
    databases := viper.GetStringMap("databases")
    for db, enabled := range databases {
        fmt.Printf("  %s: %v\n", db, enabled)
    }

    // 读取 log 配置
    fmt.Println("Log Configuration:")
    fmt.Printf("  Level: %s\n", viper.GetString("log.level"))
    fmt.Printf("  Format: %s\n", viper.GetString("log.format"))
    fmt.Printf("  Output: %s\n", viper.GetString("log.output"))

    // 读取 app 配置
    fmt.Println("App Configuration:")
    fmt.Printf("  Name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  Port: %d\n", viper.GetInt("app.port"))

    // 读取环境变量
    fmt.Println("Environment Variables:")
    // 使用环境变量覆盖配置文件内容
    // 使用viper.BindEnv()方法绑定环境变量到Viper实例
    fmt.Printf("  config_name: %s\n", viper.GetString("config_name"))
    fmt.Printf("  config_version: %s\n", viper.GetString("config_version"))
    fmt.Printf("  app.name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  app.port: %s\n", viper.GetString("app.port"))
    // 使用viper.AutomaticEnv()方法自动绑定环境变量到Viper实例
    fmt.Printf("  DATABASE_URL: %s\n", viper.GetString("VIPER_DATABASE_URL"))
}


4. 在viper管理配置文件的同时，引入cobra参数配置
  1. 使用命令行参数优先环境变量，使环境变量优先配置文件内容
  2. 这里引入要引入cobra，所以目录机构做了一些调整
├── cmd
│   └── root.go
├── config
│   └── config.go
├── etc
│   └── config.yaml
├── go.mod
├── go.sum
└── main.go
  3. 
# main.go

package main

import (
    "viper_example/cmd"
)

func main() {
    cmd.Execute()
}

# cmd.root.go
package cmd

import (
    "fmt"
    "viper_example/config"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "viper_example",
    Short: "Example of using Viper for configuration management",
    Long:  `Example of using Viper for configuration management.`,
    Run:   RunMain,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}

func init() {
    // 初始化配置文件
    config.InitConfig()
    // 使用命令行的config_version参数代替环境变量和配置文件的内容
    rootCmd.PersistentFlags().String("config_version", "default_flag_config_version", "config_version")
    rootCmd.PersistentFlags().String("app_name", "default_flag_app_name", "app_name")
    viper.BindPFlag("config_version", rootCmd.PersistentFlags().Lookup("config_version"))
    viper.BindPFlag("app.name", rootCmd.PersistentFlags().Lookup("app_name"))
}

func RunMain(cmd *cobra.Command, args []string) {

    // 读取基本配置信息
    fmt.Println("Config Name:", viper.GetString("config_name"))
    fmt.Println("Config Version:", viper.GetString("config_version"))

    // 读取 repository 列表
    repos := viper.Get("repository").([]interface{})
    fmt.Println("Repositories:")
    for _, repo := range repos {
        repoMap := repo.(map[string]interface{})
        fmt.Printf("  - Name: %s\n", repoMap["name"])
        fmt.Printf("    Dialector: %s\n", repoMap["dialector"])
        fmt.Printf("    URL: %s\n", repoMap["url"])
    }

    // 读取 databases Map 结构
    fmt.Println("Databases:")
    databases := viper.GetStringMap("databases")
    for db, enabled := range databases {
        fmt.Printf("  %s: %v\n", db, enabled)
    }

    // 读取 log 配置
    fmt.Println("Log Configuration:")
    fmt.Printf("  Level: %s\n", viper.GetString("log.level"))
    fmt.Printf("  Format: %s\n", viper.GetString("log.format"))
    fmt.Printf("  Output: %s\n", viper.GetString("log.output"))

    // 读取 app 配置
    fmt.Println("App Configuration:")
    fmt.Printf("  Name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  Port: %d\n", viper.GetInt("app.port"))

    // 读取环境变量
    fmt.Println("Environment Variables:")
    // 使用环境变量覆盖配置文件内容
    // 使用viper.BindEnv()方法绑定环境变量到Viper实例
    fmt.Printf("  config_name: %s\n", viper.GetString("config_name"))
    fmt.Printf("  config_version: %s\n", viper.GetString("config_version"))
    fmt.Printf("  app.name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  app.port: %s\n", viper.GetString("app.port"))
    // 使用viper.AutomaticEnv()方法自动绑定环境变量到Viper实例
    fmt.Printf("  DATABASE_URL: %s\n", viper.GetString("VIPER_DATABASE_URL"))
}




# config.go没有变化


# 使用
go run main.go --config_version flag_config_version --app_name flag_app_name
5. 配置文件变更，动态处理变更内容
  1. 需要做两个改造，第一个是动态监听配置文件
  2. 第二个是把我们项目改为一个deamon类型的程序
  3. config.go 增加试试监听配置文件变化
# config.go 支持监听配置文件
package config

import (
    "fmt"
    "log"
    "strings"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
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
}


  4. root.go 增加deamon进程，实时读取log.level配置
# root.go
package cmd

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
    "viper_example/config"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "viper_example",
    Short: "Example of using Viper for configuration management",
    Long:  `Example of using Viper for configuration management.`,
    Run:   RunMain,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}

func init() {
    // 初始化配置文件
    config.InitConfig()
    // 使用命令行的config_version参数代替环境变量和配置文件的内容
    rootCmd.PersistentFlags().String("config_version", "default_flag_config_version", "config_version")
    rootCmd.PersistentFlags().String("app_name", "default_flag_app_name", "app_name")
    viper.BindPFlag("config_version", rootCmd.PersistentFlags().Lookup("config_version"))
    viper.BindPFlag("app.name", rootCmd.PersistentFlags().Lookup("app_name"))
}

func RunMain(cmd *cobra.Command, args []string) {

    // 读取基本配置信息
    fmt.Println("Config Name:", viper.GetString("config_name"))
    fmt.Println("Config Version:", viper.GetString("config_version"))

    // 读取 repository 列表
    repos := viper.Get("repository").([]interface{})
    fmt.Println("Repositories:")
    for _, repo := range repos {
        repoMap := repo.(map[string]interface{})
        fmt.Printf("  - Name: %s\n", repoMap["name"])
        fmt.Printf("    Dialector: %s\n", repoMap["dialector"])
        fmt.Printf("    URL: %s\n", repoMap["url"])
    }

    // 读取 databases Map 结构
    fmt.Println("Databases:")
    databases := viper.GetStringMap("databases")
    for db, enabled := range databases {
        fmt.Printf("  %s: %v\n", db, enabled)
    }

    // 读取 log 配置
    fmt.Println("Log Configuration:")
    fmt.Printf("  Level: %s\n", viper.GetString("log.level"))
    fmt.Printf("  Format: %s\n", viper.GetString("log.format"))
    fmt.Printf("  Output: %s\n", viper.GetString("log.output"))

    // 读取 app 配置
    fmt.Println("App Configuration:")
    fmt.Printf("  Name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  Port: %d\n", viper.GetInt("app.port"))

    // 读取环境变量
    fmt.Println("Environment Variables:")
    // 使用环境变量覆盖配置文件内容
    // 使用viper.BindEnv()方法绑定环境变量到Viper实例
    fmt.Printf("  config_name: %s\n", viper.GetString("config_name"))
    fmt.Printf("  config_version: %s\n", viper.GetString("config_version"))
    fmt.Printf("  app.name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  app.port: %s\n", viper.GetString("app.port"))
    // 使用viper.AutomaticEnv()方法自动绑定环境变量到Viper实例
    fmt.Printf("  DATABASE_URL: %s\n", viper.GetString("VIPER_DATABASE_URL"))

    // 启动守护进程
    StartServer()
}

// 任务函数，守护进程将循环执行的内容
func task() {
    for {
        fmt.Println("Daemon task running...")                   // 在后台运行的任务
        fmt.Println("log.level:", viper.GetString("log.level")) // 打印配置信息
        time.Sleep(5 * time.Second)                             // 模拟任务执行间隔
    }
}

func StartServer() {
    // 创建信号通道，用于接收系统信号
    sigChannel := make(chan os.Signal, 1)
    signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM) // 捕获系统信号

    // 启动任务
    go task()

    fmt.Println("Daemon started. Press Ctrl+C to stop.")

    // 等待系统信号
    sig := <-sigChannel
    fmt.Printf("Received signal: %s. Exiting...\n", sig)
}


6. 从etcd获取配置文件
  1. 安装一个单机版的etcd用于测试
# 安装etcd
docker run --name etcd -d -p 2379:2379 quay.io/coreos/etcd:latest /usr/local/bin/etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://localhost:2379
# 进入容器
docker exec -it etcd sh
# 设置api为v3
export ETCDCTL_API=3
# 写入key
etcdctl put /viper/config '{"app_name":"etcd_name","etcd_version":"version3.30"}'
# 读取key
etcdctl get /viper/config
 {"app_name":"etcd_name","etcd_version":"version3.30"}
  2. 修改代码，config.go
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

7. etcd的配置变更，动态处理变更内容
  1. 注释掉viper.WriteConfig()相关代码，因为配置文件的优先级必remote配置高，回写配置文件后就不会再使用remote的配置内容
  2. 使用viper.WatchRemoteConfig()来读取远程配置，但是注意：这个不是运行一次就一直监听，而是要自己开个协程，定时触发这个请求
  3. 修改root.go，重点在task()内
package cmd

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
    "viper_example/config"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "viper_example",
    Short: "Example of using Viper for configuration management",
    Long:  `Example of using Viper for configuration management.`,
    Run:   RunMain,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}

func init() {
    // 初始化配置文件
    config.InitConfig()
    // 使用命令行的config_version参数代替环境变量和配置文件的内容
    rootCmd.PersistentFlags().String("config_version", "default_flag_config_version", "config_version")
    rootCmd.PersistentFlags().String("app_name", "default_flag_app_name", "app_name")
    viper.BindPFlag("config_version", rootCmd.PersistentFlags().Lookup("config_version"))
    viper.BindPFlag("app.name", rootCmd.PersistentFlags().Lookup("app_name"))
}

func RunMain(cmd *cobra.Command, args []string) {

    // 读取基本配置信息
    fmt.Println("Config Name:", viper.GetString("config_name"))
    fmt.Println("Config Version:", viper.GetString("config_version"))

    // 读取 repository 列表
    repos := viper.Get("repository").([]interface{})
    fmt.Println("Repositories:")
    for _, repo := range repos {
        repoMap := repo.(map[string]interface{})
        fmt.Printf("  - Name: %s\n", repoMap["name"])
        fmt.Printf("    Dialector: %s\n", repoMap["dialector"])
        fmt.Printf("    URL: %s\n", repoMap["url"])
    }

    // 读取 databases Map 结构
    fmt.Println("Databases:")
    databases := viper.GetStringMap("databases")
    for db, enabled := range databases {
        fmt.Printf("  %s: %v\n", db, enabled)
    }

    // 读取 log 配置
    fmt.Println("Log Configuration:")
    fmt.Printf("  Level: %s\n", viper.GetString("log.level"))
    fmt.Printf("  Format: %s\n", viper.GetString("log.format"))
    fmt.Printf("  Output: %s\n", viper.GetString("log.output"))

    // 读取 app 配置
    fmt.Println("App Configuration:")
    fmt.Printf("  Name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  Port: %d\n", viper.GetInt("app.port"))

    // 读取环境变量
    fmt.Println("Environment Variables:")
    // 使用环境变量覆盖配置文件内容
    // 使用viper.BindEnv()方法绑定环境变量到Viper实例
    fmt.Printf("  config_name: %s\n", viper.GetString("config_name"))
    fmt.Printf("  config_version: %s\n", viper.GetString("config_version"))
    fmt.Printf("  app.name: %s\n", viper.GetString("app.name"))
    fmt.Printf("  app.port: %s\n", viper.GetString("app.port"))
    // 使用viper.AutomaticEnv()方法自动绑定环境变量到Viper实例
    fmt.Printf("  DATABASE_URL: %s\n", viper.GetString("VIPER_DATABASE_URL"))

    // 启动守护进程
    StartServer()
}

// 任务函数，守护进程将循环执行的内容
func task() {
    for {
        fmt.Println("Daemon task running...")                   // 在后台运行的任务
        fmt.Println("log.level:", viper.GetString("log.level")) // 打印配置信息
        time.Sleep(5 * time.Second)                             // 模拟任务执行间隔

        // 修改配置文件，然后回写配置文件
        // viper.Set("app.name", "viper_example"+time.Now().String())
        // if err := viper.WriteConfig(); err != nil {
        //  log.Printf("写入配置文件失败: %v", err)
        // }

        // 打印远程etcd的配置文件变化
        // 监听etcd的远程配置文件变化
        err := viper.WatchRemoteConfig()
        if err != nil {
            fmt.Println("unable to read remote config:", err)
        }
        fmt.Println("Listening to remote etcd config changes...")
        fmt.Println("etcd_version:", viper.GetString("etcd_version"))

    }
}

func StartServer() {
    // 创建信号通道，用于接收系统信号
    sigChannel := make(chan os.Signal, 1)
    signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM) // 捕获系统信号

    // 启动任务
    go task()

    fmt.Println("Daemon started. Press Ctrl+C to stop.")

    // 等待系统信号
    sig := <-sigChannel
    fmt.Printf("Received signal: %s. Exiting...\n", sig)
}


8. 修改内存中的viper配置，回写配置文件
  1. 注意：回写配置文件不仅仅是自己使用viper.Set()设置的会回写，使用viper.BindEnv()，viper.BindPFlag()绑定的变量也会被回写，且配置文件顺序会有所变化。
  2. 调整task()方法
func task() {
    for {
        fmt.Println("Daemon task running...")                   // 在后台运行的任务
        fmt.Println("log.level:", viper.GetString("log.level")) // 打印配置信息
        time.Sleep(5 * time.Second)                             // 模拟任务执行间隔

        // 修改配置文件，然后回写配置文件
        viper.Set("app.name", "viper_example"+time.Now().String())
        if err := viper.WriteConfig(); err != nil {
            log.Printf("写入配置文件失败: %v", err)
        }
    }
}

其他
1. viper默认的读取配置文件内容的方式是通过viper.GetString("keyString")等方法，而不是解析成结构体，这个和我们自己实现一个通过yaml解析配置文件的方法套路完全不一样，viper的方式很像SpringBoot自动加载application.yaml的方式，然后在项目中使用@value或者@ConfigurationProperties(prefix = "repostiory")。
  1. 优点：
    灵活性高：可以动态读取任何配置值，不受结构体字段限制，尤其在配置项不确定、内容较为多变或频繁调整时非常灵活。
    减少代码耦合：不需要定义结构体，可以直接使用 viper.Get() 获取数据，代码较为简洁。
    支持嵌套数据的动态解析：对于复杂的嵌套结构，直接通过键路径（如 "app.port"）来获取值，无需额外定义嵌套结构体。
  2. 缺点：
    缺乏编译期检查：键名是字符串，拼写错误不会在编译期检查出来，容易引发运行时错误，尤其是在键名较多且复杂时。
    数据类型安全性较差：viper.Get() 返回的是接口类型，无法确定返回值的类型，可能导致类型断言错误，给开发带来不便。例如，viper.Get("app.port") 可能需要将接口类型转换为 int，增加了出错的可能。
    代码可读性和维护性差：代码中充斥着键字符串，长时间维护时不利于理解。对于新开发者或长期维护项目，这种方式的可读性较低。
2. 相关代码：https://github.com/ISADBA/viper_example
