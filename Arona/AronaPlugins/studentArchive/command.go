package studentArchive

import (
	"fmt"
	"strconv"

	"Shittim/Arona/cmd"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type NewStudentModule struct{}

func (m *NewStudentModule) Name() string {
	return "studentArchive"
}

func (m *NewStudentModule) Enter(ctx *zero.Ctx) {
	ctx.Send("已进入学生档案操作系统\n可用命令：\n- newstudent: 创建新学生（使用社团和学校名称）\n- schools: 查看所有学校\n- clubs: 查看所有社团\n- help: 查看帮助\n- exit: 退出系统")
}

func (m *NewStudentModule) Exit(ctx *zero.Ctx) {
	ctx.Send("已退出学生档案操作系统")
}

func (m *NewStudentModule) HandleCommand(cmd string, args []string, ctx *zero.Ctx) bool {
	switch cmd {
	case "help":
		ctx.Send("学生档案操作系统帮助\n- newstudent: 创建新学生（格式：newstudent 姓名 年级 年龄 社团名称 学校名称 身高 爱好 好感度）\n- schools: 查看所有学校\n- clubs: 查看所有社团\n- help: 查看帮助\n- exit: 退出系统")
		return true
	case "schools":
		schools, err := GetAllSchools()
		if err != nil {
			ctx.Send(message.Text("获取学校列表失败：", err.Error()))
			return true
		}

		msg := "学校列表：\n"
		for _, school := range schools {
			msg += fmt.Sprintf("名称: %s, 学生数量: %d\n", school.SchoolName, school.StudentCounts)
		}
		ctx.Send(message.Text(msg))
		return true
	case "clubs":
		clubs, err := GetAllClubs()
		if err != nil {
			ctx.Send(message.Text("获取社团列表失败：", err.Error()))
			return true
		}

		msg := "社团列表：\n"
		for _, club := range clubs {
			msg += fmt.Sprintf("名称: %s,所属: %s, 学生数量: %d\n", club.ClubName, club.SchoolName, club.StudentCounts)
		}
		ctx.Send(message.Text(msg))
		return true
	case "newstudent", "创建学生":
		if len(args) < 8 {
			ctx.Send("请按照以下格式输入学生信息：\nnewstudent 姓名 年级 年龄 社团名称 学校名称 身高 爱好 好感度\n例如：newstudent 小鸟游星野 三年级 17 对策委员会 阿拜多斯高等学院 139cm 睡觉 喜欢")
			return true
		}

		// 解析学生信息
		name := args[0]
		level := args[1]
		age, err := strconv.ParseUint(args[2], 10, 32)
		if err != nil {
			ctx.Send(message.Text("年龄格式错误：", err.Error()))
			return true
		}

		clubName := args[3]
		schoolName := args[4]
		height := args[5]
		love := args[6]
		affection := args[7]

		// 创建学生
		student, err := CreateStudent(name, level, uint(age), clubName, schoolName, height, love, affection)
		if err != nil {
			ctx.Send(message.Text("创建学生失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"🎉 学生创建成功！\n",
			fmt.Sprintf("姓名：%s\n", student.Name),
			fmt.Sprintf("年级：%s\n", student.Level),
			fmt.Sprintf("年龄：%d\n", student.Age),
			fmt.Sprintf("身高：%s\n", student.Height),
			fmt.Sprintf("爱好：%s\n", student.Love),
			fmt.Sprintf("好感度：%s\n", student.Affection),
			"学生档案已成功录入系统！",
		))
		return true
	default:
		return false
	}
}

// 注册学生档案模块
func RegisterModule(cmdSystem interface {
	RegisterModule(module cmd.Module)
}) {
	cmdSystem.RegisterModule(&NewStudentModule{})
}
