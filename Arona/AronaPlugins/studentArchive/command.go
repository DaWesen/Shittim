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
	ctx.Send("已进入学生档案操作系统\n可用命令：\n- newschool: 创建新学院\n- newclub: 创建新社团（需要挂靠到学院）\n- newstudent: 创建新学生（使用社团和学校名称）\n- student: 查看学生信息\n- update: 更新学生档案信息\n- attention: 更新学生关注状态\n- schools: 查看所有学校\n- clubs: 查看所有社团\n- help: 查看帮助\n- exit: 退出系统")
}

func (m *NewStudentModule) Exit(ctx *zero.Ctx) {
	ctx.Send("已退出学生档案操作系统")
}

func (m *NewStudentModule) HandleCommand(cmd string, args []string, ctx *zero.Ctx) bool {
	switch cmd {
	case "help":
		ctx.Send("学生档案操作系统帮助\n- newschool: 创建新学院（格式：newschool 学院名称）\n- newclub: 创建新社团（格式：newclub 社团名称 学院名称）\n- newstudent: 创建新学生（格式：newstudent 姓名 年级 年龄 社团名称 学校名称 身高 爱好 好感度）\n- student: 查看学生信息（格式：student 姓名）\n- update: 更新学生档案信息（格式：update 姓名 年级 年龄 社团名称 学校名称 身高 爱好 关注状态 好感度）\n- attention: 更新学生关注状态（格式：attention 姓名 状态）\n- schools: 查看所有学校\n- clubs: 查看所有社团\n- help: 查看帮助\n- exit: 退出系统")
		return true
	case "newschool", "创建学院":
		if len(args) < 1 {
			ctx.Send("请按照以下格式输入学院名称：\nnewschool 学院名称\n例如：newschool 阿拜多斯学院")
			return true
		}

		schoolName := args[0]
		school, err := CreateSchool(schoolName)
		if err != nil {
			ctx.Send(message.Text("创建学院失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"🎉 学院创建成功！\n",
			fmt.Sprintf("学院名称：%s\n", school.SchoolName),
			fmt.Sprintf("学生数量：%d\n", school.StudentCounts),
			"学院已成功创建并添加到系统！",
		))
		return true
	case "newclub", "创建社团":
		if len(args) < 2 {
			ctx.Send("请按照以下格式输入社团名称和学院名称：\nnewclub 社团名称 学院名称\n例如：newclub 阿拜多斯 阿拜多斯学院")
			return true
		}

		clubName := args[0]
		schoolName := args[1]
		club, err := CreateClub(clubName, schoolName)
		if err != nil {
			ctx.Send(message.Text("创建社团失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"🎉 社团创建成功！\n",
			fmt.Sprintf("社团名称：%s\n", club.ClubName),
			fmt.Sprintf("所属学院：%s\n", club.SchoolName),
			fmt.Sprintf("学生数量：%d\n", club.StudentCounts),
			"社团已成功创建并添加到系统！",
		))
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
	case "student", "查看学生":
		if len(args) < 1 {
			ctx.Send("请按照以下格式输入学生姓名：\nstudent 姓名\n例如：student 小鸟游星野")
			return true
		}

		name := args[0]
		student, err := GetStudent(name)
		if err != nil {
			ctx.Send(message.Text("获取学生信息失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"📋 学生信息\n",
			fmt.Sprintf("姓名：%s\n", student.Name),
			fmt.Sprintf("年级：%s\n", student.Level),
			fmt.Sprintf("年龄：%d\n", student.Age),
			fmt.Sprintf("所属社团：%s\n", student.ClubName),
			fmt.Sprintf("所属学校：%s\n", student.SchoolName),
			fmt.Sprintf("身高：%s\n", student.Height),
			fmt.Sprintf("爱好：%s\n", student.Love),
			fmt.Sprintf("关注状态：%s\n", student.UnderEye),
			fmt.Sprintf("好感度：%s\n", student.Affection),
		))
		return true
	case "update", "更新学生":
		if len(args) < 9 {
			ctx.Send("请按照以下格式输入学生信息：\nupdate 姓名 年级 年龄 社团名称 学校名称 身高 爱好 关注状态 好感度\n例如：update 小鸟游星野 三年级 17 对策委员会 阿拜多斯高等学院 139cm 睡觉 重点关注 喜欢")
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
		underEye := args[7]
		affection := args[8]

		// 验证关注状态值
		validStatus := false
		validStatuses := []string{"默认", "重点关注"}
		for _, s := range validStatuses {
			if underEye == s {
				validStatus = true
				break
			}
		}

		if !validStatus {
			ctx.Send("无效的关注状态，请使用以下状态之一：默认, 重点关注")
			return true
		}

		// 更新学生档案
		student, err := UpdateStudentArchive(name, level, uint(age), clubName, schoolName, height, love, underEye, affection)
		if err != nil {
			ctx.Send(message.Text("更新学生档案失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"🎉 学生档案更新成功！\n",
			fmt.Sprintf("姓名：%s\n", student.Name),
			fmt.Sprintf("年级：%s\n", student.Level),
			fmt.Sprintf("年龄：%d\n", student.Age),
			fmt.Sprintf("所属社团：%s\n", student.ClubName),
			fmt.Sprintf("所属学校：%s\n", student.SchoolName),
			fmt.Sprintf("身高：%s\n", student.Height),
			fmt.Sprintf("爱好：%s\n", student.Love),
			fmt.Sprintf("关注状态：%s\n", student.UnderEye),
			fmt.Sprintf("好感度：%s\n", student.Affection),
			"学生档案已成功更新！",
		))
		return true
	case "attention", "更新关注状态":
		if len(args) < 2 {
			ctx.Send("请按照以下格式输入学生姓名和关注状态：\nattention 姓名 状态\n例如：attention 小鸟游星野 重点关注\n可用状态：默认, 重点关注")
			return true
		}

		name := args[0]
		status := args[1]

		// 验证状态值
		validStatus := false
		validStatuses := []string{"默认", "重点关注"}
		for _, s := range validStatuses {
			if status == s {
				validStatus = true
				break
			}
		}

		if !validStatus {
			ctx.Send("无效的关注状态，请使用以下状态之一：默认, 重点关注")
			return true
		}

		student, err := UpdateAttentionStatus(name, status)
		if err != nil {
			ctx.Send(message.Text("更新关注状态失败：", err.Error()))
			return true
		}

		ctx.Send(message.Text(
			"🎉 关注状态更新成功！\n",
			fmt.Sprintf("学生：%s\n", student.Name),
			fmt.Sprintf("新状态：%s\n", student.UnderEye),
		))
		return true
	case "newstudent", "创建学生":
		if len(args) < 8 {
			ctx.Send("请按照以下格式输入学生信息：\nnewstudent 姓名 年级 年龄 社团名称 学校名称 身高 爱好 好感度\n例如：newstudent 小鸟游星野 三年级 17 对策委员会 阿拜多斯高等学院 139cm 音乐 喜欢")
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
			fmt.Sprintf("关注状态：%s\n", student.UnderEye),
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
