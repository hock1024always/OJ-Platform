package codegen

import (
	"fmt"
	"strings"
)

type CGenerator struct{}

func (g *CGenerator) Language() string { return "C" }

func (g *CGenerator) Generate(sig *FunctionSignature) (*GeneratedCode, error) {
	// C 不支持 ListNode/TreeNode/2D 数组等复杂类型时做简化处理
	template := g.genTemplate(sig)
	driver, err := g.genDriver(sig)
	if err != nil {
		return nil, err
	}
	return &GeneratedCode{
		Language:         "C",
		FunctionTemplate: template,
		DriverCode:       driver,
	}, nil
}

func (g *CGenerator) genTemplate(sig *FunctionSignature) string {
	params := make([]string, 0)
	for _, p := range sig.Params {
		params = append(params, cParamDecls(p.Type, p.Name)...)
	}
	ret := cType(sig.ReturnType)
	if ret == "" {
		ret = "void"
	}
	// C 返回数组需要额外的 returnSize 参数
	if sig.ReturnType == TypeIntArray {
		params = append(params, "int* returnSize")
	}
	return fmt.Sprintf("%s %s(%s) {\n    // 请在此实现你的代码\n}", ret, sig.Name, strings.Join(params, ", "))
}

func (g *CGenerator) genDriver(sig *FunctionSignature) (string, error) {
	var sb strings.Builder
	sb.WriteString("#include <stdio.h>\n#include <stdlib.h>\n#include <string.h>\n#include <stdbool.h>\n\n")

	needsListNode := false
	needsTreeNode := false
	for _, p := range sig.Params {
		if p.Type == TypeListNode {
			needsListNode = true
		}
		if p.Type == TypeTreeNode {
			needsTreeNode = true
		}
	}
	if sig.ReturnType == TypeListNode {
		needsListNode = true
	}
	if sig.ReturnType == TypeTreeNode {
		needsTreeNode = true
	}

	if needsListNode {
		sb.WriteString(cListNodeHelper())
		sb.WriteString("\n")
	}
	if needsTreeNode {
		sb.WriteString(cTreeNodeHelper())
		sb.WriteString("\n")
	}

	sb.WriteString("// 用户代码将在此处被包含\n\n")
	sb.WriteString("int main() {\n")

	var callArgs []string
	for _, p := range sig.Params {
		readCode, args, err := cReadParam(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(readCode)
		callArgs = append(callArgs, args...)
	}

	// C 返回数组需要额外参数
	if sig.ReturnType == TypeIntArray {
		sb.WriteString("    int returnSize = 0;\n")
		callArgs = append(callArgs, "&returnSize")
	}

	call := fmt.Sprintf("%s(%s)", sig.Name, strings.Join(callArgs, ", "))
	if sig.ReturnType == "" {
		sb.WriteString(fmt.Sprintf("    %s;\n", call))
	} else {
		retCType := cType(sig.ReturnType)
		sb.WriteString(fmt.Sprintf("    %s result = %s;\n", retCType, call))
		sb.WriteString(cOutputCode(sig.ReturnType))
	}

	sb.WriteString("    return 0;\n}\n")
	return sb.String(), nil
}

func cReadParam(p Param) (string, []string, error) {
	switch p.Type {
	case TypeInt:
		return fmt.Sprintf("    int %s;\n    scanf(\"%%d\", &%s);\n", p.Name, p.Name),
			[]string{p.Name}, nil
	case TypeFloat:
		return fmt.Sprintf("    double %s;\n    scanf(\"%%lf\", &%s);\n", p.Name, p.Name),
			[]string{p.Name}, nil
	case TypeString:
		return fmt.Sprintf("    char %s[100001];\n    scanf(\"%%s\", %s);\n", p.Name, p.Name),
			[]string{p.Name}, nil
	case TypeBool:
		return fmt.Sprintf("    char %s_str[10];\n    scanf(\"%%s\", %s_str);\n    bool %s = (strcmp(%s_str, \"true\") == 0);\n",
			p.Name, p.Name, p.Name, p.Name), []string{p.Name}, nil
	case TypeIntArray:
		sizeVar := p.Name + "Size"
		code := fmt.Sprintf(`    int %s = 0;
    int %s[100001];
    {
        char line[1000001];
        if (fgets(line, sizeof(line), stdin) != NULL) {
            char *tok = strtok(line, " \n");
            while (tok != NULL) {
                %s[%s++] = atoi(tok);
                tok = strtok(NULL, " \n");
            }
        }
    }
`, sizeVar, p.Name, p.Name, sizeVar)
		return code, []string{p.Name, sizeVar}, nil
	case TypeListNode:
		code := fmt.Sprintf(`    struct ListNode* %s;
    {
        char line[100001];
        fgets(line, sizeof(line), stdin);
        int vals[10001];
        int cnt = 0;
        char *tok = strtok(line, " \n");
        while (tok != NULL) {
            if (strcmp(tok, "null") == 0) { %s = NULL; cnt = -1; break; }
            vals[cnt++] = atoi(tok);
            tok = strtok(NULL, " \n");
        }
        if (cnt > 0) %s = buildList(vals, cnt);
        else if (cnt == 0) %s = NULL;
    }
`, p.Name, p.Name, p.Name, p.Name)
		return code, []string{p.Name}, nil
	default:
		return "", nil, fmt.Errorf("unsupported C type: %s", p.Type)
	}
}

func cOutputCode(retType string) string {
	switch retType {
	case TypeInt:
		return "    printf(\"%d\\n\", result);\n"
	case TypeFloat:
		return "    printf(\"%f\\n\", result);\n"
	case TypeBool:
		return "    printf(\"%s\\n\", result ? \"true\" : \"false\");\n"
	case TypeString:
		return "    printf(\"%s\\n\", result);\n"
	case TypeIntArray:
		return `    for (int i = 0; i < returnSize; i++) {
        if (i > 0) printf(" ");
        printf("%d", result[i]);
    }
    printf("\n");
`
	case TypeListNode:
		return "    printList(result);\n"
	default:
		return "    printf(\"%d\\n\", result);\n"
	}
}

func cType(t string) string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "double"
	case TypeString:
		return "char*"
	case TypeBool:
		return "bool"
	case TypeIntArray:
		return "int*"
	case TypeListNode:
		return "struct ListNode*"
	case TypeTreeNode:
		return "struct TreeNode*"
	case "":
		return ""
	default:
		return "int"
	}
}

