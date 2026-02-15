package signin

import (
	"Shittim/pkg/database"
	"Shittim/pkg/models"
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	//签到命令
	zero.OnRegex(`^签到$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		qq := ctx.Event.UserID
		nickname := ctx.Event.Sender.NickName

		reward, streak, err := DoSignin(qq, nickname)
		if err != nil {
			ctx.Send(message.Text("签到失败：", err.Error()))
			return
		}

		ctx.Send(message.Text(
			"🎉 签到成功！\n",
			fmt.Sprintf("获得奖励：%d 经验值\n", reward),
			fmt.Sprintf("连续签到：%d 天\n", streak),
			"努力成为什亭之箱的守护者吧！",
		))
	})

	//查询经验命令
	zero.OnRegex(`^经验$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		qq := ctx.Event.UserID

		//查询用户信息
		var user models.User
		result := database.GetDB().Where("qq = ?", qq).First(&user)

		if result.Error != nil {
			ctx.Send(message.Text("查询失败：用户不存在，请先签到注册！"))
			return
		}

		ctx.Send(message.Text(
			"📊 经验信息\n",
			fmt.Sprintf("当前等级：%d\n", user.Level),
			fmt.Sprintf("当前经验：%d\n", user.Exp),
			"继续签到获取更多经验值吧！",
		))
	})
}
