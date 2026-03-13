package vector

import (
	"context"
	"log"
	"math"

	"Shittim/config"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/pgvector/pgvector-go"
)

// 全局变量
var (
	embedder embedding.Embedder
)

// 初始化向量模型
func Init() error {
	log.Println("Initializing vector embedding model...")

	// 从配置中获取 Eino 配置
	einoConfig := config.AppConfig.AI.Eino

	// 检查配置是否存在
	if einoConfig.APIKey == "" {
		return nil
	}

	// 创建嵌入模型
	var err error
	embedder, err = openai.NewEmbedder(context.Background(), &openai.EmbeddingConfig{
		APIKey:  einoConfig.APIKey,
		Model:   "text-embedding-ada-002",
		BaseURL: einoConfig.BaseURL,
	})
	if err != nil {
		log.Printf("Failed to create embedding model: %v", err)
		return err
	}

	log.Println("Vector embedding model initialized successfully")
	return nil
}

// GenerateEmbedding 生成文本的向量嵌入
func GenerateEmbedding(text string) (pgvector.Vector, error) {
	if embedder == nil {
		return pgvector.Vector{}, nil
	}

	// 生成嵌入
	vectors, err := embedder.EmbedStrings(context.Background(), []string{text})
	if err != nil {
		log.Printf("Failed to generate embedding: %v", err)
		return pgvector.Vector{}, err
	}

	// 检查结果
	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return pgvector.Vector{}, nil
	}

	// 转换为float32
	float32Vector := make([]float32, len(vectors[0]))
	for i, v := range vectors[0] {
		float32Vector[i] = float32(v)
	}

	// 转换为向量
	return pgvector.NewVector(float32Vector), nil
}

// Similarity 计算两个向量的相似度
func Similarity(v1, v2 pgvector.Vector) float64 {
	// 计算余弦相似度
	// 余弦相似度 = 1 - 余弦距离
	return 1 - cosineDistance(v1, v2)
}

// cosineDistance 计算两个向量的余弦距离
func cosineDistance(v1, v2 pgvector.Vector) float64 {
	// 获取向量数据
	vec1 := v1.Slice()
	vec2 := v2.Slice()

	// 检查向量维度是否相同
	if len(vec1) != len(vec2) {
		return 1.0 // 维度不同，返回最大距离
	}

	// 计算点积
	dotProduct := 0.0
	mag1 := 0.0
	mag2 := 0.0

	for i := range vec1 {
		dotProduct += float64(vec1[i]) * float64(vec2[i])
		mag1 += float64(vec1[i]) * float64(vec1[i])
		mag2 += float64(vec2[i]) * float64(vec2[i])
	}

	// 计算模长
	mag1 = math.Sqrt(mag1)
	mag2 = math.Sqrt(mag2)

	// 避免除以零
	if mag1 == 0 || mag2 == 0 {
		return 1.0
	}

	// 计算余弦相似度
	similarity := dotProduct / (mag1 * mag2)

	// 计算余弦距离
	distance := 1 - similarity

	return distance
}
