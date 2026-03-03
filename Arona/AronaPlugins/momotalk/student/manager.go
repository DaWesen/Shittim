package student

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// StudentInfo 学生信息扩展，包含对话相关的配置
type StudentInfo struct {
	Name           string   `yaml:"name" json:"name"`
	School         string   `yaml:"school" json:"school"`
	Organization   string   `yaml:"organization" json:"organization"`
	Personality    []string `yaml:"personality" json:"personality"`
	SpeechStyle    []string `yaml:"speech_style" json:"speech_style"`
	SelfReference  []string `yaml:"self_reference" json:"self_reference"`
	SpecialAbility []string `yaml:"special_ability" json:"special_ability"`
	Memory         []string `yaml:"memory" json:"memory"`
	ImportantRules []string `yaml:"important_rules" json:"important_rules"`
	Relationships  []string `yaml:"relationships" json:"relationships"`
	Interests      []string `yaml:"interests" json:"interests"`
	Background     []string `yaml:"background" json:"background"`
	Catchphrases   []string `yaml:"catchphrases" json:"catchphrases"`
	Appearance     []string `yaml:"appearance" json:"appearance"`
	VoiceActor     string   `yaml:"voice_actor" json:"voice_actor"`

	// 对话相关配置
	Greeting       string   `json:"greeting"`
	ExampleDialogs []string `json:"example_dialogs"`
	Temperature    float64  `json:"temperature"`
}

// Manager 学生管理器
type Manager struct {
	students map[string]*StudentInfo
}

// NewManager 创建学生管理器
func NewManager() (*Manager, error) {
	manager := &Manager{
		students: make(map[string]*StudentInfo),
	}

	// 从YAML文件加载学生信息
	if err := manager.loadStudentsFromYAML(); err != nil {
		return nil, err
	}

	return manager, nil
}

// loadStudentsFromYAML 从YAML文件加载学生信息
func (m *Manager) loadStudentsFromYAML() error {
	// 定义可能的YAML文件路径
	yamlPaths := []string{
		"./student",
		"../student",
		"../../student",
		"../../../student",
	}

	// 遍历所有可能的路径
	for _, basePath := range yamlPaths {
		// 检查路径是否存在
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			continue
		}

		// 读取目录中的所有YAML文件
		files, err := ioutil.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
				// 读取YAML文件
				filePath := filepath.Join(basePath, file.Name())
				data, err := ioutil.ReadFile(filePath)
				if err != nil {
					continue
				}

				// 解析YAML文件
				var studentInfo StudentInfo
				if err := yaml.Unmarshal(data, &studentInfo); err != nil {
					continue
				}

				// 初始化对话相关配置
				studentInfo.Greeting = fmt.Sprintf("你好，我是%s，来自%s的%s成员。", studentInfo.Name, studentInfo.School, studentInfo.Organization)
				studentInfo.ExampleDialogs = studentInfo.Catchphrases
				studentInfo.Temperature = 0.7

				// 添加到学生列表
				m.students[studentInfo.Name] = &studentInfo
			}
		}
	}

	// 确保至少加载了圣园未花
	if len(m.students) == 0 {
		return fmt.Errorf("没有加载到任何学生信息，请确保student目录下有YAML文件")
	}

	return nil
}

// GetStudent 获取学生信息
func (m *Manager) GetStudent(name string) (*StudentInfo, error) {
	student, ok := m.students[name]
	if !ok {
		return nil, fmt.Errorf("student %s not found", name)
	}
	return student, nil
}

// ListStudents 列出所有学生
func (m *Manager) ListStudents() []string {
	names := make([]string, 0, len(m.students))
	for name := range m.students {
		names = append(names, name)
	}
	return names
}

// GetStudentDetails 获取学生详细信息
func (m *Manager) GetStudentDetails(name string) (map[string]interface{}, error) {
	student, err := m.GetStudent(name)
	if err != nil {
		return nil, err
	}

	details := map[string]interface{}{
		"name":        student.Name,
		"club":        student.Organization,
		"school":      student.School,
		"greeting":    student.Greeting,
		"personality": student.Personality,
	}

	return details, nil
}
