package momotalk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Shittim/Arona/AronaPlugins/momotalk/student"
	"Shittim/Arona/cmd"
	"Shittim/pkg/database"
	"Shittim/pkg/models"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 聊天模式类型
type ChatMode string

// 定义聊天模式常量
const (
	DailyChatMode       ChatMode = "daily"     // 日常聊天模式
	EventChatMode       ChatMode = "event"     // 活动聊天模式
	ExclusiveMemoryMode ChatMode = "exclusive" // 独家记忆模式
)

// Momotalk模块
type MomotalkModule struct {
	studentManager *student.Manager
	currentStudent []string
	currentMode    ChatMode
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

	// 构建聊天模式列表
	modeList := "可用聊天模式：\n"
	modeList += "- daily: 日常聊天模式（支持1-3名学生）\n"
	modeList += "- event: 活动聊天模式（支持2-4名学生）\n"
	modeList += "- exclusive: 独家记忆模式（仅支持1名学生）\n"

	ctx.Send("已进入 Momotalk 聊天系统\n" + roleList + "\n" + modeList + "\n可用命令：\n- mode <模式>: 切换聊天模式\n- select <学生姓名>: 选择聊天角色\n- chat <消息>: 与 AI 聊天\n- exit: 退出 Momotalk 系统")
}

// 退出模块时执行的操作
func (m *MomotalkModule) Exit(ctx *zero.Ctx) {
	ctx.Send("已退出 Momotalk 系统")
}

