package codegen

import (
	"fmt"
	"strings"
)

type JavaGenerator struct{}

func (g *JavaGenerator) Language() string { return "Java" }

func (g *JavaGenerator) Generate(sig *FunctionSignature) (*GeneratedCode, error) {
	template := g.genTemplate(sig)
	driver, err := g.genDriver(sig)
	if err != nil {
		return nil, err
	}
	return &GeneratedCode{
		Language:         "Java",
		FunctionTemplate: template,
		DriverCode:       driver,
	}, nil
}

func (g *JavaGenerator) genTemplate(sig *FunctionSignature) string {
	params := make([]string, len(sig.Params))
	for i, p := range sig.Params {
		params[i] = javaType(p.Type) + " " + p.Name
	}
	ret := javaType(sig.ReturnType)
	if ret == "" {
		ret = "void"
	}
	return fmt.Sprintf("class Solution {\n    public %s %s(%s) {\n        // 请在此实现你的代码\n    }\n}", ret, sig.Name, strings.Join(params, ", "))
}

func (g *JavaGenerator) genDriver(sig *FunctionSignature) (string, error) {
	var sb strings.Builder

	sb.WriteString("import java.util.*;\nimport java.io.*;\n\n")

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
		sb.WriteString(javaListNodeHelper())
		sb.WriteString("\n")
	}
	if needsTreeNode {
		sb.WriteString(javaTreeNodeHelper())
		sb.WriteString("\n")
	}

	sb.WriteString("// Solution 类由用户代码提供\n\n")
	sb.WriteString("class Main {\n    public static void main(String[] args) throws Exception {\n")
	sb.WriteString("        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));\n")
	sb.WriteString("        Solution sol = new Solution();\n")

	var callArgs []string
	for _, p := range sig.Params {
		readCode, err := javaReadParam(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(readCode)
		callArgs = append(callArgs, p.Name)
	}

	call := fmt.Sprintf("sol.%s(%s)", sig.Name, strings.Join(callArgs, ", "))
	if sig.ReturnType == "" {
		sb.WriteString(fmt.Sprintf("        %s;\n", call))
	} else {
		jType := javaType(sig.ReturnType)
		sb.WriteString(fmt.Sprintf("        %s result = %s;\n", jType, call))
		sb.WriteString(javaOutputCode(sig.ReturnType))
	}

	sb.WriteString("    }\n}\n")
	return sb.String(), nil
}

func javaReadParam(p Param) (string, error) {
	switch p.Type {
	case TypeInt:
		return fmt.Sprintf("        int %s = Integer.parseInt(br.readLine().trim());\n", p.Name), nil
	case TypeFloat:
		return fmt.Sprintf("        double %s = Double.parseDouble(br.readLine().trim());\n", p.Name), nil
	case TypeString:
		return fmt.Sprintf("        String %s = br.readLine().trim();\n", p.Name), nil
	case TypeBool:
		return fmt.Sprintf("        boolean %s = br.readLine().trim().equals(\"true\");\n", p.Name), nil
	case TypeIntArray:
		return fmt.Sprintf(`        int[] %s;
        {
            String[] parts = br.readLine().trim().split("\\s+");
            %s = new int[parts.length];
            for (int i = 0; i < parts.length; i++) {
                %s[i] = Integer.parseInt(parts[i]);
            }
        }
`, p.Name, p.Name, p.Name), nil
	case TypeStrArray:
		return fmt.Sprintf("        String[] %s = br.readLine().trim().split(\"\\\\s+\");\n", p.Name), nil
	case TypeInt2D:
		return fmt.Sprintf(`        int %sRows = Integer.parseInt(br.readLine().trim());
        int[][] %s = new int[%sRows][];
        for (int i = 0; i < %sRows; i++) {
            String[] parts = br.readLine().trim().split("\\s+");
            %s[i] = new int[parts.length];
            for (int j = 0; j < parts.length; j++) {
                %s[i][j] = Integer.parseInt(parts[j]);
            }
        }
`, p.Name, p.Name, p.Name, p.Name, p.Name, p.Name), nil
	case TypeListNode:
		return fmt.Sprintf(`        ListNode %s;
        {
            String line = br.readLine().trim();
            String[] parts = line.split("\\s+");
            %s = buildList(parts);
        }
`, p.Name, p.Name), nil
	case TypeTreeNode:
		return fmt.Sprintf(`        TreeNode %s;
        {
            String line = br.readLine().trim();
            String[] parts = line.split("\\s+");
            %s = buildTree(parts);
        }
`, p.Name, p.Name), nil
	default:
		return "", fmt.Errorf("unsupported Java type: %s", p.Type)
	}
}

