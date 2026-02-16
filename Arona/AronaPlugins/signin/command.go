package signin

import (
	"Shittim/Arona/cmd"
	"Shittim/pkg/database"
	"Shittim/pkg/models"
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 签到模块
type SigninModule struct{}

// 返回模块名称
func (m *SigninModule) Name() string {
	return "signin"
}

// Enter 进入模块时执行的操作
func (m *SigninModule) Enter(ctx *zero.Ctx) {
	ctx.Send("已进入签到系统\n可用命令：\n- signin: 执行签到\n- exp: 查询经验\n- rank: 查看经验排行榜\n- 排名: 查看我的排名\n- exit: 退出签到系统")
}

// 退出模块时执行的操作
func (m *SigninModule) Exit(ctx *zero.Ctx) {
	ctx.Send("已退出签到系统")
}

// 处理模块内的命令
func (m *SigninModule) HandleCommand(cmd string, args []string, ctx *zero.Ctx) bool {
	switch cmd {
	case "signin", "签到":
		qq := ctx.Event.UserID
		nickname := ctx.Event.Sender.NickName

		reward, streak, err := DoSignin(qq, nickname)
		if err != nil {
			ctx.Send(message.Text("签到失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"🎉 签到成功！\n",
			fmt.Sprintf("获得奖励：%d 经验值\n", reward),
			fmt.Sprintf("连续签到：%d 天\n", streak),
			"努力成为什亭之箱的守护者吧！",
		))
		return true
	case "exp", "经验":
		qq := ctx.Event.UserID

		//查询用户信息
		var user models.User
		result := database.GetDB().Where("qq = ?", qq).First(&user)

		if result.Error != nil {
			ctx.Send(message.Text("查询失败：用户不存在，请先签到注册！"))
			return true
		}

		ctx.Send(message.Text(
			"📊 经验信息\n",
			fmt.Sprintf("当前等级：%d\n", user.Level),
			fmt.Sprintf("当前经验：%d\n", user.Exp),
			"继续签到获取更多经验值吧！",
		))
		return true
	case "rank", "排行", "排行榜":
		//获取经验排行榜
		users, err := GetExpRank(10)
		if err != nil {
			ctx.Send(message.Text("获取排行榜失败：", err.Error()))
			return true
		}

		//构建排行榜消息
		rankMsg := "🏆 经验排行榜\n"
		for i, user := range users {
			rankMsg += fmt.Sprintf("%d. %s - %d 经验值\n", i+1, user.Nickname, user.Exp)
		}

		ctx.Send(message.Text(rankMsg))
		return true
	case "我的排名", "排名":
		//获取用户排名
		rank, err := GetUserRank(ctx.Event.UserID)
		if err != nil {
			ctx.Send(message.Text("获取排名失败：用户不存在，请先签到注册！"))
			return true
		}

		//查询用户信息
		var user models.User
		database.GetDB().Where("qq = ?", ctx.Event.UserID).First(&user)

		ctx.Send(message.Text(
			"📈 我的排名\n",
			fmt.Sprintf("当前排名：第 %d 名\n", rank),
			fmt.Sprintf("当前经验：%d\n", user.Exp),
			"继续签到提升排名吧！",
		))
		return true
	default:
		return false
	}
}

// 注册签到模块
func RegisterModule(cmdSystem interface {
	RegisterModule(module cmd.Module)
}) {
	cmdSystem.RegisterModule(&SigninModule{})
}
