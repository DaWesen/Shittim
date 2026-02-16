package balogo

import (
	"Shittim/Arona/cmd"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 蔚蓝档案logo生成模块
type BaLogoModule struct{}

// 返回模块名称
func (m *BaLogoModule) Name() string {
	return "balogo"
}

// 进入模块时执行的操作
func (m *BaLogoModule) Enter(ctx *zero.Ctx) {
	ctx.Send("已进入蔚蓝档案Logo生成系统\n可用命令：\n- gen: 生成Logo\n- help: 查看帮助\n- exit: 退出系统")
}

// 退出模块时执行的操作
func (m *BaLogoModule) Exit(ctx *zero.Ctx) {
	ctx.Send("已退出蔚蓝档案Logo生成系统")
}

// 处理模块内的命令
func (m *BaLogoModule) HandleCommand(cmd string, args []string, ctx *zero.Ctx) bool {
	switch cmd {
	case "gen", "生成":
		if len(args) < 2 {
			ctx.Send(message.Text("请提供左右两部分文字哦~格式：gen 左文字 右文字"))
			return true
		}

		left := args[0]
		right := strings.Join(args[1:], " ")
		logoURL := GetBALOGO(left, right)

		ctx.Send(message.Message{
			message.Image(logoURL),
		})
		return true
	case "help", "帮助":
		ctx.Send("蔚蓝档案Logo生成系统帮助\n\n命令格式：\n- gen 左文字 右文字：生成Logo图片\n- exit：退出系统\n\n示例：\n- gen 什亭 之箱：生成左右文字分别为'什亭'和'之箱'的Logo")
		return true
	default:
		return false
	}
}

// 注册蔚蓝档案Logo生成模块
func RegisterModule(cmdSystem interface {
	RegisterModule(module cmd.Module)
}) {
	cmdSystem.RegisterModule(&BaLogoModule{})
}
