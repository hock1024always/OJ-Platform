package skill

import (
	"context"
	"encoding/json"
	"fmt"
)

// Skill 定义
type Skill struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
	Handler     SkillHandler           `json:"-"`
}

// SkillHandler Skill 处理函数
type SkillHandler func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)

// Registry Skill 注册表
type Registry struct {
	skills map[string]*Skill
}

func NewRegistry() *Registry {
	r := &Registry{
		skills: make(map[string]*Skill),
	}
	r.registerBuiltInSkills()
	return r
}

// Register 注册 Skill
func (r *Registry) Register(skill *Skill) error {
	if _, exists := r.skills[skill.Name]; exists {
		return fmt.Errorf("skill %s already registered", skill.Name)
	}
	r.skills[skill.Name] = skill
	return nil
}

// Get 获取 Skill
func (r *Registry) Get(name string) (*Skill, bool) {
	skill, exists := r.skills[name]
	return skill, exists
}

// List 列出所有 Skills
func (r *Registry) List() []*Skill {
	list := make([]*Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		list = append(list, skill)
	}
	return list
}

// Execute 执行 Skill
func (r *Registry) Execute(ctx context.Context, name string, input map[string]interface{}) (map[string]interface{}, error) {
	skill, exists := r.skills[name]
	if !exists {
		return nil, fmt.Errorf("skill %s not found", name)
	}

	// 验证输入
	if err := r.validateInput(skill, input); err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// 执行 handler
	return skill.Handler(ctx, input)
}

// validateInput 验证输入参数
func (r *Registry) validateInput(skill *Skill, input map[string]interface{}) error {
	// 简化验证：检查必需字段是否存在
	for field, fieldType := range skill.InputSchema {
		if fieldType == "required" {
			if _, exists := input[field]; !exists {
				return fmt.Errorf("required field %s missing", field)
			}
		}
	}
	return nil
}

// registerBuiltInSkills 注册内置 Skills
func (r *Registry) registerBuiltInSkills() {
	// 代码生成
	r.Register(&Skill{
		Name:        "code_generation",
		Description: "根据需求生成代码",
		InputSchema: map[string]interface{}{
			"language":    "required",
			"requirement": "required",
			"context":     "optional",
		},
		Handler: CodeGenerationHandler,
	})

	// 代码审查
	r.Register(&Skill{
		Name:        "code_review",
		Description: "审查代码质量",
		InputSchema: map[string]interface{}{
			"code":     "required",
			"language": "required",
		},
		Handler: CodeReviewHandler,
	})

	// 调试
	r.Register(&Skill{
		Name:        "debug",
		Description: "调试代码问题",
		InputSchema: map[string]interface{}{
			"code":  "required",
			"error": "required",
		},
		Handler: DebugHandler,
	})

	// 测试生成
	r.Register(&Skill{
		Name:        "test_generation",
		Description: "生成测试用例",
		InputSchema: map[string]interface{}{
			"code":      "required",
			"test_type": "required", // unit/integration/e2e
		},
		Handler: TestGenerationHandler,
	})

	// 系统设计
	r.Register(&Skill{
		Name:        "system_design",
		Description: "设计系统架构",
		InputSchema: map[string]interface{}{
			"requirements": "required",
			"constraints":  "optional",
		},
		Handler: SystemDesignHandler,
	})

	// 部署
	r.Register(&Skill{
		Name:        "deploy",
		Description: "部署应用",
		InputSchema: map[string]interface{}{
			"artifact":    "required",
			"environment": "required",
		},
		Handler: DeployHandler,
	})

	// API 调用
	r.Register(&Skill{
		Name:        "api_call",
		Description: "调用外部 API",
		InputSchema: map[string]interface{}{
			"url":     "required",
			"method":  "required",
			"headers": "optional",
			"body":    "optional",
		},
		Handler: APICallHandler,
	})

	// 数据库查询
	r.Register(&Skill{
		Name:        "db_query",
		Description: "执行数据库查询",
		InputSchema: map[string]interface{}{
			"sql":        "required",
			"connection": "optional",
		},
		Handler: DBQueryHandler,
	})

	// 文件操作
	r.Register(&Skill{
		Name:        "file_operation",
		Description: "文件读写操作",
		InputSchema: map[string]interface{}{
			"action":  "required", // read/write/delete/list
			"path":    "required",
			"content": "optional",
		},
		Handler: FileOperationHandler,
	})

	// 代码搜索
	r.Register(&Skill{
		Name:        "code_search",
		Description: "在代码库中搜索",
		InputSchema: map[string]interface{}{
			"query":    "required",
			"language": "optional",
		},
		Handler: CodeSearchHandler,
	})
}

