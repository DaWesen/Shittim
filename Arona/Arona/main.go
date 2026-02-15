package main

import (
	"Shittim/config"
	"Shittim/pkg/database"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	// 导入插件
	_ "Shittim/Arona/AronaPlugins/signin"
)

func main() {
	// 加载配置文件
	config.InitConfig()

	// 初始化数据库连接
	database.InitDatabase()

	// 自动迁移数据库模型
	database.AutoMigrate()

	// 运行机器人
	zero.RunAndBlock(&zero.Config{
		NickName:      config.AppConfig.Arona.NickName,
		CommandPrefix: config.AppConfig.Arona.CommandPrefix,
		SuperUsers:    config.AppConfig.Arona.SuperUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketServer(
				config.AppConfig.Arona.WebSocket.Port,
				config.AppConfig.Arona.WebSocket.URL,
				"",
			),
		},
	}, nil)
}
