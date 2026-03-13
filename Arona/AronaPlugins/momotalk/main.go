package momotalk

import (
	"context"
	"fmt"
	"log"

	"Shittim/Arona/AronaPlugins/momotalk/dao"
	"Shittim/config"
	"Shittim/pkg/vector"

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

	// 初始化向量模型
	if err := vector.Init(); err != nil {
		log.Printf("Failed to initialize vector model: %v", err)
		// 向量模型初始化失败不影响主功能
	}

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

	// 从配置中获取默认角色信息
	defaultChar := config.AppConfig.Momotalk.Default
	charName := defaultChar.Name
	charDesc := defaultChar.Description

	// 创建 ChatModelAgent
	agent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        charName,
		Description: charDesc,
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
func Chat(ctx context.Context, userID int64, message string, studentName string, onChunk func(string) error) error {
	// 检查是否初始化
	if runner == nil {
		return fmt.Errorf("Momotalk AI not initialized")
	}

	// 生成用户消息的向量嵌入
	messageVector, err := vector.GenerateEmbedding(message)
	if err != nil {
		log.Printf("Failed to generate message embedding: %v", err)
	}

	// 搜索相似的对话和记忆
	searchService := vector.NewSearchService()

	// 搜索相似的对话
	similarConversations, err := searchService.SearchSimilarConversations(userID, messageVector, 3)
	if err != nil {
		log.Printf("Failed to search similar conversations: %v", err)
	}

	// 搜索相似的记忆
	similarMemories, err := searchService.SearchSimilarMemories(studentName, messageVector, 3)
	if err != nil {
		log.Printf("Failed to search similar memories: %v", err)
	}

	// 搜索相似的日常故事
	similarDailyStories, err := searchService.SearchSimilarDailyStories(messageVector, 2)
	if err != nil {
		log.Printf("Failed to search similar daily stories: %v", err)
	}

	// 搜索相似的活动故事
	similarEventStories, err := searchService.SearchSimilarEventStories(messageVector, 2)
	if err != nil {
		log.Printf("Failed to search similar event stories: %v", err)
	}

	// 构建上下文信息
	var contextInfo string
	if len(similarConversations) > 0 {
		contextInfo += "## 相关对话\n"
		for _, conv := range similarConversations {
			for _, msg := range conv.Messages {
				if msg.IsUser {
					contextInfo += "用户: " + msg.Content + "\n"
				} else {
					contextInfo += studentName + ": " + msg.Content + "\n"
				}
			}
			contextInfo += "\n"
		}
	}

	if len(similarMemories) > 0 {
		contextInfo += "## 相关记忆\n"
		for _, memory := range similarMemories {
			contextInfo += memory.Content + "\n\n"
		}
	}

	if len(similarDailyStories) > 0 {
		contextInfo += "## 相关日常故事\n"
		for _, story := range similarDailyStories {
			contextInfo += "标题: " + story.Title + "\n"
			contextInfo += "内容: " + story.Content + "\n\n"
		}
	}

	if len(similarEventStories) > 0 {
		contextInfo += "## 相关活动故事\n"
		for _, story := range similarEventStories {
			contextInfo += "标题: " + story.Title + "\n"
			contextInfo += "内容: " + story.Content + "\n\n"
		}
	}

	// 构建完整的消息
	fullMessage := message
	if contextInfo != "" {
		fullMessage = "# 上下文信息\n" + contextInfo + "\n# 用户问题\n" + message
	}

	// 发送消息并获取响应
	iter := runner.Query(ctx, fullMessage)

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

	// 生成响应的向量嵌入
	responseVector, err := vector.GenerateEmbedding(fullResponse)
	if err != nil {
		log.Printf("Failed to generate response embedding: %v", err)
	}

	// 生成对话主题的向量嵌入
	topicVector, err := vector.GenerateEmbedding("日常对话")
	if err != nil {
		log.Printf("Failed to generate topic embedding: %v", err)
	}

	// 保存对话到数据库
	conversationDAO := dao.NewConversationDAO()
	if err := conversationDAO.SaveConversationWithVectors(userID, message, fullResponse, messageVector, responseVector, topicVector); err != nil {
		log.Printf("Failed to save conversation: %v", err)
	}

	return nil
}