func javaOutputCode(retType string) string {
	switch retType {
	case TypeInt, TypeFloat, TypeBool:
		return "        System.out.println(result);\n"
	case TypeString:
		return "        System.out.println(result);\n"
	case TypeIntArray:
		return `        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < result.length; i++) {
            if (i > 0) sb.append(" ");
            sb.append(result[i]);
        }
        System.out.println(sb.toString());
`
	case TypeStrArray:
		return "        System.out.println(String.join(\" \", result));\n"
	case TypeInt2D:
		return `        for (int[] row : result) {
            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < row.length; i++) {
                if (i > 0) sb.append(" ");
                sb.append(row[i]);
            }
            System.out.println(sb.toString());
        }
`
	case TypeListNode:
		return "        printList(result);\n"
	case TypeTreeNode:
		return "        printTree(result);\n"
	default:
		return "        System.out.println(result);\n"
	}
}

func javaType(t string) string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "double"
	case TypeString:
		return "String"
	case TypeBool:
		return "boolean"
	case TypeIntArray:
		return "int[]"
	case TypeStrArray:
		return "String[]"
	case TypeInt2D:
		return "int[][]"
	case TypeListNode:
		return "ListNode"
	case TypeTreeNode:
		return "TreeNode"
	case "":
		return ""
	default:
		return t
	}
}

func javaListNodeHelper() string {
	return `class ListNode {
    int val;
    ListNode next;
    ListNode() {}
    ListNode(int val) { this.val = val; }

    static ListNode buildList(String[] vals) {
        if (vals.length == 0 || vals[0].equals("null")) return null;
        ListNode head = new ListNode(Integer.parseInt(vals[0]));
        ListNode cur = head;
        for (int i = 1; i < vals.length; i++) {
            cur.next = new ListNode(Integer.parseInt(vals[i]));
            cur = cur.next;
        }
        return head;
    }

    static void printList(ListNode head) {
        if (head == null) { System.out.println("null"); return; }
        StringBuilder sb = new StringBuilder();
        boolean first = true;
        while (head != null) {
            if (!first) sb.append(" ");
            sb.append(head.val);
            first = false;
            head = head.next;
        }
        System.out.println(sb.toString());
    }
}

`
}

func javaTreeNodeHelper() string {
	return `class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    TreeNode() {}
    TreeNode(int val) { this.val = val; }

    static TreeNode buildTree(String[] vals) {
        if (vals.length == 0 || vals[0].equals("null")) return null;
        TreeNode root = new TreeNode(Integer.parseInt(vals[0]));
        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);
        int i = 1;
        while (!queue.isEmpty() && i < vals.length) {
            TreeNode node = queue.poll();
            if (i < vals.length && !vals[i].equals("null")) {
                node.left = new TreeNode(Integer.parseInt(vals[i]));
                queue.add(node.left);
            }
            i++;
            if (i < vals.length && !vals[i].equals("null")) {
                node.right = new TreeNode(Integer.parseInt(vals[i]));
                queue.add(node.right);
            }
            i++;
        }
        return root;
    }

    static void printTree(TreeNode root) {
        if (root == null) { System.out.println("null"); return; }
        List<String> vals = new ArrayList<>();
        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);
        while (!queue.isEmpty()) {
            TreeNode node = queue.poll();
            if (node == null) {
                vals.add("null");
            } else {
                vals.add(String.valueOf(node.val));
                queue.add(node.left);
                queue.add(node.right);
            }
        }
        while (!vals.isEmpty() && vals.get(vals.size()-1).equals("null")) {
            vals.remove(vals.size()-1);
        }
        System.out.println(String.join(" ", vals));
    }
}

`
}
