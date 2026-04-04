package codegen

import (
	"fmt"
	"strings"
)

type GoGenerator struct{}

func (g *GoGenerator) Language() string { return "Go" }

func (g *GoGenerator) Generate(sig *FunctionSignature) (*GeneratedCode, error) {
	template := g.genTemplate(sig)
	driver, err := g.genDriver(sig)
	if err != nil {
		return nil, err
	}
	return &GeneratedCode{
		Language:         "Go",
		FunctionTemplate: template,
		DriverCode:       driver,
	}, nil
}

func (g *GoGenerator) genTemplate(sig *FunctionSignature) string {
	params := make([]string, len(sig.Params))
	for i, p := range sig.Params {
		params[i] = p.Name + " " + goType(p.Type)
	}
	ret := goType(sig.ReturnType)
	if ret != "" {
		ret = " " + ret
	}
	return fmt.Sprintf("func %s(%s)%s {\n    // 请在此实现你的代码\n}", sig.Name, strings.Join(params, ", "), ret)
}

func (g *GoGenerator) genDriver(sig *FunctionSignature) (string, error) {
	var imports []string
	importSet := map[string]bool{"fmt": true}

	var readLines []string
	var callArgs []string

	needsBufio := false
	needsStrconv := false
	needsStrings := false

	for _, p := range sig.Params {
		switch p.Type {
		case TypeInt:
			readLines = append(readLines, fmt.Sprintf("\tvar %s int\n\tfmt.Scan(&%s)", p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeFloat:
			readLines = append(readLines, fmt.Sprintf("\tvar %s float64\n\tfmt.Scan(&%s)", p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeString:
			needsBufio = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf("\t%sLine, _ := reader.ReadString('\\n')\n\t%s := strings.TrimSpace(%sLine)", p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeBool:
			readLines = append(readLines, fmt.Sprintf("\tvar %sStr string\n\tfmt.Scan(&%sStr)\n\t%s := %sStr == \"true\"", p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeIntArray:
			needsBufio = true
			needsStrconv = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	%sLine, _ := reader.ReadString('\n')
	%sLine = strings.TrimSpace(%sLine)
	%sParts := strings.Fields(%sLine)
	%s := make([]int, len(%sParts))
	for i, p := range %sParts {
		%s[i], _ = strconv.Atoi(p)
	}`, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeStrArray:
			needsBufio = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	%sLine, _ := reader.ReadString('\n')
	%sLine = strings.TrimSpace(%sLine)
	%s := strings.Fields(%sLine)`, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeByteArray:
			needsBufio = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	%sLine, _ := reader.ReadString('\n')
	%sLine = strings.TrimSpace(%sLine)
	%s := []byte(%sLine)`, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeInt2D:
			needsBufio = true
			needsStrconv = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	%sRows, _ := strconv.Atoi(strings.TrimSpace(func() string { s, _ := reader.ReadString('\n'); return s }()))
	%s := make([][]int, %sRows)
	for i := 0; i < %sRows; i++ {
		rowLine, _ := reader.ReadString('\n')
		rowLine = strings.TrimSpace(rowLine)
		rowParts := strings.Fields(rowLine)
		row := make([]int, len(rowParts))
		for j, p := range rowParts {
			row[j], _ = strconv.Atoi(p)
		}
		%s[i] = row
	}`, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeByte2D:
			needsBufio = true
			needsStrconv = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	%sRows, _ := strconv.Atoi(strings.TrimSpace(func() string { s, _ := reader.ReadString('\n'); return s }()))
	%s := make([][]byte, %sRows)
	for i := 0; i < %sRows; i++ {
		rowLine, _ := reader.ReadString('\n')
		rowLine = strings.TrimSpace(rowLine)
		%s[i] = []byte(rowLine)
	}`, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeListNode:
			needsStrconv = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	var %sLine string
	fmt.Scanln(&%sLine)
	%sVals := strings.Fields(%sLine)
	%s := buildList(%sVals)`, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		case TypeTreeNode:
			needsBufio = true
			needsStrconv = true
			needsStrings = true
			readLines = append(readLines, fmt.Sprintf(`	%sLine, _ := reader.ReadString('\n')
	%sLine = strings.TrimSpace(%sLine)
	%sVals := strings.Fields(%sLine)
	%s := buildTree(%sVals)`, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name))
			callArgs = append(callArgs, p.Name)
		default:
			return "", fmt.Errorf("unsupported Go type: %s", p.Type)
		}
	}

	if needsBufio {
		importSet["bufio"] = true
		importSet["os"] = true
	}
	if needsStrconv {
		importSet["strconv"] = true
	}
	if needsStrings {
		importSet["strings"] = true
	}

	for k := range importSet {
		imports = append(imports, fmt.Sprintf("\t\"%s\"", k))
	}

	// 生成输出代码
	outputCode := g.genOutput(sig.ReturnType)

	// 构建完整 driver code
	var sb strings.Builder
	sb.WriteString("package main\n\nimport (\n")
	sb.WriteString(strings.Join(imports, "\n"))
	sb.WriteString("\n)\n\n")

	// 添加辅助结构体和函数
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
		sb.WriteString(goListNodeHelper())
		sb.WriteString("\n")
	}
	if needsTreeNode {
		sb.WriteString(goTreeNodeHelper())
		sb.WriteString("\n")
	}

	sb.WriteString("func main() {\n")
	if needsBufio {
		sb.WriteString("\treader := bufio.NewReader(os.Stdin)\n")
	}
	for _, line := range readLines {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// 调用函数
	call := fmt.Sprintf("%s(%s)", sig.Name, strings.Join(callArgs, ", "))
	if sig.ReturnType == "" {
		sb.WriteString(fmt.Sprintf("\t%s\n", call))
	} else {
		sb.WriteString(fmt.Sprintf("\tresult := %s\n", call))
		sb.WriteString(outputCode)
	}
	sb.WriteString("}\n")

	return sb.String(), nil
}

func (g *GoGenerator) genOutput(retType string) string {
	switch retType {
	case TypeInt, TypeFloat, TypeBool:
		return "\tfmt.Println(result)\n"
	case TypeString:
		return "\tfmt.Println(result)\n"
	case TypeIntArray:
		return `	strs := make([]string, len(result))
	for i, v := range result {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
`
	case TypeStrArray:
		return "\tfmt.Println(strings.Join(result, \" \"))\n"
	case TypeInt2D:
		return `	for _, row := range result {
		strs := make([]string, len(row))
		for i, v := range row {
			strs[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(strs, " "))
	}
`
	case TypeListNode:
		return "\tfmt.Println(printList(result))\n"
	case TypeTreeNode:
		return "\tfmt.Println(printTree(result))\n"
	default:
		return "\tfmt.Println(result)\n"
	}
}

func goType(t string) string {
	switch t {
	case TypeFloat:
		return "float64"
	case TypeListNode:
		return "*ListNode"
	case TypeTreeNode:
		return "*TreeNode"
	case "":
		return ""
	default:
		return t // int, string, bool, []int, []string, [][]int, []byte, [][]byte 已经是 Go 类型
	}
}

func goListNodeHelper() string {
	return `type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(vals []string) *ListNode {
	if len(vals) == 0 || vals[0] == "null" {
		return nil
	}
	head := &ListNode{}
	head.Val, _ = strconv.Atoi(vals[0])
	cur := head
	for i := 1; i < len(vals); i++ {
		cur.Next = &ListNode{}
		cur.Next.Val, _ = strconv.Atoi(vals[i])
		cur = cur.Next
	}
	return head
}

func printList(head *ListNode) string {
	if head == nil {
		return "null"
	}
	vals := []string{}
	for head != nil {
		vals = append(vals, strconv.Itoa(head.Val))
		head = head.Next
	}
	return strings.Join(vals, " ")
}
`
}

func goTreeNodeHelper() string {
	return `type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func buildTree(vals []string) *TreeNode {
	if len(vals) == 0 || vals[0] == "null" {
		return nil
	}
	root := &TreeNode{}
	root.Val, _ = strconv.Atoi(vals[0])
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != "null" {
			node.Left = &TreeNode{}
			node.Left.Val, _ = strconv.Atoi(vals[i])
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != "null" {
			node.Right = &TreeNode{}
			node.Right.Val, _ = strconv.Atoi(vals[i])
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

func printTree(root *TreeNode) string {
	if root == nil {
		return "null"
	}
	vals := []string{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			vals = append(vals, "null")
		} else {
			vals = append(vals, strconv.Itoa(node.Val))
			queue = append(queue, node.Left)
			queue = append(queue, node.Right)
		}
	}
	// 去掉尾部的 null
	for len(vals) > 0 && vals[len(vals)-1] == "null" {
		vals = vals[:len(vals)-1]
	}
	return strings.Join(vals, " ")
}
`
}
