package codegen

import (
	"fmt"
	"strings"
)

type CppGenerator struct{}

func (g *CppGenerator) Language() string { return "C++" }

func (g *CppGenerator) Generate(sig *FunctionSignature) (*GeneratedCode, error) {
	template := g.genTemplate(sig)
	driver, err := g.genDriver(sig)
	if err != nil {
		return nil, err
	}
	return &GeneratedCode{
		Language:         "C++",
		FunctionTemplate: template,
		DriverCode:       driver,
	}, nil
}

func (g *CppGenerator) genTemplate(sig *FunctionSignature) string {
	params := make([]string, len(sig.Params))
	for i, p := range sig.Params {
		params[i] = cppParamDecl(p.Type, p.Name)
	}
	ret := cppType(sig.ReturnType)
	if ret == "" {
		ret = "void"
	}
	return fmt.Sprintf("%s %s(%s) {\n    // 请在此实现你的代码\n}", ret, sig.Name, strings.Join(params, ", "))
}

func (g *CppGenerator) genDriver(sig *FunctionSignature) (string, error) {
	var sb strings.Builder

	sb.WriteString("#include <iostream>\n#include <vector>\n#include <string>\n#include <sstream>\n#include <unordered_map>\n")

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
	if needsTreeNode {
		sb.WriteString("#include <queue>\n")
	}

	sb.WriteString("using namespace std;\n\n")

	if needsListNode {
		sb.WriteString(cppListNodeHelper())
		sb.WriteString("\n")
	}
	if needsTreeNode {
		sb.WriteString(cppTreeNodeHelper())
		sb.WriteString("\n")
	}

	// 前置声明用户函数
	sb.WriteString("// 用户代码将在此处被包含\n\n")

	sb.WriteString("int main() {\n")

	var callArgs []string
	for _, p := range sig.Params {
		readCode, err := cppReadParam(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(readCode)
		callArgs = append(callArgs, p.Name)
	}

	call := fmt.Sprintf("%s(%s)", sig.Name, strings.Join(callArgs, ", "))
	if sig.ReturnType == "" {
		sb.WriteString(fmt.Sprintf("    %s;\n", call))
	} else {
		sb.WriteString(fmt.Sprintf("    auto result = %s;\n", call))
		sb.WriteString(cppOutputCode(sig.ReturnType))
	}

	sb.WriteString("    return 0;\n}\n")
	return sb.String(), nil
}

func cppReadParam(p Param) (string, error) {
	switch p.Type {
	case TypeInt:
		return fmt.Sprintf("    int %s;\n    cin >> %s;\n", p.Name, p.Name), nil
	case TypeFloat:
		return fmt.Sprintf("    double %s;\n    cin >> %s;\n", p.Name, p.Name), nil
	case TypeString:
		return fmt.Sprintf("    string %s;\n    getline(cin, %s);\n", p.Name, p.Name), nil
	case TypeBool:
		return fmt.Sprintf("    string %s_str;\n    cin >> %s_str;\n    bool %s = (%s_str == \"true\");\n", p.Name, p.Name, p.Name, p.Name), nil
	case TypeIntArray:
		return fmt.Sprintf(`    vector<int> %s;
    {
        string line;
        getline(cin, line);
        istringstream iss(line);
        int val;
        while (iss >> val) %s.push_back(val);
    }
`, p.Name, p.Name), nil
	case TypeStrArray:
		return fmt.Sprintf(`    vector<string> %s;
    {
        string line;
        getline(cin, line);
        istringstream iss(line);
        string val;
        while (iss >> val) %s.push_back(val);
    }
`, p.Name, p.Name), nil
	case TypeByteArray:
		return fmt.Sprintf("    string %s_str;\n    getline(cin, %s_str);\n    vector<char> %s(%s_str.begin(), %s_str.end());\n",
			p.Name, p.Name, p.Name, p.Name, p.Name), nil
	case TypeInt2D:
		return fmt.Sprintf(`    int %s_rows;
    cin >> %s_rows;
    cin.ignore();
    vector<vector<int>> %s(%s_rows);
    for (int i = 0; i < %s_rows; i++) {
        string line;
        getline(cin, line);
        istringstream iss(line);
        int val;
        while (iss >> val) %s[i].push_back(val);
    }
`, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name), nil
	case TypeByte2D:
		return fmt.Sprintf(`    int %s_rows;
    cin >> %s_rows;
    cin.ignore();
    vector<vector<char>> %s(%s_rows);
    for (int i = 0; i < %s_rows; i++) {
        string line;
        getline(cin, line);
        for (char c : line) %s[i].push_back(c);
    }
`, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name), nil
	case TypeListNode:
		return fmt.Sprintf(`    ListNode* %s;
    {
        string line;
        getline(cin, line);
        istringstream iss(line);
        vector<string> vals;
        string v;
        while (iss >> v) vals.push_back(v);
        %s = buildList(vals);
    }
`, p.Name, p.Name), nil
	case TypeTreeNode:
		return fmt.Sprintf(`    TreeNode* %s;
    {
        string line;
        getline(cin, line);
        istringstream iss(line);
        vector<string> vals;
        string v;
        while (iss >> v) vals.push_back(v);
        %s = buildTree(vals);
    }
`, p.Name, p.Name), nil
	default:
		return "", fmt.Errorf("unsupported C++ type: %s", p.Type)
	}
}

func cppOutputCode(retType string) string {
	switch retType {
	case TypeInt, TypeFloat, TypeBool:
		return "    cout << result << endl;\n"
	case TypeString:
		return "    cout << result << endl;\n"
	case TypeIntArray:
		return `    for (int i = 0; i < (int)result.size(); i++) {
        if (i > 0) cout << " ";
        cout << result[i];
    }
    cout << endl;
`
	case TypeStrArray:
		return `    for (int i = 0; i < (int)result.size(); i++) {
        if (i > 0) cout << " ";
        cout << result[i];
    }
    cout << endl;
`
	case TypeInt2D:
		return `    for (auto& row : result) {
        for (int i = 0; i < (int)row.size(); i++) {
            if (i > 0) cout << " ";
            cout << row[i];
        }
        cout << endl;
    }
`
	case TypeListNode:
		return "    printList(result);\n"
	case TypeTreeNode:
		return "    printTree(result);\n"
	default:
		return "    cout << result << endl;\n"
	}
}

func cppType(t string) string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "double"
	case TypeString:
		return "string"
	case TypeBool:
		return "bool"
	case TypeIntArray:
		return "vector<int>"
	case TypeStrArray:
		return "vector<string>"
	case TypeByteArray:
		return "vector<char>"
	case TypeInt2D:
		return "vector<vector<int>>"
	case TypeByte2D:
		return "vector<vector<char>>"
	case TypeListNode:
		return "ListNode*"
	case TypeTreeNode:
		return "TreeNode*"
	case "":
		return ""
	default:
		return t
	}
}

