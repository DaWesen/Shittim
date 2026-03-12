package momotalk

import (
	"context"
	"fmt"
	"log"

	"Shittim/Arona/AronaPlugins/momotalk/dao"
	"Shittim/config"
	"Shittim/pkg/database"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// 全局变量
var (
	chatModel *openai.ChatModel
	runner    *adk.Runner
)

// Init 初始化 Momotalk AI
func Init() error {
	fmt.Println("Initializing Momotalk AI with Eino...")

	// 初始化数据库
	database.InitDatabase()
	database.AutoMigrate()

	// 从配置中获取 Eino 配置
	einoConfig := config.AppConfig.AI.Eino

	// 检查配置是否存在
	if einoConfig.APIKey == "" {
		return fmt.Errorf("Eino API key not configured")
	}

	// 创建 OpenAI 聊天模型
	var err error
	chatModel, err = openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:   einoConfig.Model,
		APIKey:  einoConfig.APIKey,
		BaseURL: einoConfig.BaseURL,
	})
	if err != nil {
		log.Printf("Failed to create chat model: %v", err)
		return err
	}

	// 创建 ChatModelAgent
	agent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "圣园未花",
		Description: "圣园未花是三一综合学园茶会的成员，性格乐观开朗，热心助人，有点冒失，但内心非常关心朋友。她说话时句尾喜欢加上\"♪\"或\"！\"，语气活力满满、略带调皮。",
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to create agent: %v", err)
		return err
	}

	// 创建运行器
	runner = adk.NewRunner(context.Background(), adk.RunnerConfig{Agent: agent})

	return nil
}

// Chat 处理 AI 对话（流式）
func Chat(ctx context.Context, userID int64, message string, onChunk func(string) error) error {
	// 检查是否初始化
	if runner == nil {
		return fmt.Errorf("Momotalk AI not initialized")
	}

	// 发送消息并获取响应
	iter := runner.Query(ctx, message)

	// 收集完整的响应
	var fullResponse string

	// 流式处理响应
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		// 处理事件并流式输出
		if event != nil {
			// 使用GetMessage函数从AgentEvent中获取消息
			msg, _, err := adk.GetMessage(event)
			if err == nil && msg != nil {
				// 获取消息内容
				content := msg.Content
				if content != "" {
					fullResponse += content
					if err := onChunk(content); err != nil {
						return err
					}
				}
			}
		}
	}

	// 保存对话到数据库
	conversationDAO := dao.NewConversationDAO()
	if err := conversationDAO.SaveConversation(userID, message, fullResponse); err != nil {
		log.Printf("Failed to save conversation: %v", err)
	}

	return nil
}
