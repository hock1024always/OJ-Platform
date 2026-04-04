package codegen

import "fmt"

// GenerateAll 为所有支持的语言生成代码
func GenerateAll(sig *FunctionSignature) (map[string]*GeneratedCode, error) {
	generators := []CodeGenerator{
		&GoGenerator{},
		&CGenerator{},
		&CppGenerator{},
		&JavaGenerator{},
	}

	result := make(map[string]*GeneratedCode)
	for _, gen := range generators {
		code, err := gen.Generate(sig)
		if err != nil {
			return nil, fmt.Errorf("generate %s: %w", gen.Language(), err)
		}
		result[gen.Language()] = code
	}
	return result, nil
}

// GenerateForLanguage 为指定语言生成代码
func GenerateForLanguage(sig *FunctionSignature, language string) (*GeneratedCode, error) {
	var gen CodeGenerator
	switch language {
	case "Go":
		gen = &GoGenerator{}
	case "C":
		gen = &CGenerator{}
	case "C++":
		gen = &CppGenerator{}
	case "Java":
		gen = &JavaGenerator{}
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return gen.Generate(sig)
}