// 处理模块内的命令
func (m *MomotalkModule) HandleCommand(cmd string, args []string, ctx *zero.Ctx) bool {
	// 检查是否是exit命令
	if cmd == "exit" || cmd == "/exit" {
		m.Exit(ctx)
		m.currentStudent = []string{}
		m.currentMode = ""
		return true
	}

	// 处理模式切换命令
	if cmd == "mode" || cmd == "模式" {
		if len(args) == 0 {
			ctx.Send("请输入要切换的聊天模式，例如：mode daily")
			return true
		}

		mode := ChatMode(args[0])
		// 验证模式是否有效
		switch mode {
		case DailyChatMode, EventChatMode, ExclusiveMemoryMode:
			m.currentMode = mode
			m.currentStudent = []string{} // 切换模式时清空已选择的角色
			var modeDescription string
			switch mode {
			case DailyChatMode:
				modeDescription = "日常聊天模式（支持1-3名学生）"
			case EventChatMode:
				modeDescription = "活动聊天模式（支持2-4名学生）"
			case ExclusiveMemoryMode:
				modeDescription = "独家记忆模式（仅支持1名学生）"
			}
			ctx.Send(fmt.Sprintf("已切换到%s，请重新选择角色", modeDescription))
		default:
			ctx.Send("无效的聊天模式，请选择：daily、event 或 exclusive")
		}
		return true
	}

	// 处理选择角色命令
	if cmd == "select" || cmd == "选择" {
		if len(args) == 0 {
			ctx.Send("请输入要选择的学生姓名，例如：select 圣园未花")
			return true
		}

		// 检查是否已经选择了模式
		if m.currentMode == "" {
			ctx.Send("请先选择聊天模式，例如：mode daily")
			return true
		}

		studentName := strings.Join(args, " ")
		studentInfo, err := m.studentManager.GetStudent(studentName)
		if err != nil {
			ctx.Send("学生不存在: " + err.Error())
			return true
		}

		// 检查角色数量是否符合当前模式的要求
		maxStudents := 1
		switch m.currentMode {
		case DailyChatMode:
			maxStudents = 3
		case EventChatMode:
			maxStudents = 4
		case ExclusiveMemoryMode:
			maxStudents = 1
		}

		// 检查是否已经达到最大角色数量
		if len(m.currentStudent) >= maxStudents {
			ctx.Send(fmt.Sprintf("当前模式最多只能选择%d名学生", maxStudents))
			return true
		}

		// 检查是否已经选择了该角色
		for _, name := range m.currentStudent {
			if name == studentInfo.Name {
				ctx.Send("该角色已经被选择")
				return true
			}
		}

		// 添加角色
		m.currentStudent = append(m.currentStudent, studentInfo.Name)

		// 构建角色列表消息
		roleList := "已选择角色：\n"
		for _, name := range m.currentStudent {
			sInfo, _ := m.studentManager.GetStudent(name)
			if sInfo != nil {
				roleList += fmt.Sprintf("- %s (所属组织: %s, 所属学校: %s)\n",
					sInfo.Name, sInfo.Organization, sInfo.School)
			}
		}

		// 检查是否达到当前模式的最小角色数量
		minStudents := 1
		if m.currentMode == EventChatMode {
			minStudents = 2
		}

		if len(m.currentStudent) < minStudents {
			ctx.Send(roleList + fmt.Sprintf("\n当前模式至少需要选择%d名学生，请继续选择", minStudents))
		} else {
			ctx.Send(roleList + "\n角色选择完成，可以开始聊天了")
		}

		return true
	}

	// 处理聊天命令或直接消息
	if cmd == "chat" || cmd == "聊天" || len(m.currentStudent) > 0 {
		// 构建消息
		var messageText string
		if cmd == "chat" || cmd == "聊天" {
			if len(args) == 0 {
				ctx.Send("请输入要发送的消息，例如：chat 你好")
				return true
			}
			messageText = strings.Join(args, " ")
		} else {
			messageText = cmd
			if len(args) > 0 {
				messageText += " " + strings.Join(args, " ")
			}
		}

		// 检查是否已经选择了模式
		if m.currentMode == "" {
			ctx.Send("请先选择聊天模式，例如：mode daily")
			return true
		}

		// 检查是否已经选择了角色
		if len(m.currentStudent) == 0 {
			ctx.Send("请先选择角色，例如：select 圣园未花")
			return true
		}

		// 检查角色数量是否符合当前模式的要求
		minStudents := 1
		if m.currentMode == EventChatMode {
			minStudents = 2
		}

		if len(m.currentStudent) < minStudents {
			ctx.Send(fmt.Sprintf("当前模式至少需要选择%d名学生，请继续选择", minStudents))
			return true
		}

		// 构建提示词
		prompt := m.buildPrompt(messageText)

		// 调用 AI 聊天功能（流式）
		go func() {
			// 流式接收并输出消息
			var fullResponse string
			var lastResponse string
			err := Chat(context.Background(), prompt, func(chunk string) error {
				fullResponse += chunk
				// 只有当内容发生变化时才发送消息
				if fullResponse != lastResponse {
					ctx.Send(message.Text(fullResponse))
					lastResponse = fullResponse
				}
				return nil
			})

			if err != nil {
				// 发送错误消息
				ctx.Send(message.Text("AI 处理失败: " + err.Error()))
				return
			}

			// 如果是独家记忆模式，保存聊天内容为记忆
			if m.currentMode == ExclusiveMemoryMode && len(m.currentStudent) == 1 {
				memoryContent := fmt.Sprintf("用户: %s\n%s: %s", messageText, m.currentStudent[0], fullResponse)
				emotionTag := "聊天记忆"
				err := m.saveExclusiveMemory(m.currentStudent[0], memoryContent, emotionTag)
				if err != nil {
					fmt.Printf("保存独家记忆失败: %v\n", err)
				}
			}
		}()

		return true
	}

	return false
}

