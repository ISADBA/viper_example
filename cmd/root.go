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
		// 	log.Printf("写入配置文件失败: %v", err)
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