// Skill Handlers

func CodeGenerationHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	language := input["language"].(string)
	requirement := input["requirement"].(string)

	// TODO: 调用 LLM 生成代码
	return map[string]interface{}{
		"code": fmt.Sprintf("// Generated %s code for: %s\n// TODO: Implement", language, requirement),
		"language": language,
	}, nil
}

func CodeReviewHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	code := input["code"].(string)

	// TODO: 调用 LLM 审查代码
	return map[string]interface{}{
		"issues": []map[string]interface{}{
			{"severity": "info", "message": "Code review completed"},
		},
		"suggestions": []string{"Consider adding more comments"},
		"code": code,
	}, nil
}

func DebugHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	code := input["code"].(string)
	errMsg := input["error"].(string)

	// TODO: 调用 LLM 分析错误
	return map[string]interface{}{
		"analysis": fmt.Sprintf("Error analysis for: %s", errMsg),
		"suggestions": []string{"Check variable initialization", "Verify input parameters"},
		"fixed_code": code,
	}, nil
}

func TestGenerationHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	code := input["code"].(string)
	testType := input["test_type"].(string)

	// TODO: 调用 LLM 生成测试
	return map[string]interface{}{
		"test_code": fmt.Sprintf("// %s tests generated\n// TODO: Implement tests for:\n%s", testType, code),
		"test_type": testType,
	}, nil
}

func SystemDesignHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	requirements := input["requirements"].(string)

	// TODO: 调用 LLM 设计系统
	return map[string]interface{}{
		"architecture": "System architecture design",
		"components": []string{"API Gateway", "Service Mesh", "Database", "Cache"},
		"diagram": fmt.Sprintf("Design for: %s", requirements),
	}, nil
}

func DeployHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	artifact := input["artifact"].(string)
	environment := input["environment"].(string)

	// TODO: 执行部署
	return map[string]interface{}{
		"status":      "deployed",
		"artifact":    artifact,
		"environment": environment,
		"url":         fmt.Sprintf("https://%s.example.com", environment),
	}, nil
}

func APICallHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	url := input["url"].(string)
	method := input["method"].(string)

	// TODO: 实际调用 API
	return map[string]interface{}{
		"status_code": 200,
		"method":      method,
		"url":         url,
		"response":    map[string]interface{}{"message": "API call simulated"},
	}, nil
}

func DBQueryHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	sql := input["sql"].(string)

	// TODO: 实际执行查询
	return map[string]interface{}{
		"sql":    sql,
		"rows":   []map[string]interface{}{},
		"count":  0,
	}, nil
}

func FileOperationHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	action := input["action"].(string)
	path := input["path"].(string)

	switch action {
	case "read":
		return map[string]interface{}{
			"action":  action,
			"path":    path,
			"content": "// File content",
		}, nil
	case "write":
		content := input["content"].(string)
		return map[string]interface{}{
			"action":  action,
			"path":    path,
			"size":    len(content),
			"success": true,
		}, nil
	case "delete":
		return map[string]interface{}{
			"action":  action,
			"path":    path,
			"success": true,
		}, nil
	case "list":
		return map[string]interface{}{
			"action": action,
			"path":   path,
			"files":  []string{},
		}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func CodeSearchHandler(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	query := input["query"].(string)

	// TODO: 实际搜索代码
	return map[string]interface{}{
		"query":   query,
		"results": []map[string]interface{}{},
		"count":   0,
	}, nil
}

// MCP Server 接口

// MCPServer MCP 服务器
type MCPServer struct {
	registry *Registry
}

func NewMCPServer(registry *Registry) *MCPServer {
	return &MCPServer{registry: registry}
}

// HandleRequest 处理 MCP 请求
func (s *MCPServer) HandleRequest(ctx context.Context, requestJSON []byte) ([]byte, error) {
	var request struct {
		SessionID string                 `json:"session_id"`
		Type      string                 `json:"type"`
		Skill     string                 `json:"skill"`
		Params    map[string]interface{} `json:"params"`
	}

	if err := json.Unmarshal(requestJSON, &request); err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"session_id": request.SessionID,
		"type":       request.Type,
	}

	switch request.Type {
	case "list_skills":
		skills := s.registry.List()
		response["success"] = true
		response["data"] = skills

	case "execute":
		result, err := s.registry.Execute(ctx, request.Skill, request.Params)
		if err != nil {
			response["success"] = false
			response["error"] = err.Error()
		} else {
			response["success"] = true
			response["data"] = result
		}

	default:
		response["success"] = false
		response["error"] = "unknown request type"
	}

	return json.Marshal(response)
}