func cppParamDecl(t, name string) string {
	ct := cppType(t)
	// 复杂类型用 const 引用
	switch t {
	case TypeIntArray, TypeStrArray, TypeInt2D, TypeByte2D, TypeByteArray, TypeString:
		return fmt.Sprintf("const %s& %s", ct, name)
	default:
		return fmt.Sprintf("%s %s", ct, name)
	}
}

func cppListNodeHelper() string {
	return `struct ListNode {
    int val;
    ListNode *next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
};

ListNode* buildList(const vector<string>& vals) {
    if (vals.empty() || vals[0] == "null") return nullptr;
    ListNode* head = new ListNode(stoi(vals[0]));
    ListNode* cur = head;
    for (int i = 1; i < (int)vals.size(); i++) {
        cur->next = new ListNode(stoi(vals[i]));
        cur = cur->next;
    }
    return head;
}

void printList(ListNode* head) {
    if (!head) { cout << "null" << endl; return; }
    bool first = true;
    while (head) {
        if (!first) cout << " ";
        cout << head->val;
        first = false;
        head = head->next;
    }
    cout << endl;
}
`
}

func cppTreeNodeHelper() string {
	return `struct TreeNode {
    int val;
    TreeNode *left;
    TreeNode *right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

TreeNode* buildTree(const vector<string>& vals) {
    if (vals.empty() || vals[0] == "null") return nullptr;
    TreeNode* root = new TreeNode(stoi(vals[0]));
    queue<TreeNode*> q;
    q.push(root);
    int i = 1;
    while (!q.empty() && i < (int)vals.size()) {
        TreeNode* node = q.front(); q.pop();
        if (i < (int)vals.size() && vals[i] != "null") {
            node->left = new TreeNode(stoi(vals[i]));
            q.push(node->left);
        }
        i++;
        if (i < (int)vals.size() && vals[i] != "null") {
            node->right = new TreeNode(stoi(vals[i]));
            q.push(node->right);
        }
        i++;
    }
    return root;
}

void printTree(TreeNode* root) {
    if (!root) { cout << "null" << endl; return; }
    vector<string> vals;
    queue<TreeNode*> q;
    q.push(root);
    while (!q.empty()) {
        TreeNode* node = q.front(); q.pop();
        if (node) {
            vals.push_back(to_string(node->val));
            q.push(node->left);
            q.push(node->right);
        } else {
            vals.push_back("null");
        }
    }
    while (!vals.empty() && vals.back() == "null") vals.pop_back();
    for (int i = 0; i < (int)vals.size(); i++) {
        if (i > 0) cout << " ";
        cout << vals[i];
    }
    cout << endl;
}
`
}
