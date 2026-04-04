package codegen

import "encoding/json"

// FunctionSignature 结构化函数签名定义
type FunctionSignature struct {
	Name       string  `json:"name"`
	Params     []Param `json:"params"`
	ReturnType string  `json:"return_type"`
}

// Param 函数参数
type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// InputConstraint 输入约束（用于自动生成测试数据）
type InputConstraint struct {
	ParamName string `json:"param_name"`
	MinVal    int    `json:"min_val,omitempty"`
	MaxVal    int    `json:"max_val,omitempty"`
	MinLen    int    `json:"min_len,omitempty"`
	MaxLen    int    `json:"max_len,omitempty"`
	// 链表/树的节点数范围
	MinNodes int `json:"min_nodes,omitempty"`
	MaxNodes int `json:"max_nodes,omitempty"`
}

// GeneratedCode 生成的代码
type GeneratedCode struct {
	Language         string `json:"language"`
	FunctionTemplate string `json:"function_template"`
	DriverCode       string `json:"driver_code"`
}

// ParseSignature 从 JSON 字符串解析函数签名
func ParseSignature(s string) (*FunctionSignature, error) {
	var sig FunctionSignature
	if err := json.Unmarshal([]byte(s), &sig); err != nil {
		return nil, err
	}
	return &sig, nil
}

// ToJSON 序列化为 JSON
func (sig *FunctionSignature) ToJSON() string {
	data, _ := json.Marshal(sig)
	return string(data)
}

// 支持的类型常量
const (
	TypeInt       = "int"
	TypeFloat     = "float"
	TypeString    = "string"
	TypeBool      = "bool"
	TypeIntArray  = "[]int"
	TypeStrArray  = "[]string"
	TypeInt2D     = "[][]int"
	TypeStr2D     = "[][]string"
	TypeByteArray = "[]byte"
	TypeByte2D    = "[][]byte"
	TypeListNode  = "ListNode"
	TypeTreeNode  = "TreeNode"
)

// CodeGenerator 代码生成器接口
type CodeGenerator interface {
	Generate(sig *FunctionSignature) (*GeneratedCode, error)
	Language() string
}
