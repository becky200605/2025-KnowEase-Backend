package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // 允许使用环境变量覆盖配置

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: Unable to load .env file")
	}
}

func ProvideDBConnectionString() string {
	dsn := viper.GetString("DSN") // 读取 DSN 配置
	if dsn == "" {
		panic("未设置数据库连接，请检查环境变量")
	}
	return dsn
}

// @title KnowEase API
// @version 1.0
// @description 小知，你的校园助手
// @host localhost:8080
// @BasePath /api
func main() {
	app := InitializeApp()
	app.Run()
}
