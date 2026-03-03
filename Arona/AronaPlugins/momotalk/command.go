package momotalk

import (
	"context"
	"fmt"
	"strings"

	"Shittim/Arona/AronaPlugins/momotalk/student"
	"Shittim/Arona/cmd"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// Momotalk模块
type MomotalkModule struct {
	studentManager *student.Manager
	currentStudent string
}

// 返回模块名称
func (m *MomotalkModule) Name() string {
	return "momotalk"
}

// Enter 进入模块时执行的操作
func (m *MomotalkModule) Enter(ctx *zero.Ctx) {
	// 初始化学生管理器
	var err error
	m.studentManager, err = student.NewManager()
	if err != nil {
		ctx.Send("初始化学生管理器失败: " + err.Error())
		return
	}

	// 从学生管理器中加载学生列表
	studentNames := m.studentManager.ListStudents()

	if len(studentNames) == 0 {
		ctx.Send("没有可用的学生角色")
		return
	}

	// 构建角色列表消息
	roleList := "可用角色：\n"
	for _, name := range studentNames {
		details, _ := m.studentManager.GetStudentDetails(name)
		roleList += fmt.Sprintf("- %s (所属组织: %s, 所属学校: %s)\n",
			details["name"], details["club"], details["school"])
	}

	ctx.Send("已进入 Momotalk 聊天系统\n" + roleList + "\n可用命令：\n- select <学生姓名>: 选择聊天角色\n- chat <消息>: 与 AI 聊天\n- exit: 退出 Momotalk 系统")
}

// 退出模块时执行的操作
func (m *MomotalkModule) Exit(ctx *zero.Ctx) {
	ctx.Send("已退出 Momotalk 系统")
}

// 处理模块内的命令
func (m *MomotalkModule) HandleCommand(cmd string, args []string, ctx *zero.Ctx) bool {
	// 如果已经选择了角色，直接将输入作为聊天消息处理
	if m.currentStudent != "" {
		// 构建消息
		messageText := cmd
		if len(args) > 0 {
			for _, arg := range args {
				messageText += " " + arg
			}
		}

		// 获取当前学生信息
		studentInfo, err := m.studentManager.GetStudent(m.currentStudent)
		if err != nil {
			ctx.Send("获取学生信息失败: " + err.Error())
			return true
		}

		// 构建角色提示
		prompt := fmt.Sprintf("你是 %s，一名来自%s的%s成员。\n你的性格：%s\n你的说话风格：%s\n你的自称：%s\n你的特殊能力：%s\n你的兴趣：%s\n你的背景：%s\n示例对话：%s\n\n用户说：%s",
			studentInfo.Name,
			studentInfo.School,
			studentInfo.Organization,
			strings.Join(studentInfo.Personality, "、"),
			strings.Join(studentInfo.SpeechStyle, "、"),
			strings.Join(studentInfo.SelfReference, "、"),
			strings.Join(studentInfo.SpecialAbility, "、"),
			strings.Join(studentInfo.Interests, "、"),
			strings.Join(studentInfo.Background, "。"),
			strings.Join(studentInfo.ExampleDialogs, " "),
			messageText)

		// 调用 AI 聊天功能（流式）
		go func() {
			// 发送 "正在思考..." 提示
			ctx.Send(message.Text("正在思考..."))

			// 流式接收并输出消息
			var fullResponse string
			err := Chat(context.Background(), prompt, func(chunk string) error {
				fullResponse += chunk
				// 每次收到新内容就发送一次
				ctx.Send(message.Text(fullResponse))
				return nil
			})

			if err != nil {
				// 发送错误消息
				ctx.Send(message.Text("AI 处理失败: " + err.Error()))
				return
			}
		}()

		return true
	}

	// 未选择角色时，处理命令
	switch cmd {
	case "select", "选择":
		if len(args) == 0 {
			ctx.Send("请输入要选择的学生姓名，例如：select 三宅美玖")
			return true
		}

		studentName := strings.Join(args, " ")
		studentInfo, err := m.studentManager.GetStudent(studentName)
		if err != nil {
			ctx.Send("学生不存在: " + err.Error())
			return true
		}

		m.currentStudent = studentInfo.Name
		ctx.Send(fmt.Sprintf("已选择角色：%s\n%s\n所属组织: %s\n所属学校: %s",
			studentInfo.Name, studentInfo.Greeting, studentInfo.Organization, studentInfo.School))
		return true

	case "chat", "聊天":
		if len(args) == 0 {
			ctx.Send("请输入要发送的消息，例如：chat 你好")
			return true
		}

		// 构建消息
		messageText := ""
		for _, arg := range args {
			messageText += arg + " "
		}

		// 获取当前学生信息
		studentInfo, err := m.studentManager.GetStudent(m.currentStudent)
		if err != nil {
			ctx.Send("请先选择一个角色，例如：select 圣园未花")
			return true
		}

		// 构建角色提示
		prompt := fmt.Sprintf("你是 %s，一名来自%s的%s成员。\n你的性格：%s\n你的说话风格：%s\n你的自称：%s\n你的特殊能力：%s\n你的兴趣：%s\n你的背景：%s\n示例对话：%s\n\n用户说：%s",
			studentInfo.Name,
			studentInfo.School,
			studentInfo.Organization,
			strings.Join(studentInfo.Personality, "、"),
			strings.Join(studentInfo.SpeechStyle, "、"),
			strings.Join(studentInfo.SelfReference, "、"),
			strings.Join(studentInfo.SpecialAbility, "、"),
			strings.Join(studentInfo.Interests, "、"),
			strings.Join(studentInfo.Background, "。"),
			strings.Join(studentInfo.ExampleDialogs, " "),
			messageText)

		// 调用 AI 聊天功能（流式）
		go func() {
			// 发送 "正在思考..." 提示
			ctx.Send(message.Text("正在思考..."))

			// 流式接收并输出消息
			var fullResponse string
			err := Chat(context.Background(), prompt, func(chunk string) error {
				fullResponse += chunk
				// 每次收到新内容就发送一次
				ctx.Send(message.Text(fullResponse))
				return nil
			})

			if err != nil {
				// 发送错误消息
				ctx.Send(message.Text("AI 处理失败: " + err.Error()))
				return
			}
		}()

		return true
	default:
		return false
	}
}

// 注册 Momotalk 模块
func RegisterModule(cmdSystem interface {
	RegisterModule(module cmd.Module)
}) {
	cmdSystem.RegisterModule(&MomotalkModule{})
}