// 构建提示词
func (m *MomotalkModule) buildPrompt(messageText string) string {
	if len(m.currentStudent) == 1 {
		// 单人聊天
		studentInfo, err := m.studentManager.GetStudent(m.currentStudent[0])
		if err != nil {
			return "获取学生信息失败"
		}

		// 构建基础提示
		prompt := studentInfo.Prompt

		// 如果是独家记忆模式，添加存储的记忆
		if m.currentMode == ExclusiveMemoryMode {
			memories := m.loadExclusiveMemories(studentInfo.Name)
			if len(memories) > 0 {
				prompt += "\n\n## 独家记忆\n"
				for _, memory := range memories {
					prompt += fmt.Sprintf("- %s (情感标签: %s)\n", memory.Content, memory.EmotionTag)
				}
			}
		}

		// 检查是否需要查询数据库
		if strings.Contains(messageText, "学生") || strings.Contains(messageText, "学校") || strings.Contains(messageText, "社团") {
			// 添加数据库查询结果
			prompt += "\n\n## 数据库信息\n"
			prompt += m.getDatabaseInfo(messageText)
		}

		prompt += "\n\n用户说：" + messageText
		return prompt
	} else {
		// 多人聊天，构建群聊提示
		prompt := "# 群聊场景\n\n"
		prompt += "参与聊天的学生："
		for _, name := range m.currentStudent {
			studentInfo, err := m.studentManager.GetStudent(name)
			if err == nil {
				prompt += studentInfo.Name + "、"
			}
		}
		prompt = strings.TrimSuffix(prompt, "、") + "\n\n"
		prompt += "用户说：" + messageText
		return prompt
	}
}

// 加载独家记忆
func (m *MomotalkModule) loadExclusiveMemories(studentName string) []models.ExclusiveMemory {
	var memories []models.ExclusiveMemory
	db := database.GetDB()
	result := db.Where("student_name = ?", studentName).Find(&memories)
	if result.Error != nil {
		fmt.Printf("加载独家记忆失败: %v\n", result.Error)
		return []models.ExclusiveMemory{}
	}
	return memories
}

// 保存独家记忆
func (m *MomotalkModule) saveExclusiveMemory(studentName, content, emotionTag string) error {
	memory := models.ExclusiveMemory{
		StoryBase: models.StoryBase{
			Title:     fmt.Sprintf("与%s的记忆", studentName),
			Content:   content,
			StoryType: "exclusive_memory",
		},
		StudentName:   studentName,
		IntimacyLevel: 1,
		MemoryType:    "chat_memory",
		SpecialDate:   time.Now().Format("2006-01-02"),
		EmotionTag:    emotionTag,
		LastRecallAt:  time.Now(),
	}

	db := database.GetDB()
	result := db.Create(&memory)
	return result.Error
}

// 获取数据库信息
func (m *MomotalkModule) getDatabaseInfo(query string) string {
	// 获取数据库连接
	db := database.GetDB()

	// 构建查询
	dbQuery := db.Model(&models.Student{})

	// 根据查询内容添加条件
	if strings.Contains(query, "三一") {
		dbQuery = dbQuery.Where("school_name LIKE ?", "%三一%")
	} else if strings.Contains(query, "格黑娜") {
		dbQuery = dbQuery.Where("school_name LIKE ?", "%格黑娜%")
	} else if strings.Contains(query, "阿拜多斯") {
		dbQuery = dbQuery.Where("school_name LIKE ?", "%阿拜多斯%")
	} else if strings.Contains(query, "茶会") {
		dbQuery = dbQuery.Where("club_name LIKE ?", "%茶会%")
	}

	// 执行查询
	var students []models.Student
	result := dbQuery.Find(&students)
	if result.Error != nil || len(students) == 0 {
		return "未找到相关学生信息"
	}

	// 构建响应
	response := ""
	for i, student := range students {
		if i > 0 {
			response += "\n"
		}
		response += fmt.Sprintf("- %s (学校: %s, 社团: %s, 年龄: %d)", student.Name, student.SchoolName, student.ClubName, student.Age)
	}

	return response
}

// 注册 Momotalk 模块
func RegisterModule(cmdSystem interface {
	RegisterModule(module cmd.Module)
}) {
	cmdSystem.RegisterModule(&MomotalkModule{})
}
