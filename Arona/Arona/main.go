package main

import (
	"Shittim/Arona/AronaPlugins/balogo"
	"Shittim/Arona/AronaPlugins/signin"
	"Shittim/Arona/cmd"
	"Shittim/config"
	"Shittim/pkg/database"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	// 加载配置文件
	config.InitConfig()

	// 初始化数据库连接
	database.InitDatabase()

	// 自动迁移数据库模型
	database.AutoMigrate()

	// 创建命令系统实例
	cmdSystem := cmd.NewCommandSystem()

	// 注册签到模块
	signin.RegisterModule(cmdSystem)

	// 注册蔚蓝档案Logo生成模块
	balogo.RegisterModule(cmdSystem)

	// 注册消息处理器
	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		// 获取用户输入
		input := ctx.Event.RawMessage

		// 处理用户输入
		cmdSystem.ProcessInput(input, ctx)
	})

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