func cParamDecls(t, name string) []string {
	switch t {
	case TypeInt:
		return []string{fmt.Sprintf("int %s", name)}
	case TypeFloat:
		return []string{fmt.Sprintf("double %s", name)}
	case TypeString:
		return []string{fmt.Sprintf("char* %s", name)}
	case TypeBool:
		return []string{fmt.Sprintf("bool %s", name)}
	case TypeIntArray:
		return []string{fmt.Sprintf("int* %s", name), fmt.Sprintf("int %sSize", name)}
	case TypeListNode:
		return []string{fmt.Sprintf("struct ListNode* %s", name)}
	case TypeTreeNode:
		return []string{fmt.Sprintf("struct TreeNode* %s", name)}
	default:
		return []string{fmt.Sprintf("int %s", name)}
	}
}

func cListNodeHelper() string {
	return `struct ListNode {
    int val;
    struct ListNode *next;
};

struct ListNode* buildList(int* vals, int size) {
    if (size == 0) return NULL;
    struct ListNode* head = (struct ListNode*)malloc(sizeof(struct ListNode));
    head->val = vals[0];
    head->next = NULL;
    struct ListNode* cur = head;
    for (int i = 1; i < size; i++) {
        cur->next = (struct ListNode*)malloc(sizeof(struct ListNode));
        cur->next->val = vals[i];
        cur->next->next = NULL;
        cur = cur->next;
    }
    return head;
}

void printList(struct ListNode* head) {
    if (!head) { printf("null\n"); return; }
    int first = 1;
    while (head) {
        if (!first) printf(" ");
        printf("%d", head->val);
        first = 0;
        head = head->next;
    }
    printf("\n");
}
`
}

func cTreeNodeHelper() string {
	return `struct TreeNode {
    int val;
    struct TreeNode *left;
    struct TreeNode *right;
};

struct TreeNode* newTreeNode(int val) {
    struct TreeNode* node = (struct TreeNode*)malloc(sizeof(struct TreeNode));
    node->val = val;
    node->left = NULL;
    node->right = NULL;
    return node;
}
`
}
