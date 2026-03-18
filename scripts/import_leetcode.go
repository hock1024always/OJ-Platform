//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/your-org/oj-platform/internal/database"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/pkg/config"
)

func main() {
	if err := config.Load("./config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := database.Init(&config.AppConfig.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	// 清空旧数据
	database.DB.Exec("DELETE FROM test_cases")
	database.DB.Exec("DELETE FROM problems")

	type TC struct {
		Input    string
		Output   string
		IsPublic bool
	}
	type Prob struct {
		Title            string
		Description      string
		Difficulty       string
		Tags             string
		FunctionTemplate string
		DriverCode       string
		TestCases        []TC
	}

	problems := []Prob{
		// ==================== 两数之和 ====================
		{
			Title:      "两数之和",
			Difficulty: "Easy",
			Tags:       "数组,哈希表",
			Description: `给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值 target的那两个整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。

你可以按任意顺序返回答案。

示例 1：
输入：nums = [2,7,11,15], target = 9
输出：[0,1]

示例 2：
输入：nums = [3,2,4], target = 6
输出：[1,2]`,
			FunctionTemplate: `func twoSum(nums []int, target int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	parts := strings.Fields(line1)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	target, _ := strconv.Atoi(line2)
	res := twoSum(nums, target)
	fmt.Println(res[0], res[1])
}`,
			TestCases: []TC{
				// 基础用例
				{"2 7 11 15\n9", "0 1", true},
				{"3 2 4\n6", "1 2", true},
				// 边界：最小数组
				{"1 2\n3", "0 1", false},
				// 边界：相同元素
				{"3 3\n6", "0 1", false},
				// 边界：负数
				{"-1 -2 -3 -4\n-7", "2 3", false},
				// 边界：正负混合
				{"-3 4 3 90\n0", "0 2", false},
				// 边界：大数
				{"1000000000 1000000000\n2000000000", "0 1", false},
			},
		},
		// ==================== 爬楼梯 ====================
		{
			Title:      "爬楼梯",
			Difficulty: "Easy",
			Tags:       "记忆化搜索,数学,动态规划",
			Description: `假设你正在爬楼梯。需要 n 阶你才能到达楼顶。

每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶？

示例 1：
输入：n = 2
输出：2

示例 2：
输入：n = 3
输出：3`,
			FunctionTemplate: `func climbStairs(n int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	fmt.Println(climbStairs(n))
}`,
			TestCases: []TC{
				// 基础用例
				{"1", "1", true},
				{"2", "2", true},
				{"3", "3", true},
				// 边界：最小值
				{"1", "1", false},
				// 中等值
				{"5", "8", false},
				{"10", "89", false},
				// 边界：较大值（测试性能）
				{"20", "10946", false},
				{"30", "1346269", false},
				{"45", "1836311903", false},
			},
		},
		// ==================== 最大子数组和 ====================
		{
			Title:      "最大子数组和",
			Difficulty: "Medium",
			Tags:       "数组,分治,动态规划",
			Description: `给你一个整数数组 nums，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

示例 1：
输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
输出：6

示例 2：
输入：nums = [1]
输出：1`,
			FunctionTemplate: `func maxSubArray(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(maxSubArray(nums))
}`,
			TestCases: []TC{
				// 基础用例
				{"-2 1 -3 4 -1 2 1 -5 4", "6", true},
				{"1", "1", true},
				// 边界：全负数
				{"-1 -2 -3 -4", "-1", false},
				// 边界：全正数
				{"1 2 3 4 5", "15", false},
				// 边界：单个元素
				{"-1", "-1", false},
				{"0", "0", false},
				// 边界：正负交替
				{"5 -4 3 -2 1", "5", false},
				{"-1 2 3 -4 5", "6", false},
			},
		},
		// ==================== 买卖股票的最佳时机 ====================
		{
			Title:      "买卖股票的最佳时机",
			Difficulty: "Easy",
			Tags:       "数组,动态规划",
			Description: `给定一个数组 prices，它的第 i 个元素 prices[i] 是一支给定股票第 i 天的价格。

如果你最多只允许完成一笔交易（即买入和卖出一只股票一次），设计一个算法来计算你所能获取的最大利润。

示例 1：
输入：[7,1,5,3,6,4]
输出：5

示例 2：
输入：prices = [7,6,4,3,1]
输出：0`,
			FunctionTemplate: `func maxProfit(prices []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	prices := make([]int, len(parts))
	for i, p := range parts {
		prices[i], _ = strconv.Atoi(p)
	}
	fmt.Println(maxProfit(prices))
}`,
			TestCases: []TC{
				// 基础用例
				{"7 1 5 3 6 4", "5", true},
				{"7 6 4 3 1", "0", true},
				// 边界：最小数组
				{"1 2", "1", false},
				{"2 1", "0", false},
				// 边界：单天无法交易
				{"1", "0", false},
				// 边界：先跌后涨
				{"3 2 1 4 5", "4", false},
				// 边界：波动
				{"2 4 1 5 3 6", "5", false},
			},
		},
		// ==================== 只出现一次的数字 ====================
		{
			Title:      "只出现一次的数字",
			Difficulty: "Easy",
			Tags:       "位运算,数组",
			Description: `给你一个 非空 整数数组 nums，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。

示例 1：
输入：nums = [2,2,1]
输出：1

示例 2：
输入：nums = [4,1,2,1,2]
输出：4`,
			FunctionTemplate: `func singleNumber(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(singleNumber(nums))
}`,
			TestCases: []TC{
				// 基础用例
				{"2 2 1", "1", true},
				{"4 1 2 1 2", "4", true},
				// 边界：最小数组
				{"1", "1", false},
				// 边界：负数
				{"-1 -1 -2", "-2", false},
				// 边界：零
				{"0 1 1", "0", false},
				// 边界：较大数组
				{"1 2 3 4 5 6 7 8 9 10 1 2 3 4 5 6 7 8 9", "10", false},
			},
		},
		// ==================== 多数元素 ====================
		{
			Title:      "多数元素",
			Difficulty: "Easy",
			Tags:       "数组,哈希表,分治,计数",
			Description: `给定一个大小为 n 的数组 nums，返回其中的多数元素。多数元素是指在数组中出现次数大于⌊n/2⌋的元素。

你可以假设数组是非空的，并且给定的数组总是存在多数元素。

示例 1：
输入：nums = [3,2,3]
输出：3

示例 2：
输入：nums = [2,2,1,1,1,2,2]
输出：2`,
			FunctionTemplate: `func majorityElement(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(majorityElement(nums))
}`,
			TestCases: []TC{
				// 基础用例
				{"3 2 3", "3", true},
				{"2 2 1 1 1 2 2", "2", true},
				// 边界：单元素
				{"1", "1", false},
				// 边界：全部相同
				{"5 5 5 5", "5", false},
				// 边界：刚好过半
				{"1 2 1", "1", false},
				{"1 1 2 3", "1", false},
			},
		},
		// ==================== 移动零 ====================
		{
			Title:      "移动零",
			Difficulty: "Easy",
			Tags:       "数组,双指针",
			Description: `给定一个数组 nums，编写一个函数将所有 0 移动到数组的末尾，同时保持非零元素的相对顺序。

请注意，必须在不复制数组的情况下原地对数组进行操作。

示例 1：
输入：nums = [0,1,0,3,12]
输出：[1,3,12,0,0]

示例 2：
输入：nums = [0]
输出：[0]`,
			FunctionTemplate: `func moveZeroes(nums []int) {
    // 请原地修改 nums，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	moveZeroes(nums)
	strs := make([]string, len(nums))
	for i, v := range nums {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				// 基础用例
				{"0 1 0 3 12", "1 3 12 0 0", true},
				{"0", "0", true},
				// 边界：无零
				{"1 2 3", "1 2 3", false},
				// 边界：全零
				{"0 0 0", "0 0 0", false},
				// 边界：零在末尾
				{"1 2 3 0", "1 2 3 0", false},
				// 边界：零在开头
				{"0 0 1", "1 0 0", false},
			},
		},
		// ==================== 合并两个有序数组 ====================
		{
			Title:      "合并两个有序数组",
			Difficulty: "Easy",
			Tags:       "数组,双指针,排序",
			Description: `给你两个按非递减顺序排列的整数数组 nums1 和 nums2，另有两个整数 m 和 n，分别表示 nums1 和 nums2 中的元素数目。

请你合并 nums2 到 nums1 中，使合并后的数组同样按非递减顺序排列。

注意：最终合并后数组不应由函数返回，而是存储在数组 nums1 中。

示例 1：
输入：nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
输出：[1,2,2,3,5,6]`,
			FunctionTemplate: `func merge(nums1 []int, m int, nums2 []int, n int) {
    // 请原地修改 nums1，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	mLine, _ := reader.ReadString('\n')
	mLine = strings.TrimSpace(mLine)
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	nLine, _ := reader.ReadString('\n')
	nLine = strings.TrimSpace(nLine)
	m, _ := strconv.Atoi(mLine)
	n, _ := strconv.Atoi(nLine)
	parts1 := strings.Fields(line1)
	nums1 := make([]int, len(parts1))
	for i, p := range parts1 {
		nums1[i], _ = strconv.Atoi(p)
	}
	parts2 := strings.Fields(line2)
	nums2 := make([]int, len(parts2))
	for i, p := range parts2 {
		nums2[i], _ = strconv.Atoi(p)
	}
	merge(nums1, m, nums2, n)
	strs := make([]string, m+n)
	for i := 0; i < m+n; i++ {
		strs[i] = strconv.Itoa(nums1[i])
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				// 基础用例
				{"1 2 3 0 0 0\n3\n2 5 6\n3", "1 2 2 3 5 6", true},
				{"1\n1\n\n0", "1", true},
				// 边界：nums2 为空
				{"1 2 3\n3\n\n0", "1 2 3", false},
				// 边界：nums1 有效部分为空
				{"0 0 0\n0\n1 2 3\n3", "1 2 3", false},
				// 边界：交错插入
				{"1 3 5 0 0 0\n3\n2 4 6\n3", "1 2 3 4 5 6", false},
			},
		},
		// ==================== 验证回文串 ====================
		{
			Title:      "验证回文串",
			Difficulty: "Easy",
			Tags:       "双指针,字符串",
			Description: `如果在将所有大写字符转换为小写字符、并移除所有非字母数字字符之后，短语正着读和反着读都一样。则可以认为该短语是回文串。

给你一个字符串 s，如果它是回文串，返回 true；否则，返回 false。

示例 1：
输入：s = "A man, a plan, a canal: Panama"
输出：true

示例 2：
输入：s = "race a car"
输出：false`,
			FunctionTemplate: `func isPalindrome(s string) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimRight(line, "\n")
	fmt.Println(isPalindrome(line))
}`,
			TestCases: []TC{
				// 基础用例
				{"A man, a plan, a canal: Panama", "true", true},
				{"race a car", "false", true},
				// 边界：空串
				{" ", "true", false},
				// 边界：单字符
				{"a", "true", false},
				// 边界：纯数字
				{"12321", "true", false},
				{"12345", "false", false},
				// 边界：特殊字符
				{"!@#$%^&*()", "true", false},
			},
		},
		// ==================== 有效的括号 ====================
		{
			Title:      "有效的括号",
			Difficulty: "Easy",
			Tags:       "栈,字符串",
			Description: `给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s，判断字符串是否有效。

有效字符串需满足：
1. 左括号必须用相同类型的右括号闭合。
2. 左括号必须以正确的顺序闭合。

示例 1：
输入：s = "()"
输出：true

示例 2：
输入：s = "()[]{}"
输出：true

示例 3：
输入：s = "(]"
输出：false`,
			FunctionTemplate: `func isValid(s string) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	fmt.Println(isValid(line))
}`,
			TestCases: []TC{
				// 基础用例
				{"()", "true", true},
				{"()[]{}", "true", true},
				{"(]", "false", true},
				// 边界：空串
				{"", "true", false},
				// 边界：单括号
				{"(", "false", false},
				{")", "false", false},
				// 边界：嵌套
				{"{[]}", "true", false},
				{"([)]", "false", false},
				// 边界：多层嵌套
				{"((()))", "true", false},
				{"((())", "false", false},
			},
		},
		// ==================== 最长公共前缀 ====================
		{
			Title:      "最长公共前缀",
			Difficulty: "Easy",
			Tags:       "字符串",
			Description: `编写一个函数来查找字符串数组中的最长公共前缀。

如果不存在公共前缀，返回空字符串 ""。

示例 1：
输入：strs = ["flower","flow","flight"]
输出："fl"

示例 2：
输入：strs = ["dog","racecar","car"]
输出：""`,
			FunctionTemplate: `func longestCommonPrefix(strs []string) string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	strs := strings.Fields(line)
	fmt.Println(longestCommonPrefix(strs))
}`,
			TestCases: []TC{
				// 基础用例
				{"flower flow flight", "fl", true},
				{"dog racecar car", "", true},
				// 边界：单字符串
				{"hello", "hello", false},
				// 边界：完全相同
				{"abc abc abc", "abc", false},
				// 边界：空字符串在其中
				{"abc \"\" def", "", false},
				// 边界：无公共前缀
				{"a b c", "", false},
			},
		},
		// ==================== 二叉树的最大深度 ====================
		{
			Title:      "二叉树的最大深度",
			Difficulty: "Easy",
			Tags:       "树,深度优先搜索,广度优先搜索",
			Description: `给定一个二叉树 root，返回其最大深度。

二叉树的最大深度是指从根节点到最远叶子节点的最长路径上的节点数。

示例 1：
输入：root = [3,9,20,null,null,15,7]
输出：3

示例 2：
输入：root = [1,null,2]
输出：2`,
			FunctionTemplate: `type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func maxDepth(root *TreeNode) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	fmt.Println(maxDepth(root))
}`,
			TestCases: []TC{
				// 基础用例
				{"3 9 20 null null 15 7", "3", true},
				{"1 null 2", "2", true},
				// 边界：空树
				{"null", "0", false},
				// 边界：只有根
				{"1", "1", false},
				// 边界：左斜树
				{"1 2 null 3 null 4", "4", false},
				// 边界：完全二叉树
				{"1 2 3 4 5 6 7", "3", false},
			},
		},
		// ==================== 找到字符串中所有字母异位词 ====================
		{
			Title:      "找到字符串中所有字母异位词",
			Difficulty: "Medium",
			Tags:       "哈希表,字符串,滑动窗口",
			Description: `给定两个字符串 s 和 p，找到 s 中所有 p 的异位词的子串，返回这些子串的起始索引。不考虑答案输出的顺序。

示例 1：
输入：s = "cbaebabacd", p = "abc"
输出：[0,6]

示例 2：
输入：s = "abab", p = "ab"
输出：[0,1,2]`,
			FunctionTemplate: `func findAnagrams(s string, p string) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	p, _ := reader.ReadString('\n')
	p = strings.TrimSpace(p)
	res := findAnagrams(s, p)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				// 基础用例
				{"cbaebabacd\nabc", "0 6", true},
				{"abab\nab", "0 1 2", true},
				// 边界：s 比 p 短
				{"a\nab", "", false},
				// 边界：无匹配
				{"abcdef\nxyz", "", false},
				// 边界：完全匹配
				{"abc\nabc", "0", false},
				// 边界：单字符模式
				{"aaaa\na", "0 1 2 3", false},
			},
		},
		// ==================== 二叉树的中序遍历 ====================
		{
			Title:      "二叉树的中序遍历",
			Difficulty: "Easy",
			Tags:       "栈,树,深度优先搜索",
			Description: `给定一个二叉树的根节点 root，返回它的中序遍历。

示例 1：
输入：root = [1,null,2,3]
输出：[1,3,2]

示例 2：
输入：root = []
输出：[]`,
			FunctionTemplate: `type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func inorderTraversal(root *TreeNode) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	res := inorderTraversal(root)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				// 基础用例
				{"1 null 2 3", "1 3 2", true},
				{"1", "1", true},
				// 边界：空树
				{"null", "", false},
				// 边界：只有左子树
				{"3 2 null 1", "1 2 3", false},
				// 边界：只有右子树
				{"1 null 2 null 3", "1 2 3", false},
			},
		},
		// ==================== 对称二叉树 ====================
		{
			Title:      "对称二叉树",
			Difficulty: "Easy",
			Tags:       "树,深度优先搜索,广度优先搜索",
			Description: `给你一个二叉树的根节点 root，检查它是否轴对称。

示例 1：
输入：root = [1,2,2,3,4,4,3]
输出：true

示例 2：
输入：root = [1,2,2,null,3,null,3]
输出：false`,
			FunctionTemplate: `type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func isSymmetric(root *TreeNode) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	fmt.Println(isSymmetric(root))
}`,
			TestCases: []TC{
				// 基础用例
				{"1 2 2 3 4 4 3", "true", true},
				{"1 2 2 null 3 null 3", "false", true},
				// 边界：空树
				{"null", "true", false},
				// 边界：单节点
				{"1", "true", false},
				// 边界：只有一层
				{"1 2 2", "true", false},
				{"1 2 3", "false", false},
			},
		},
		// ==================== 反转链表 ====================
		{
			Title:      "反转链表",
			Difficulty: "Easy",
			Tags:       "链表,递归",
			Description: `给你单链表的头节点 head，请你反转链表，并返回反转后的链表。

示例 1：
输入：head = [1,2,3,4,5]
输出：[5,4,3,2,1]

示例 2：
输入：head = [1,2]
输出：[2,1]`,
			FunctionTemplate: `func reverseList(head *ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
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
	vals := []string{}
	for head != nil {
		vals = append(vals, strconv.Itoa(head.Val))
		head = head.Next
	}
	return strings.Join(vals, " ")
}

func main() {
	var line string
	fmt.Scanln(&line)
	vals := strings.Fields(line)
	head := buildList(vals)
	reversed := reverseList(head)
	fmt.Println(printList(reversed))
}`,
			TestCases: []TC{
				{"1 2 3 4 5", "5 4 3 2 1", true},
				{"1 2", "2 1", true},
				{"null", "null", false},
				{"1", "1", false},
				{"1 2 3", "3 2 1", false},
			},
		},
		// ==================== 合并两个有序链表 ====================
		{
			Title:      "合并两个有序链表",
			Difficulty: "Easy",
			Tags:       "链表,递归",
			Description: `将两个升序链表合并为一个新的升序链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。

示例 1：
输入：l1 = [1,2,4], l2 = [1,3,4]
输出：[1,1,2,3,4,4]

示例 2：
输入：l1 = [], l2 = []
输出：[]`,
			FunctionTemplate: `func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(vals []string) *ListNode {
	if len(vals) == 0 || vals[0] == "null" || vals[0] == "" {
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

func main() {
	var line1, line2 string
	fmt.Scanln(&line1)
	fmt.Scanln(&line2)
	l1 := buildList(strings.Fields(line1))
	l2 := buildList(strings.Fields(line2))
	merged := mergeTwoLists(l1, l2)
	fmt.Println(printList(merged))
}`,
			TestCases: []TC{
				{"1 2 4\n1 3 4", "1 1 2 3 4 4", true},
				{"\n", "null", true},
				{"\n1", "1", false},
				{"1\n", "1", false},
				{"1 3 5\n2 4 6", "1 2 3 4 5 6", false},
			},
		},
		// ==================== 环形链表 ====================
		{
			Title:      "环形链表",
			Difficulty: "Easy",
			Tags:       "链表,双指针",
			Description: `给你一个链表的头节点 head，判断链表中是否有环。

如果链表中有某个节点，可以通过连续跟踪 next 指针再次到达，则链表中存在环。

示例 1：
输入：head = [3,2,0,-4], pos = 1
输出：true

示例 2：
输入：head = [1,2], pos = 0
输出：true`,
			FunctionTemplate: `func hasCycle(head *ListNode) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildCycleList(vals []string, pos int) *ListNode {
	if len(vals) == 0 || vals[0] == "null" {
		return nil
	}
	nodes := make([]*ListNode, len(vals))
	for i, v := range vals {
		nodes[i] = &ListNode{}
		nodes[i].Val, _ = strconv.Atoi(v)
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].Next = nodes[i+1]
	}
	if pos >= 0 && pos < len(nodes) {
		nodes[len(nodes)-1].Next = nodes[pos]
	}
	return nodes[0]
}

func main() {
	var line string
	var pos int
	fmt.Scanln(&line)
	fmt.Scan(&pos)
	vals := strings.Fields(line)
	head := buildCycleList(vals, pos)
	fmt.Println(hasCycle(head))
}`,
			TestCases: []TC{
				{"3 2 0 -4\n1", "true", true},
				{"1 2\n0", "true", true},
				{"1\n-1", "false", true},
				{"1 2 3 4\n-1", "false", false},
				{"1 2 3 4 5\n2", "true", false},
			},
		},
		// ==================== 无重复字符的最长子串 ====================
		{
			Title:      "无重复字符的最长子串",
			Difficulty: "Medium",
			Tags:       "字符串,滑动窗口,哈希表",
			Description: `给定一个字符串 s，请你找出其中不含有重复字符的最长子串的长度。

示例 1：
输入：s = "abcabcbb"
输出：3

示例 2：
输入：s = "bbbbb"
输出：1`,
			FunctionTemplate: `func lengthOfLongestSubstring(s string) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var s string
	fmt.Scanln(&s)
	fmt.Println(lengthOfLongestSubstring(s))
}`,
			TestCases: []TC{
				{"abcabcbb", "3", true},
				{"bbbbb", "1", true},
				{"pwwkew", "3", true},
				{"", "0", false},
				{"a", "1", false},
				{"abcdef", "6", false},
				{"abba", "2", false},
			},
		},
		// ==================== 盛最多水的容器 ====================
		{
			Title:      "盛最多水的容器",
			Difficulty: "Medium",
			Tags:       "数组,双指针,贪心",
			Description: `给定一个长度为 n 的整数数组 height。有 n 条垂线，第 i 条线的两个端点是 (i, 0) 和 (i, height[i])。

找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。

返回容器可以储存的最大水量。

示例 1：
输入：height = [1,8,6,2,5,4,8,3,7]
输出：49`,
			FunctionTemplate: `func maxArea(height []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	height := make([]int, len(parts))
	for i, p := range parts {
		height[i], _ = strconv.Atoi(p)
	}
	fmt.Println(maxArea(height))
}`,
			TestCases: []TC{
				{"1 8 6 2 5 4 8 3 7", "49", true},
				{"1 1", "1", true},
				{"1 2 1", "2", false},
				{"4 3 2 1 4", "16", false},
				{"1 2 4 3", "4", false},
			},
		},
		// ==================== 三数之和 ====================
		{
			Title:      "三数之和",
			Difficulty: "Medium",
			Tags:       "数组,双指针,排序",
			Description: `给你一个整数数组 nums，判断是否存在三元组 [nums[i], nums[j], nums[k]] 满足 i != j、i != k 且 j != k，同时还满足 nums[i] + nums[j] + nums[k] == 0。

请你返回所有和为 0 且不重复的三元组。

示例 1：
输入：nums = [-1,0,1,2,-1,-4]
输出：[[-1,-1,2],[-1,0,1]]`,
			FunctionTemplate: `func threeSum(nums []int) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	result := threeSum(nums)
	// 按字典序排序输出
	sort.Slice(result, func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return false
	})
	for _, tri := range result {
		fmt.Printf("%d %d %d\n", tri[0], tri[1], tri[2])
	}
}`,
			TestCases: []TC{
				{"-1 0 1 2 -1 -4", "-1 -1 2\n-1 0 1", true},
				{"0 1 1", "", true},
				{"0 0 0", "0 0 0", true},
				{"-2 0 1 1 2", "-2 0 2\n-2 1 1", false},
			},
		},
		// ==================== 翻转二叉树 ====================
		{
			Title:      "翻转二叉树",
			Difficulty: "Easy",
			Tags:       "树,递归,DFS",
			Description: `给你一棵二叉树的根节点 root，翻转这棵二叉树，并返回其根节点。

示例 1：
输入：root = [4,2,7,1,3,6,9]
输出：[4,7,2,9,6,3,1]`,
			FunctionTemplate: `func invertTree(root *TreeNode) *TreeNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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
			continue
		}
		vals = append(vals, strconv.Itoa(node.Val))
		queue = append(queue, node.Left, node.Right)
	}
	// 去掉末尾的 null
	for len(vals) > 0 && vals[len(vals)-1] == "null" {
		vals = vals[:len(vals)-1]
	}
	return strings.Join(vals, " ")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	inverted := invertTree(root)
	fmt.Println(printTree(inverted))
}`,
			TestCases: []TC{
				{"4 2 7 1 3 6 9", "4 7 2 9 6 3 1", true},
				{"2 1 3", "2 3 1", true},
				{"null", "null", false},
				{"1", "1", false},
			},
		},
		// ==================== 二叉树的层序遍历 ====================
		{
			Title:      "二叉树的层序遍历",
			Difficulty: "Medium",
			Tags:       "树,BFS,队列",
			Description: `给你二叉树的根节点 root，返回其节点值的层序遍历。（即逐层地，从左到右访问所有节点）。

示例 1：
输入：root = [3,9,20,null,null,15,7]
输出：[[3],[9,20],[15,7]]`,
			FunctionTemplate: `func levelOrder(root *TreeNode) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	result := levelOrder(root)
	for _, level := range result {
		for i, v := range level {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Println()
	}
}`,
			TestCases: []TC{
				{"3 9 20 null null 15 7", "3\n9 20\n15 7", true},
				{"1", "1", true},
				{"null", "", false},
				{"1 2 3 4 5", "1\n2 3\n4 5", false},
			},
		},
		// ==================== 全排列 ====================
		{
			Title:      "全排列",
			Difficulty: "Medium",
			Tags:       "数组,回溯",
			Description: `给定一个不含重复数字的数组 nums，返回其所有可能的全排列。你可以按任意顺序返回答案。

示例 1：
输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]`,
			FunctionTemplate: `func permute(nums []int) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	result := permute(nums)
	// 按字典序排序输出
	sort.Slice(result, func(i, j int) bool {
		for k := 0; k < len(result[i]); k++ {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return false
	})
	for _, perm := range result {
		parts := make([]string, len(perm))
		for i, v := range perm {
			parts[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(parts, " "))
	}
}`,
			TestCases: []TC{
				{"1 2 3", "1 2 3\n1 3 2\n2 1 3\n2 3 1\n3 1 2\n3 2 1", true},
				{"0 1", "0 1\n1 0", true},
				{"1", "1", false},
			},
		},
		// ==================== 最长回文子串 ====================
		{
			Title:      "最长回文子串",
			Difficulty: "Medium",
			Tags:       "字符串,动态规划",
			Description: `给你一个字符串 s，找到 s 中最长的回文子串。

示例 1：
输入：s = "babad"
输出："bab"

示例 2：
输入：s = "cbbd"
输出："bb"`,
			FunctionTemplate: `func longestPalindrome(s string) string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var s string
	fmt.Scanln(&s)
	fmt.Println(longestPalindrome(s))
}`,
			TestCases: []TC{
				{"babad", "bab", true},
				{"cbbd", "bb", true},
				{"a", "a", true},
				{"ac", "a", false},
				{"racecar", "racecar", false},
			},
		},
		// ==================== 接雨水 ====================
		{
			Title:      "接雨水",
			Difficulty: "Hard",
			Tags:       "数组,双指针,栈,动态规划",
			Description: `给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。

示例 1：
输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
输出：6`,
			FunctionTemplate: `func trap(height []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	height := make([]int, len(parts))
	for i, p := range parts {
		height[i], _ = strconv.Atoi(p)
	}
	fmt.Println(trap(height))
}`,
			TestCases: []TC{
				{"0 1 0 2 1 0 1 3 2 1 2 1", "6", true},
				{"4 2 0 3 2 5", "9", true},
				{"1 0 1", "1", false},
				{"0 0 0 0", "0", false},
			},
		},
		// ==================== 删除链表的倒数第N个节点 ====================
		{
			Title:      "删除链表的倒数第N个节点",
			Difficulty: "Medium",
			Tags:       "链表,双指针",
			Description: `给你一个链表，删除链表的倒数第 n 个节点，并且返回链表的头节点。

示例 1：
输入：head = [1,2,3,4,5], n = 2
输出：[1,2,3,5]`,
			FunctionTemplate: `func removeNthFromEnd(head *ListNode, n int) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
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

func main() {
	var line string
	var n int
	fmt.Scanln(&line)
	fmt.Scan(&n)
	vals := strings.Fields(line)
	head := buildList(vals)
	result := removeNthFromEnd(head, n)
	fmt.Println(printList(result))
}`,
			TestCases: []TC{
				{"1 2 3 4 5\n2", "1 2 3 5", true},
				{"1\n1", "null", true},
				{"1 2\n1", "1", false},
				{"1 2 3\n3", "2 3", false},
			},
		},
		// ==================== 路径总和 ====================
		{
			Title:      "路径总和",
			Difficulty: "Easy",
			Tags:       "树,DFS,递归",
			Description: `给你二叉树的根节点 root 和一个表示目标和的整数 targetSum。判断该树中是否存在根节点到叶子节点的路径，这条路径上所有节点值相加等于目标和 targetSum。

示例 1：
输入：root = [5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum = 22
输出：true`,
			FunctionTemplate: `func hasPathSum(root *TreeNode, targetSum int) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var target int
	fmt.Scan(&target)
	vals := strings.Fields(line)
	root := buildTree(vals)
	fmt.Println(hasPathSum(root, target))
}`,
			TestCases: []TC{
				{"5 4 8 11 null 13 4 7 2 null null null 1\n22", "true", true},
				{"1 2 3\n5", "false", true},
				{"null\n0", "false", false},
				{"1 2\n1", "false", false},
			},
		},
		// ==================== 子集 ====================
		{
			Title:      "子集",
			Difficulty: "Medium",
			Tags:       "数组,回溯,位运算",
			Description: `给你一个整数数组 nums，数组中的元素互不相同。返回该数组所有可能的子集（幂集）。

解集不能包含重复的子集。你可以按任意顺序返回解集。

示例 1：
输入：nums = [1,2,3]
输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]`,
			FunctionTemplate: `func subsets(nums []int) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	result := subsets(nums)
	// 按长度和字典序排序
	sort.Slice(result, func(i, j int) bool {
		if len(result[i]) != len(result[j]) {
			return len(result[i]) < len(result[j])
		}
		for k := 0; k < len(result[i]); k++ {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return false
	})
	for _, subset := range result {
		parts := make([]string, len(subset))
		for i, v := range subset {
			parts[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(parts, " "))
	}
}`,
			TestCases: []TC{
				{"1 2 3", "\n1\n2\n3\n1 2\n1 3\n2 3\n1 2 3", true},
				{"0", "\n0", true},
				{"1 2", "\n1\n2\n1 2", false},
			},
		},
		// ==================== 不同路径 ====================
		{
			Title:      "不同路径",
			Difficulty: "Medium",
			Tags:       "数学,动态规划,组合数学",
			Description: `一个机器人位于一个 m x n 网格的左上角。机器人每次只能向下或者向右移动一步。机器人试图达到网格的右下角。

问总共有多少条不同的路径？

示例 1：
输入：m = 3, n = 7
输出：28`,
			FunctionTemplate: `func uniquePaths(m int, n int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var m, n int
	fmt.Scan(&m, &n)
	fmt.Println(uniquePaths(m, n))
}`,
			TestCases: []TC{
				{"3 7", "28", true},
				{"3 2", "3", true},
				{"1 1", "1", false},
				{"7 3", "28", false},
			},
		},
		// ==================== 跳跃游戏 ====================
		{
			Title:      "跳跃游戏",
			Difficulty: "Medium",
			Tags:       "数组,贪心",
			Description: `给定一个非负整数数组 nums，你最初位于数组的第一个下标。数组中的每个元素代表你在该位置可以跳跃的最大长度。

判断你是否能够到达最后一个下标。

示例 1：
输入：nums = [2,3,1,1,4]
输出：true`,
			FunctionTemplate: `func canJump(nums []int) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(canJump(nums))
}`,
			TestCases: []TC{
				{"2 3 1 1 4", "true", true},
				{"3 2 1 0 4", "false", true},
				{"0", "true", false},
				{"1 0", "true", false},
			},
		},
		// ==================== 岛屿数量 ====================
		{
			Title:      "岛屿数量",
			Difficulty: "Medium",
			Tags:       "数组,DFS,BFS,矩阵",
			Description: `给你一个由 '1'（陆地）和 '0'（水）组成的的二维网格，请你计算网格中岛屿的数量。

岛屿总是被水包围，并且每座岛屿只能由水平方向和/或竖直方向上相邻的陆地连接形成。

示例 1：
输入：grid = [
  ["1","1","1","1","0"],
  ["1","1","0","1","0"],
  ["1","1","0","0","0"],
  ["0","0","0","0","0"]
]
输出：1`,
			FunctionTemplate: `func numIslands(grid [][]byte) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var grid [][]byte
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\n" || line == "" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		row := []byte{}
		for _, c := range line {
			if c == '0' || c == '1' {
				row = append(row, byte(c))
			}
		}
		if len(row) > 0 {
			grid = append(grid, row)
		}
	}
	fmt.Println(numIslands(grid))
}`,
			TestCases: []TC{
				{"11110\n11010\n11000\n00000", "1", true},
				{"11000\n11000\n00100\n00011", "3", true},
				{"0", "0", false},
				{"1", "1", false},
			},
		},
		// ==================== 合并区间 ====================
		{
			Title:      "合并区间",
			Difficulty: "Medium",
			Tags:       "数组,排序",
			Description: `以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi]。

请你合并所有重叠的区间，并返回一个不重叠的区间数组。

示例 1：
输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
输出：[[1,6],[8,10],[15,18]]`,
			FunctionTemplate: `func merge(intervals [][]int) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var intervals [][]int
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\n" || line == "" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			intervals = append(intervals, []int{start, end})
		}
	}
	result := merge(intervals)
	for _, interval := range result {
		fmt.Printf("%d %d\n", interval[0], interval[1])
	}
}`,
			TestCases: []TC{
				{"1 3\n2 6\n8 10\n15 18", "1 6\n8 10\n15 18", true},
				{"1 4\n4 5", "1 5", true},
				{"1 4\n0 4", "0 4", false},
			},
		},
		// ==================== 字母异位词分组 ====================
		{
			Title:      "字母异位词分组",
			Difficulty: "Medium",
			Tags:       "哈希表,字符串,排序",
			Description: `给你一个字符串数组，请你将字母异位词组合在一起。可以按任意顺序返回结果列表。

字母异位词是由重新排列源单词的所有字母得到的一个新单词。

示例 1：
输入：strs = ["eat","tea","tan","ate","nat","bat"]
输出：[["bat"],["nat","tan"],["ate","eat","tea"]]

示例 2：
输入：strs = [""]
输出：[[""]]`,
			FunctionTemplate: `func groupAnagrams(strs []string) [][]string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	strs := strings.Fields(line)
	result := groupAnagrams(strs)
	// 对每组内部排序，然后对组按第一个元素排序
	for _, g := range result {
		sort.Strings(g)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0]
	})
	for _, g := range result {
		fmt.Println(strings.Join(g, " "))
	}
}`,
			TestCases: []TC{
				{"eat tea tan ate nat bat", "bat\nnat tan\nate eat tea", true},
				{"a", "a", true},
				{"abc bca cab xyz", "abc bca cab\nxyz", false},
			},
		},
		// ==================== 最长连续序列 ====================
		{
			Title:      "最长连续序列",
			Difficulty: "Medium",
			Tags:       "哈希表,数组,并查集",
			Description: `给定一个未排序的整数数组 nums，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。

请你设计并实现时间复杂度为 O(n) 的算法解决此问题。

示例 1：
输入：nums = [100,4,200,1,3,2]
输出：4

示例 2：
输入：nums = [0,3,7,2,5,8,4,6,0,1]
输出：9`,
			FunctionTemplate: `func longestConsecutive(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(longestConsecutive(nums))
}`,
			TestCases: []TC{
				{"100 4 200 1 3 2", "4", true},
				{"0 3 7 2 5 8 4 6 0 1", "9", true},
				{"", "0", false},
				{"1", "1", false},
				{"1 2 0 1", "3", false},
			},
		},
		// ==================== 两数相加 ====================
		{
			Title:      "两数相加",
			Difficulty: "Medium",
			Tags:       "链表,数学,递归",
			Description: `给你两个非空的链表，表示两个非负的整数。它们每位数字都是按照逆序方式存储的，并且每个节点只能存储一位数字。

请你将两个数相加，并以相同形式返回一个表示和的链表。

示例 1：
输入：l1 = [2,4,3], l2 = [5,6,4]
输出：[7,0,8]

示例 2：
输入：l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
输出：[8,9,9,9,0,0,0,1]`,
			FunctionTemplate: `func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(s string) *ListNode {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil
	}
	dummy := &ListNode{}
	cur := dummy
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func printList(head *ListNode) string {
	var vals []string
	for head != nil {
		vals = append(vals, strconv.Itoa(head.Val))
		head = head.Next
	}
	return strings.Join(vals, " ")
}

func main() {
	var line1, line2 string
	fmt.Scanln(&line1)
	fmt.Scanln(&line2)
	l1 := buildList(strings.ReplaceAll(line1, ",", " "))
	l2 := buildList(strings.ReplaceAll(line2, ",", " "))
	fmt.Println(printList(addTwoNumbers(l1, l2)))
}`,
			TestCases: []TC{
				{"2 4 3\n5 6 4", "7 0 8", true},
				{"0\n0", "0", true},
				{"9 9 9 9 9 9 9\n9 9 9 9", "8 9 9 9 0 0 0 1", true},
				{"1\n9 9", "0 0 1", false},
				{"5\n5", "0 1", false},
			},
		},
		// ==================== 颜色分类 ====================
		{
			Title:      "颜色分类",
			Difficulty: "Medium",
			Tags:       "数组,双指针,排序",
			Description: `给定一个包含红色、白色和蓝色、共 n 个元素的数组 nums，原地对它们进行排序，使得相同颜色的元素相邻，并按照红色、白色、蓝色顺序排列。

我们使用整数 0、1 和 2 分别表示红色、白色和蓝色。

示例 1：
输入：nums = [2,0,2,1,1,0]
输出：[0,0,1,1,2,2]

示例 2：
输入：nums = [2,0,1]
输出：[0,1,2]`,
			FunctionTemplate: `func sortColors(nums []int) {
    // 请原地修改 nums，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	sortColors(nums)
	out := make([]string, len(nums))
	for i, v := range nums {
		out[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(out, " "))
}`,
			TestCases: []TC{
				{"2 0 2 1 1 0", "0 0 1 1 2 2", true},
				{"2 0 1", "0 1 2", true},
				{"0", "0", false},
				{"1 1 1", "1 1 1", false},
				{"2 2 0 0 1 1", "0 0 1 1 2 2", false},
			},
		},
		// ==================== 数组中的第K个最大元素 ====================
		{
			Title:      "数组中的第K个最大元素",
			Difficulty: "Medium",
			Tags:       "数组,分治,快速排序,堆",
			Description: `给定整数数组 nums 和整数 k，请返回数组中第 k 个最大的元素。

请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。

示例 1：
输入：nums = [3,2,1,5,6,4], k = 2
输出：5

示例 2：
输入：nums = [3,2,3,1,2,4,5,5,6], k = 4
输出：4`,
			FunctionTemplate: `func findKthLargest(nums []int, k int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	var k int
	fmt.Scan(&k)
	fmt.Println(findKthLargest(nums, k))
}`,
			TestCases: []TC{
				{"3 2 1 5 6 4\n2", "5", true},
				{"3 2 3 1 2 4 5 5 6\n4", "4", true},
				{"1\n1", "1", false},
				{"7 6 5 4 3 2 1\n3", "5", false},
				{"5 2 4 1 3\n1", "5", false},
			},
		},
		// ==================== 前K个高频元素 ====================
		{
			Title:      "前K个高频元素",
			Difficulty: "Medium",
			Tags:       "数组,哈希表,堆,桶排序",
			Description: `给你一个整数数组 nums 和一个整数 k，请你返回其中出现频率前 k 高的元素。你可以按任意顺序返回答案。

示例 1：
输入：nums = [1,1,1,2,2,3], k = 2
输出：[1,2]

示例 2：
输入：nums = [1], k = 1
输出：[1]`,
			FunctionTemplate: `func topKFrequent(nums []int, k int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	var k int
	fmt.Scan(&k)
	result := topKFrequent(nums, k)
	sort.Ints(result)
	out := make([]string, len(result))
	for i, v := range result {
		out[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(out, " "))
}`,
			TestCases: []TC{
				{"1 1 1 2 2 3\n2", "1 2", true},
				{"1\n1", "1", true},
				{"4 4 4 3 3 2\n2", "3 4", false},
				{"5 3 1 1 1 3 5 5 5\n2", "1 5", false},
			},
		},
		// ==================== 矩阵置零 ====================
		{
			Title:      "矩阵置零",
			Difficulty: "Medium",
			Tags:       "数组,矩阵",
			Description: `给定一个 m x n 的矩阵，如果一个元素为 0，则将其所在行和列的所有元素都设为 0。请使用原地算法。

示例 1：
输入：matrix = [[1,1,1],[1,0,1],[1,1,1]]
输出：[[1,0,1],[0,0,0],[1,0,1]]`,
			FunctionTemplate: `func setZeroes(matrix [][]int) {
    // 请原地修改 matrix，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var matrix [][]int
	for {
		line, err := reader.ReadString('\n')
		if line == "\n" || line == "" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		row := make([]int, len(parts))
		for i, p := range parts {
			row[i], _ = strconv.Atoi(p)
		}
		matrix = append(matrix, row)
		if err != nil {
			break
		}
	}
	setZeroes(matrix)
	for _, row := range matrix {
		out := make([]string, len(row))
		for i, v := range row {
			out[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(out, " "))
	}
}`,
			TestCases: []TC{
				{"1 1 1\n1 0 1\n1 1 1", "1 0 1\n0 0 0\n1 0 1", true},
				{"0 1 2 0\n3 4 5 2\n1 3 1 5", "0 0 0 0\n0 4 5 0\n0 3 1 0", true},
				{"1 2 3\n4 5 6\n7 8 9", "1 2 3\n4 5 6\n7 8 9", false},
			},
		},
		// ==================== 螺旋矩阵 ====================
		{
			Title:      "螺旋矩阵",
			Difficulty: "Medium",
			Tags:       "数组,矩阵,模拟",
			Description: `给你一个 m 行 n 列的矩阵 matrix，请按照顺时针螺旋顺序，返回矩阵中的所有元素。

示例 1：
输入：matrix = [[1,2,3],[4,5,6],[7,8,9]]
输出：[1,2,3,6,9,8,7,4,5]`,
			FunctionTemplate: `func spiralOrder(matrix [][]int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var matrix [][]int
	for {
		line, err := reader.ReadString('\n')
		if line == "\n" || line == "" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		row := make([]int, len(parts))
		for i, p := range parts {
			row[i], _ = strconv.Atoi(p)
		}
		matrix = append(matrix, row)
		if err != nil {
			break
		}
	}
	result := spiralOrder(matrix)
	out := make([]string, len(result))
	for i, v := range result {
		out[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(out, " "))
}`,
			TestCases: []TC{
				{"1 2 3\n4 5 6\n7 8 9", "1 2 3 6 9 8 7 4 5", true},
				{"1 2 3 4\n5 6 7 8\n9 10 11 12", "1 2 3 4 8 12 11 10 9 5 6 7", true},
				{"1", "1", false},
				{"1 2\n3 4", "1 2 4 3", false},
			},
		},
		// ==================== 旋转图像 ====================
		{
			Title:      "旋转图像",
			Difficulty: "Medium",
			Tags:       "数组,数学,矩阵",
			Description: `给定一个 n × n 的二维矩阵 matrix 表示一个图像。请你将图像顺时针旋转 90 度。

你必须在原地旋转图像，这意味着你需要直接修改输入的二维矩阵。请不要使用另一个矩阵来旋转图像。

示例 1：
输入：matrix = [[1,2,3],[4,5,6],[7,8,9]]
输出：[[7,4,1],[8,5,2],[9,6,3]]`,
			FunctionTemplate: `func rotate(matrix [][]int) {
    // 请原地修改 matrix，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var matrix [][]int
	for {
		line, err := reader.ReadString('\n')
		if line == "\n" || line == "" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		row := make([]int, len(parts))
		for i, p := range parts {
			row[i], _ = strconv.Atoi(p)
		}
		matrix = append(matrix, row)
		if err != nil {
			break
		}
	}
	rotate(matrix)
	for _, row := range matrix {
		out := make([]string, len(row))
		for i, v := range row {
			out[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(out, " "))
	}
}`,
			TestCases: []TC{
				{"1 2 3\n4 5 6\n7 8 9", "7 4 1\n8 5 2\n9 6 3", true},
				{"5 1 9 11\n2 4 8 10\n13 3 6 7\n15 14 12 16", "15 13 2 5\n14 3 4 1\n12 6 8 9\n16 7 10 11", true},
				{"1", "1", false},
			},
		},
		// ==================== 零钱兑换 ====================
		{
			Title:      "零钱兑换",
			Difficulty: "Medium",
			Tags:       "数组,动态规划,广度优先搜索",
			Description: `给你一个整数数组 coins，代表不同面额的硬币；以及一个整数 amount，代表总金额。

计算并返回可以凑成总金额所需的最少的硬币个数。如果没有任何一种硬币组合能组成总金额，返回 -1。

示例 1：
输入：coins = [1,2,5], amount = 11
输出：3

示例 2：
输入：coins = [2], amount = 3
输出：-1`,
			FunctionTemplate: `func coinChange(coins []int, amount int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	coins := make([]int, len(parts))
	for i, p := range parts {
		coins[i], _ = strconv.Atoi(p)
	}
	var amount int
	fmt.Scan(&amount)
	fmt.Println(coinChange(coins, amount))
}`,
			TestCases: []TC{
				{"1 2 5\n11", "3", true},
				{"2\n3", "-1", true},
				{"1\n0", "0", true},
				{"1 5 10 25\n30", "2", false},
				{"2 5 10\n6", "2", false},
				{"3\n7", "-1", false},
			},
		},
		// ==================== 打家劫舍 ====================
		{
			Title:      "打家劫舍",
			Difficulty: "Medium",
			Tags:       "数组,动态规划",
			Description: `你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。

给定一个代表每个房屋存放金额的非负整数数组，计算你不触动警报装置的情况下，一夜之内能够偷窃到的最高金额。

示例 1：
输入：nums = [1,2,3,1]
输出：4

示例 2：
输入：nums = [2,7,9,3,1]
输出：12`,
			FunctionTemplate: `func rob(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(rob(nums))
}`,
			TestCases: []TC{
				{"1 2 3 1", "4", true},
				{"2 7 9 3 1", "12", true},
				{"0", "0", false},
				{"1 1", "1", false},
				{"5 3 4 11 2", "16", false},
			},
		},
		// ==================== 完全平方数 ====================
		{
			Title:      "完全平方数",
			Difficulty: "Medium",
			Tags:       "数学,动态规划,广度优先搜索",
			Description: `给你一个整数 n，返回和为 n 的完全平方数的最少数量。

完全平方数是一个整数，其值等于另一个整数的平方；换句话说，其值等于一个整数自乘的积。

示例 1：
输入：n = 12
输出：3
解释：12 = 4 + 4 + 4

示例 2：
输入：n = 13
输出：2
解释：13 = 4 + 9`,
			FunctionTemplate: `func numSquares(n int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	fmt.Println(numSquares(n))
}`,
			TestCases: []TC{
				{"12", "3", true},
				{"13", "2", true},
				{"1", "1", false},
				{"4", "1", false},
				{"9", "1", false},
				{"2", "2", false},
				{"3", "3", false},
			},
		},
		// ==================== 单词拆分 ====================
		{
			Title:      "单词拆分",
			Difficulty: "Medium",
			Tags:       "字符串,哈希表,动态规划",
			Description: `给你一个字符串 s 和一个字符串列表 wordDict 作为字典。如果可以利用字典中出现的一个或多个单词拼接出 s 则返回 true。

注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。

示例 1：
输入：s = "leetcode", wordDict = ["leet","code"]
输出：true

示例 2：
输入：s = "applepenapple", wordDict = ["apple","pen"]
输出：true`,
			FunctionTemplate: `func wordBreak(s string, wordDict []string) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	dictLine, _ := reader.ReadString('\n')
	dictLine = strings.TrimSpace(dictLine)
	wordDict := strings.Fields(dictLine)
	fmt.Println(wordBreak(s, wordDict))
}`,
			TestCases: []TC{
				{"leetcode\nleet code", "true", true},
				{"applepenapple\napple pen", "true", true},
				{"catsandog\ncats dog sand and cat", "false", true},
				{"a\na", "true", false},
				{"dogs\ndog s", "true", false},
			},
		},
		// ==================== 最长递增子序列 ====================
		{
			Title:      "最长递增子序列",
			Difficulty: "Medium",
			Tags:       "数组,动态规划,二分查找",
			Description: `给你一个整数数组 nums，找到其中最长严格递增子序列的长度。

示例 1：
输入：nums = [10,9,2,5,3,7,101,18]
输出：4

示例 2：
输入：nums = [0,1,0,3,2,3]
输出：4`,
			FunctionTemplate: `func lengthOfLIS(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(lengthOfLIS(nums))
}`,
			TestCases: []TC{
				{"10 9 2 5 3 7 101 18", "4", true},
				{"0 1 0 3 2 3", "4", true},
				{"7 7 7 7 7", "1", true},
				{"1 2 3 4 5", "5", false},
				{"5 4 3 2 1", "1", false},
				{"1 3 2 4 3 5", "4", false},
			},
		},
		// ==================== 乘积最大子数组 ====================
		{
			Title:      "乘积最大子数组",
			Difficulty: "Medium",
			Tags:       "数组,动态规划",
			Description: `给你一个整数数组 nums，请你找出数组中乘积最大的非空连续子数组（该子数组中至少包含一个数字），并返回该子数组所对应的乘积。

示例 1：
输入：nums = [2,3,-2,4]
输出：6

示例 2：
输入：nums = [-2,0,-1]
输出：0`,
			FunctionTemplate: `func maxProduct(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(maxProduct(nums))
}`,
			TestCases: []TC{
				{"2 3 -2 4", "6", true},
				{"-2 0 -1", "0", true},
				{"2", "2", false},
				{"-2", "-2", false},
				{"3 -1 4", "4", false},
				{"-2 3 -4", "24", false},
			},
		},
		// ==================== 验证二叉搜索树 ====================
		{
			Title:      "验证二叉搜索树",
			Difficulty: "Medium",
			Tags:       "树,DFS,二叉搜索树",
			Description: `给你一个二叉树的根节点 root，判断其是否是一个有效的二叉搜索树。

有效 二叉搜索树定义如下：
- 节点的左子树只包含小于当前节点的数。
- 节点的右子树只包含大于当前节点的数。
- 所有左子树和右子树自身必须也是二叉搜索树。

示例 1：
输入：root = [2,1,3]
输出：true

示例 2：
输入：root = [5,1,4,null,null,3,6]
输出：false`,
			FunctionTemplate: `func isValidBST(root *TreeNode) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	root := buildTree(strings.Fields(line))
	fmt.Println(isValidBST(root))
}`,
			TestCases: []TC{
				{"2 1 3", "true", true},
				{"5 1 4 null null 3 6", "false", true},
				{"1", "true", false},
				{"5 4 6 null null 3 7", "false", false},
				{"2 1 3 null null null null", "true", false},
			},
		},
		// ==================== 二叉搜索树中第K小的元素 ====================
		{
			Title:      "二叉搜索树中第K小的元素",
			Difficulty: "Medium",
			Tags:       "树,中序遍历,二叉搜索树",
			Description: `给定一个二叉搜索树的根节点 root，和一个整数 k，请你设计一个算法查找其中第 k 个最小元素（从 1 开始计数）。

示例 1：
输入：root = [3,1,4,null,2], k = 1
输出：1

示例 2：
输入：root = [5,3,6,2,4,null,null,1], k = 3
输出：3`,
			FunctionTemplate: `func kthSmallest(root *TreeNode, k int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	root := buildTree(strings.Fields(line))
	var k int
	fmt.Scan(&k)
	fmt.Println(kthSmallest(root, k))
}`,
			TestCases: []TC{
				{"3 1 4 null 2\n1", "1", true},
				{"5 3 6 2 4 null null 1\n3", "3", true},
				{"1\n1", "1", false},
				{"3 1 4 null 2\n3", "3", false},
			},
		},
		// ==================== 二叉树的右视图 ====================
		{
			Title:      "二叉树的右视图",
			Difficulty: "Medium",
			Tags:       "树,DFS,BFS",
			Description: `给定一个二叉树的根节点 root，想象自己站在它的右侧，按照从顶部到底部的顺序，返回从右侧所能看到的节点值。

示例 1：
输入：root = [1,2,3,null,5,null,4]
输出：[1,3,4]

示例 2：
输入：root = [1,null,3]
输出：[1,3]`,
			FunctionTemplate: `func rightSideView(root *TreeNode) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	root := buildTree(strings.Fields(line))
	result := rightSideView(root)
	out := make([]string, len(result))
	for i, v := range result {
		out[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(out, " "))
}`,
			TestCases: []TC{
				{"1 2 3 null 5 null 4", "1 3 4", true},
				{"1 null 3", "1 3", true},
				{"null", "", false},
				{"1 2 3 4", "1 3 4", false},
			},
		},
		// ==================== 二叉树的最近公共祖先 ====================
		{
			Title:      "二叉树的最近公共祖先",
			Difficulty: "Medium",
			Tags:       "树,DFS,递归",
			Description: `给定一个二叉树, 找到该树中两个指定节点的最近公共祖先。

最近公共祖先的定义为："对于有根树 T 的两个节点 p、q，最近公共祖先表示为一个节点 x，满足 x 是 p、q 的祖先且 x 的深度尽可能大。"

示例 1：
输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
输出：3`,
			FunctionTemplate: `func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func findNode(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == val {
		return root
	}
	if l := findNode(root.Left, val); l != nil {
		return l
	}
	return findNode(root.Right, val)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var pVal, qVal int
	fmt.Scan(&pVal, &qVal)
	root := buildTree(strings.Fields(line))
	p := findNode(root, pVal)
	q := findNode(root, qVal)
	result := lowestCommonAncestor(root, p, q)
	if result == nil {
		fmt.Println(-1)
	} else {
		fmt.Println(result.Val)
	}
}`,
			TestCases: []TC{
				{"3 5 1 6 2 0 8 null null 7 4\n5 1", "3", true},
				{"3 5 1 6 2 0 8 null null 7 4\n5 4", "5", true},
				{"1 2\n1 2", "1", false},
				{"3 5 1\n5 3", "3", false},
			},
		},
		// ==================== 二叉树的直径 ====================
		{
			Title:      "二叉树的直径",
			Difficulty: "Easy",
			Tags:       "树,DFS",
			Description: `给你一棵二叉树的根节点，返回该树的直径。

二叉树的直径是指树中任意两个节点之间最长路径的长度。这条路径可能经过也可能不经过根节点。

两节点之间路径的长度由它们之间边数表示。

示例 1：
输入：root = [1,2,3,4,5]
输出：3`,
			FunctionTemplate: `func diameterOfBinaryTree(root *TreeNode) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	root := buildTree(strings.Fields(line))
	fmt.Println(diameterOfBinaryTree(root))
}`,
			TestCases: []TC{
				{"1 2 3 4 5", "3", true},
				{"1 2", "1", true},
				{"1", "0", false},
				{"1 2 3 4 5 null null 8", "5", false},
			},
		},
		// ==================== 二叉树展开为链表 ====================
		{
			Title:      "二叉树展开为链表",
			Difficulty: "Medium",
			Tags:       "树,DFS,链表",
			Description: `给你二叉树的根结点 root，请你将它展开为一个单链表：

展开后的单链表应该同样使用 TreeNode，其中 right 子指针指向链表中下一个结点，而左子指针始终为 null。
展开后的单链表应该与二叉树先序遍历顺序相同。

示例 1：
输入：root = [1,2,5,3,4,null,6]
输出：[1,null,null,null,2,null,null,null,3,null,null,null,4,null,null,null,5,null,null,null,6]`,
			FunctionTemplate: `func flatten(root *TreeNode) {
    // 请原地修改 root，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	root := buildTree(strings.Fields(line))
	flatten(root)
	result := []string{}
	for root != nil {
		result = append(result, strconv.Itoa(root.Val))
		root = root.Right
	}
	fmt.Println(strings.Join(result, " "))
}`,
			TestCases: []TC{
				{"1 2 5 3 4 null 6", "1 2 3 4 5 6", true},
				{"null", "", true},
				{"1", "1", false},
				{"1 2 null 3", "1 2 3", false},
			},
		},
		// ==================== 从前序与中序遍历序列构造二叉树 ====================
		{
			Title:      "从前序与中序遍历序列构造二叉树",
			Difficulty: "Medium",
			Tags:       "树,DFS,哈希表",
			Description: `给定两个整数数组 preorder 和 inorder，其中 preorder 是二叉树的先序遍历，inorder 是同一棵树的中序遍历，请构造二叉树并返回其根节点。

示例 1：
输入：preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
输出：[3,9,20,null,null,15,7]`,
			FunctionTemplate: `func buildTree(preorder []int, inorder []int) *TreeNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func printTree(root *TreeNode) string {
	if root == nil {
		return ""
	}
	vals := []string{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			vals = append(vals, "null")
			continue
		}
		vals = append(vals, strconv.Itoa(node.Val))
		queue = append(queue, node.Left, node.Right)
	}
	for len(vals) > 0 && vals[len(vals)-1] == "null" {
		vals = vals[:len(vals)-1]
	}
	return strings.Join(vals, " ")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	toInts := func(s string) []int {
		parts := strings.Fields(s)
		nums := make([]int, len(parts))
		for i, p := range parts {
			nums[i], _ = strconv.Atoi(p)
		}
		return nums
	}
	preorder := toInts(line1)
	inorder := toInts(line2)
	root := buildTree(preorder, inorder)
	fmt.Println(printTree(root))
}`,
			TestCases: []TC{
				{"3 9 20 15 7\n9 3 15 20 7", "3 9 20 null null 15 7", true},
				{"1\n1", "1", true},
				{"1 2\n2 1", "1 2", false},
			},
		},
		// ==================== 课程表 ====================
		{
			Title:      "课程表",
			Difficulty: "Medium",
			Tags:       "图,拓扑排序,BFS,DFS",
			Description: `你这个学期必须选修 numCourses 门课程，记为 0 到 numCourses - 1。
在选修某些课程之前需要一些先修课程。给你一个数组 prerequisites，其中 prerequisites[i] = [ai, bi] 表示如果要学习课程 ai 则必须先学习课程 bi。

请你判断是否可能完成所有课程的学习？

示例 1：
输入：numCourses = 2, prerequisites = [[1,0]]
输出：true

示例 2：
输入：numCourses = 2, prerequisites = [[1,0],[0,1]]
输出：false`,
			FunctionTemplate: `func canFinish(numCourses int, prerequisites [][]int) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	numCourses, _ := strconv.Atoi(line1)
	var prereqs [][]int
	for {
		line, err := reader.ReadString('\n')
		if line == "\n" || line == "" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			a, _ := strconv.Atoi(parts[0])
			b, _ := strconv.Atoi(parts[1])
			prereqs = append(prereqs, []int{a, b})
		}
		if err != nil {
			break
		}
	}
	fmt.Println(canFinish(numCourses, prereqs))
}`,
			TestCases: []TC{
				{"2\n1 0", "true", true},
				{"2\n1 0\n0 1", "false", true},
				{"1", "true", false},
				{"3\n1 0\n2 1", "true", false},
				{"4\n1 0\n2 1\n3 2\n0 3", "false", false},
			},
		},
		// ==================== 实现Trie前缀树 ====================
		{
			Title:      "实现Trie前缀树",
			Difficulty: "Medium",
			Tags:       "字典树,设计,哈希表,字符串",
			Description: `Trie（发音类似 "try"）或者说前缀树是一种树形数据结构，用于高效地存储和检索字符串数据集中的键。这一数据结构有相当多的应用情景，例如自动补全和拼写检查。

请你实现 Trie 类：
- Trie() 初始化前缀树对象。
- void insert(String word) 向前缀树中插入字符串 word。
- boolean search(String word) 如果字符串 word 在前缀树中，返回 true；否则，返回 false。
- boolean startsWith(String prefix) 如果之前已经插入的字符串 word 的前缀之一为 prefix，返回 true；否则，返回 false。

示例：
输入：["Trie","insert","search","search","startsWith","insert","search"]
     [[],["apple"],["apple"],["app"],["app"],["app"],["app"]]
输出：[null,null,true,false,true,null,true]`,
			FunctionTemplate: `type Trie struct {
    // 请在此定义数据结构
}

func Constructor() Trie {
    // 请在此实现你的代码
}

func (t *Trie) Insert(word string) {
    // 请在此实现你的代码
}

func (t *Trie) Search(word string) bool {
    // 请在此实现你的代码
}

func (t *Trie) StartsWith(prefix string) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	trie := Constructor()
	reader := bufio.NewReader(os.Stdin)
	results := []string{}
	for {
		line, err := reader.ReadString('\n')
		if line == "" || line == "\n" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		if len(parts) < 1 {
			break
		}
		op := parts[0]
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}
		switch op {
		case "insert":
			trie.Insert(arg)
		case "search":
			if trie.Search(arg) {
				results = append(results, "true")
			} else {
				results = append(results, "false")
			}
		case "startsWith":
			if trie.StartsWith(arg) {
				results = append(results, "true")
			} else {
				results = append(results, "false")
			}
		}
		if err != nil {
			break
		}
	}
	fmt.Println(strings.Join(results, " "))
}`,
			TestCases: []TC{
				{"insert apple\nsearch apple\nsearch app\nstartsWith app\ninsert app\nsearch app", "true false true true", true},
				{"insert hello\nsearch hello\nsearch hell\nstartsWith hell", "true false true", true},
				{"insert a\nsearch a\nstartsWith a\nsearch b", "true true false", false},
			},
		},
		// ==================== 全排列II ====================
		{
			Title:      "全排列II",
			Difficulty: "Medium",
			Tags:       "数组,回溯",
			Description: `给定一个可包含重复数字的序列 nums，按任意顺序返回所有不重复的全排列。

示例 1：
输入：nums = [1,1,2]
输出：[[1,1,2],[1,2,1],[2,1,1]]

示例 2：
输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]`,
			FunctionTemplate: `func permuteUnique(nums []int) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	result := permuteUnique(nums)
	sort.Slice(result, func(i, j int) bool {
		for k := 0; k < len(result[i]) && k < len(result[j]); k++ {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return len(result[i]) < len(result[j])
	})
	for _, perm := range result {
		pp := make([]string, len(perm))
		for i, v := range perm {
			pp[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(pp, " "))
	}
}`,
			TestCases: []TC{
				{"1 1 2", "1 1 2\n1 2 1\n2 1 1", true},
				{"1 2 3", "1 2 3\n1 3 2\n2 1 3\n2 3 1\n3 1 2\n3 2 1", true},
				{"1", "1", false},
			},
		},
		// ==================== 组合总和 ====================
		{
			Title:      "组合总和",
			Difficulty: "Medium",
			Tags:       "数组,回溯",
			Description: `给你一个无重复元素的整数数组 candidates 和一个目标整数 target，找出 candidates 中可以使数字和为目标数 target 的所有不同组合，并以列表形式返回。

candidates 中的同一个数字可以无限制重复被选取。

示例 1：
输入：candidates = [2,3,6,7], target = 7
输出：[[2,2,3],[7]]`,
			FunctionTemplate: `func combinationSum(candidates []int, target int) [][]int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	candidates := make([]int, len(parts))
	for i, p := range parts {
		candidates[i], _ = strconv.Atoi(p)
	}
	var target int
	fmt.Scan(&target)
	result := combinationSum(candidates, target)
	for _, combo := range result {
		sort.Ints(combo)
	}
	sort.Slice(result, func(i, j int) bool {
		for k := 0; k < len(result[i]) && k < len(result[j]); k++ {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return len(result[i]) < len(result[j])
	})
	for _, combo := range result {
		pp := make([]string, len(combo))
		for i, v := range combo {
			pp[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(pp, " "))
	}
}`,
			TestCases: []TC{
				{"2 3 6 7\n7", "2 2 3\n7", true},
				{"2 3 5\n8", "2 2 2 2\n2 3 3\n3 5", true},
				{"2\n1", "", false},
				{"1 2\n3", "1 1 1\n1 2", false},
			},
		},
		// ==================== 电话号码的字母组合 ====================
		{
			Title:      "电话号码的字母组合",
			Difficulty: "Medium",
			Tags:       "哈希表,字符串,回溯",
			Description: `给定一个仅包含数字 2-9 的字符串，返回所有它能表示的字母组合。答案可以按任意顺序返回。

数字到字母的映射与电话按键相同：
2->abc, 3->def, 4->ghi, 5->jkl, 6->mno, 7->pqrs, 8->tuv, 9->wxyz

示例 1：
输入：digits = "23"
输出：["ad","ae","af","bd","be","bf","cd","ce","cf"]`,
			FunctionTemplate: `func letterCombinations(digits string) []string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	var digits string
	fmt.Scanln(&digits)
	result := letterCombinations(digits)
	sort.Strings(result)
	fmt.Println(strings.Join(result, " "))
}`,
			TestCases: []TC{
				{"23", "ad ae af bd be bf cd ce cf", true},
				{"2", "a b c", true},
				{"", "", false},
				{"9", "w x y z", false},
			},
		},
		// ==================== 括号生成 ====================
		{
			Title:      "括号生成",
			Difficulty: "Medium",
			Tags:       "字符串,动态规划,回溯",
			Description: `数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且有效的括号组合。

示例 1：
输入：n = 3
输出：["((()))","(()())","(())()","()(())","()()()"]`,
			FunctionTemplate: `func generateParenthesis(n int) []string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	var n int
	fmt.Scan(&n)
	result := generateParenthesis(n)
	sort.Strings(result)
	fmt.Println(strings.Join(result, " "))
}`,
			TestCases: []TC{
				{"3", "((()))  (()()) (())() ()(()) ()()()", true},
				{"1", "()", true},
				{"2", "(())()(())", false},
			},
		},
		// ==================== 下一个排列 ====================
		{
			Title:      "下一个排列",
			Difficulty: "Medium",
			Tags:       "数组,双指针",
			Description: `整数数组的一个排列就是将其所有成员以序列或线性顺序排列。

下一个排列是指其整数的下一个字典序更大的排列。

必须原地修改，只允许使用额外常数空间。

示例 1：
输入：nums = [1,2,3]
输出：[1,3,2]

示例 2：
输入：nums = [3,2,1]
输出：[1,2,3]`,
			FunctionTemplate: `func nextPermutation(nums []int) {
    // 请原地修改 nums，无需返回值
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	nextPermutation(nums)
	out := make([]string, len(nums))
	for i, v := range nums {
		out[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(out, " "))
}`,
			TestCases: []TC{
				{"1 2 3", "1 3 2", true},
				{"3 2 1", "1 2 3", true},
				{"1 1 5", "1 5 1", true},
				{"1 3 2", "2 1 3", false},
				{"2 3 1", "3 1 2", false},
			},
		},
		// ==================== 寻找重复数 ====================
		{
			Title:      "寻找重复数",
			Difficulty: "Medium",
			Tags:       "数组,双指针,二分查找",
			Description: `给定一个包含 n + 1 个整数的数组 nums，其数字都在 [1, n] 范围内（包括 1 和 n），可知至少存在一个重复的整数。

假设 nums 只有一个重复的整数，返回这个重复的数。

示例 1：
输入：nums = [1,3,4,2,2]
输出：2

示例 2：
输入：nums = [3,1,3,4,2]
输出：3`,
			FunctionTemplate: `func findDuplicate(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(findDuplicate(nums))
}`,
			TestCases: []TC{
				{"1 3 4 2 2", "2", true},
				{"3 1 3 4 2", "3", true},
				{"1 1", "1", false},
				{"2 2 2 2 2", "2", false},
				{"1 2 3 4 5 6 3", "3", false},
			},
		},
		// ==================== 最小覆盖子串 ====================
		{
			Title:      "最小覆盖子串",
			Difficulty: "Hard",
			Tags:       "哈希表,字符串,滑动窗口",
			Description: `给你一个字符串 s、一个字符串 t。返回 s 中涵盖 t 所有字符的最小子串。如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 ""。

示例 1：
输入：s = "ADOBECODEBANC", t = "ABC"
输出："BANC"

示例 2：
输入：s = "a", t = "a"
输出："a"`,
			FunctionTemplate: `func minWindow(s string, t string) string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	t, _ := reader.ReadString('\n')
	t = strings.TrimSpace(t)
	fmt.Println(minWindow(s, t))
}`,
			TestCases: []TC{
				{"ADOBECODEBANC\nABC", "BANC", true},
				{"a\na", "a", true},
				{"a\naa", "", true},
				{"AABC\nABC", "ABC", false},
				{"ab\nb", "b", false},
			},
		},
		// ==================== 柱状图中最大的矩形 ====================
		{
			Title:      "柱状图中最大的矩形",
			Difficulty: "Hard",
			Tags:       "数组,单调栈",
			Description: `给定 n 个非负整数，用来表示柱状图中各个柱子的高度。每个柱子彼此相邻，且宽度为 1。

求在该柱状图中，能够勾勒出来的矩形的最大面积。

示例 1：
输入：heights = [2,1,5,6,2,3]
输出：10

示例 2：
输入：heights = [2,4]
输出：4`,
			FunctionTemplate: `func largestRectangleInHistogram(heights []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	heights := make([]int, len(parts))
	for i, p := range parts {
		heights[i], _ = strconv.Atoi(p)
	}
	fmt.Println(largestRectangleInHistogram(heights))
}`,
			TestCases: []TC{
				{"2 1 5 6 2 3", "10", true},
				{"2 4", "4", true},
				{"1", "1", false},
				{"2 2 2 2", "8", false},
				{"1 2 3 4 5", "9", false},
			},
		},
		// ==================== 最长有效括号 ====================
		{
			Title:      "最长有效括号",
			Difficulty: "Hard",
			Tags:       "字符串,动态规划,栈",
			Description: `给你一个只包含 '(' 和 ')' 的字符串，找出最长有效（格式正确且连续）括号子串的长度。

示例 1：
输入：s = "(()"
输出：2

示例 2：
输入：s = ")()())"
输出：4`,
			FunctionTemplate: `func longestValidParentheses(s string) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var s string
	fmt.Scanln(&s)
	fmt.Println(longestValidParentheses(s))
}`,
			TestCases: []TC{
				{"(()", "2", true},
				{")()())", "4", true},
				{"", "0", true},
				{"()()", "4", false},
				{"(())", "4", false},
				{"((((", "0", false},
			},
		},
		// ==================== 搜索插入位置 ====================
		{
			Title:      "搜索插入位置",
			Difficulty: "Easy",
			Tags:       "数组,二分查找",
			Description: `给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。如果目标值不存在于数组中，返回它将会被按顺序插入的位置。

示例 1：
输入：nums = [1,3,5,6], target = 5
输出：2`,
			FunctionTemplate: `func searchInsert(nums []int, target int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	var target int
	fmt.Scan(&target)
	fmt.Println(searchInsert(nums, target))
}`,
			TestCases: []TC{
				{"1 3 5 6\n5", "2", true},
				{"1 3 5 6\n2", "1", true},
				{"1 3 5 6\n7", "4", true},
				{"1\n0", "0", false},
			},
		},
		// ==================== 排序链表 ====================
		{
			Title:      "排序链表",
			Difficulty: "Medium",
			Tags:       "链表,排序,归并排序",
			Description: `给你链表的头结点 head，请将其按升序排列并返回排序后的链表。

示例 1：
输入：head = [4,2,1,3]
输出：[1,2,3,4]

示例 2：
输入：head = [-1,5,3,4,0]
输出：[-1,0,3,4,5]`,
			FunctionTemplate: `func sortList(head *ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(s string) *ListNode {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil
	}
	dummy := &ListNode{}
	cur := dummy
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func printList(head *ListNode) string {
	var vals []string
	for head != nil {
		vals = append(vals, strconv.Itoa(head.Val))
		head = head.Next
	}
	return strings.Join(vals, " ")
}

func main() {
	var line string
	fmt.Scanln(&line)
	head := buildList(strings.ReplaceAll(line, ",", " "))
	fmt.Println(printList(sortList(head)))
}`,
			TestCases: []TC{
				{"4 2 1 3", "1 2 3 4", true},
				{"-1 5 3 4 0", "-1 0 3 4 5", true},
				{"1", "1", false},
				{"2 1", "1 2", false},
			},
		},
		// ==================== K个一组翻转链表 ====================
		{
			Title:      "K个一组翻转链表",
			Difficulty: "Hard",
			Tags:       "链表,递归",
			Description: `给你链表的头节点 head，每 k 个节点一组进行翻转，请你返回修改后的链表。

k 是一个正整数，它的值小于或等于链表的长度。如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。

示例 1：
输入：head = [1,2,3,4,5], k = 2
输出：[2,1,4,3,5]

示例 2：
输入：head = [1,2,3,4,5], k = 3
输出：[3,2,1,4,5]`,
			FunctionTemplate: `func reverseKGroup(head *ListNode, k int) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(s string) *ListNode {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil
	}
	dummy := &ListNode{}
	cur := dummy
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func printList(head *ListNode) string {
	var vals []string
	for head != nil {
		vals = append(vals, strconv.Itoa(head.Val))
		head = head.Next
	}
	return strings.Join(vals, " ")
}

func main() {
	var line string
	var k int
	fmt.Scanln(&line)
	fmt.Scan(&k)
	head := buildList(strings.ReplaceAll(line, ",", " "))
	fmt.Println(printList(reverseKGroup(head, k)))
}`,
			TestCases: []TC{
				{"1 2 3 4 5\n2", "2 1 4 3 5", true},
				{"1 2 3 4 5\n3", "3 2 1 4 5", true},
				{"1 2 3 4 5\n1", "1 2 3 4 5", false},
				{"1 2 3 4 5\n5", "5 4 3 2 1", false},
			},
		},
		// ==================== 随机链表的复制 ====================
		{
			Title:      "随机链表的复制",
			Difficulty: "Medium",
			Tags:       "链表,哈希表",
			Description: `给你一个长度为 n 的链表，每个节点包含一个额外增加的随机指针 random，该指针可以指向链表中的任何节点或空节点。

构造这个链表的深拷贝。

示例 1：
输入：head = [[7,null],[13,0],[11,4],[10,2],[1,0]]
输出：[[7,null],[13,0],[11,4],[10,2],[1,0]]`,
			FunctionTemplate: `func copyRandomList(head *Node) *Node {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	n := len(parts) / 2
	if n == 0 {
		fmt.Println("null")
		return
	}
	nodes := make([]*Node, n)
	for i := range nodes {
		v, _ := strconv.Atoi(parts[i*2])
		nodes[i] = &Node{Val: v}
	}
	for i := 0; i < n-1; i++ {
		nodes[i].Next = nodes[i+1]
	}
	for i := 0; i < n; i++ {
		rStr := parts[i*2+1]
		if rStr != "null" {
			ri, _ := strconv.Atoi(rStr)
			nodes[i].Random = nodes[ri]
		}
	}
	result := copyRandomList(nodes[0])
	out := []string{}
	for result != nil {
		rIdx := "null"
		if result.Random != nil {
			// find index
			cur := nodes[0]
			for idx := 0; cur != nil; idx++ {
				if cur.Val == result.Random.Val {
					rIdx = strconv.Itoa(idx)
					break
				}
				cur = cur.Next
			}
		}
		out = append(out, fmt.Sprintf("%d %s", result.Val, rIdx))
		result = result.Next
	}
	fmt.Println(strings.Join(out, " | "))
}`,
			TestCases: []TC{
				{"7 null 13 0 11 4 10 2 1 0", "7 null | 13 0 | 11 4 | 10 2 | 1 0", true},
				{"1 1 2 1", "1 1 | 2 1", true},
				{"3 null", "3 null", false},
			},
		},
		// ==================== LRU缓存 ====================
		{
			Title:      "LRU缓存",
			Difficulty: "Medium",
			Tags:       "设计,链表,哈希表",
			Description: `请你设计并实现一个满足 LRU（最近最少使用）缓存约束的数据结构。

实现 LRUCache 类：
- LRUCache(int capacity) 以正整数作为容量 capacity 初始化 LRU 缓存
- int get(int key) 如果关键字存在于缓存中，则返回关键字的值，否则返回 -1
- void put(int key, int value) 如果关键字已经存在，则变更其数据值；如果不存在，则向缓存中插入该组数据。如果超出容量则逐出最久未使用的关键字。

示例 1：
输入：capacity=2, ops=[put(1,1),put(2,2),get(1),put(3,3),get(2),put(4,4),get(1),get(3),get(4)]
输出：[1,-1,1,3,4]`,
			FunctionTemplate: `type LRUCache struct {
    // 请在此定义数据结构
}

func Constructor(capacity int) LRUCache {
    // 请在此实现你的代码
}

func (c *LRUCache) Get(key int) int {
    // 请在此实现你的代码
}

func (c *LRUCache) Put(key int, value int) {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	cap, _ := strconv.Atoi(strings.TrimSpace(line1))
	cache := Constructor(cap)
	results := []string{}
	for {
		line, err := reader.ReadString('\n')
		if line == "" || line == "\n" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "put":
			k, _ := strconv.Atoi(parts[1])
			v, _ := strconv.Atoi(parts[2])
			cache.Put(k, v)
		case "get":
			k, _ := strconv.Atoi(parts[1])
			results = append(results, strconv.Itoa(cache.Get(k)))
		}
		if err != nil {
			break
		}
	}
	fmt.Println(strings.Join(results, " "))
}`,
			TestCases: []TC{
				{"2\nput 1 1\nput 2 2\nget 1\nput 3 3\nget 2\nput 4 4\nget 1\nget 3\nget 4", "1 -1 1 3 4", true},
				{"1\nput 2 1\nget 2\nput 3 2\nget 2\nget 3", "1 -1 2", true},
				{"2\nput 1 1\nget 1\nget 2", "1 -1", false},
			},
		},
		// ==================== 数组转二叉树（二叉树的后序遍历）====================
		{
			Title:      "二叉树的后序遍历",
			Difficulty: "Easy",
			Tags:       "栈,树,DFS",
			Description: `给你一棵二叉树的根节点 root，返回其节点值的后序遍历。

示例 1：
输入：root = [1,null,2,3]
输出：[3,2,1]

示例 2：
输入：root = []
输出：[]`,
			FunctionTemplate: `func postorderTraversal(root *TreeNode) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	res := postorderTraversal(root)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				{"1 null 2 3", "3 2 1", true},
				{"null", "", true},
				{"1", "1", false},
				{"1 2 3 4 5", "4 5 2 3 1", false},
			},
		},
		// ==================== 二叉树的前序遍历 ====================
		{
			Title:      "二叉树的前序遍历",
			Difficulty: "Easy",
			Tags:       "栈,树,DFS",
			Description: `给你二叉树的根节点 root，返回它节点值的前序遍历。

示例 1：
输入：root = [1,null,2,3]
输出：[1,2,3]`,
			FunctionTemplate: `func preorderTraversal(root *TreeNode) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	vals := strings.Fields(line)
	root := buildTree(vals)
	res := preorderTraversal(root)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				{"1 null 2 3", "1 2 3", true},
				{"null", "", true},
				{"1", "1", false},
				{"1 2 3 4 5", "1 2 4 5 3", false},
			},
		},
		// ==================== 相交链表 ====================
		{
			Title:      "相交链表",
			Difficulty: "Easy",
			Tags:       "链表,哈希表,双指针",
			Description: `给你两个单链表的头节点 headA 和 headB，请你找出并返回两个单链表相交的起始节点。如果两个链表不存在相交节点，返回 null。

示例 1：
输入：intersectVal = 8, listA = [4,1,8,4,5], listB = [5,6,1,8,4,5], skipA = 2, skipB = 3
输出：8`,
			FunctionTemplate: `func getIntersectionNode(headA, headB *ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	var lineA, lineB string
	var skipA, skipB int
	fmt.Scanln(&lineA)
	fmt.Scanln(&lineB)
	fmt.Scan(&skipA, &skipB)

	buildList := func(s string) []*ListNode {
		parts := strings.Fields(s)
		nodes := make([]*ListNode, len(parts))
		for i, p := range parts {
			v, _ := strconv.Atoi(p)
			nodes[i] = &ListNode{Val: v}
		}
		for i := 0; i < len(nodes)-1; i++ {
			nodes[i].Next = nodes[i+1]
		}
		return nodes
	}

	nodesA := buildList(lineA)
	nodesB := buildList(lineB)

	// Build shared tail starting at skipA in A and skipB in B
	if skipA < len(nodesA) && skipB < len(nodesB) {
		// link B's node at skipB to A's node at skipA
		nodesB[skipB].Next = nodesA[skipA]
		for i := skipB + 1; i < len(nodesB); i++ {
			nodesB[i] = nodesA[skipA+(i-skipB)]
		}
	}

	var headA, headB *ListNode
	if len(nodesA) > 0 {
		headA = nodesA[0]
	}
	if len(nodesB) > 0 {
		headB = nodesB[0]
	}

	result := getIntersectionNode(headA, headB)
	if result == nil {
		fmt.Println("null")
	} else {
		fmt.Println(result.Val)
	}
}`,
			TestCases: []TC{
				{"4 1 8 4 5\n5 6 1\n2 3", "8", true},
				{"1 9 1 2 4\n3\n3 1", "2", true},
				{"2 6 4\n1 5\n-1 -1", "null", true},
			},
		},
		// ==================== 回文链表 ====================
		{
			Title:      "回文链表",
			Difficulty: "Easy",
			Tags:       "链表,递归,双指针",
			Description: `给你一个单链表的头节点 head，请你判断该链表是否为回文链表。如果是，返回 true；否则，返回 false。

示例 1：
输入：head = [1,2,2,1]
输出：true

示例 2：
输入：head = [1,2]
输出：false`,
			FunctionTemplate: `func isPalindrome(head *ListNode) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(s string) *ListNode {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil
	}
	dummy := &ListNode{}
	cur := dummy
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func main() {
	var line string
	fmt.Scanln(&line)
	head := buildList(strings.ReplaceAll(line, ",", " "))
	fmt.Println(isPalindrome(head))
}`,
			TestCases: []TC{
				{"1 2 2 1", "true", true},
				{"1 2", "false", true},
				{"1", "true", false},
				{"1 2 1", "true", false},
				{"1 2 3 2 1", "true", false},
			},
		},
		// ==================== 用栈实现队列 ====================
		{
			Title:      "用栈实现队列",
			Difficulty: "Easy",
			Tags:       "设计,栈,队列",
			Description: `请你仅使用两个栈实现先入先出队列。

实现 MyQueue 类：
- void push(int x) 将元素 x 推到队列的末尾
- int pop() 从队列的开头移除并返回元素
- int peek() 返回队列开头的元素
- boolean empty() 如果队列为空，返回 true；否则，返回 false

示例 1：
输入：["MyQueue","push","push","peek","pop","empty"]
输出：[null,null,null,1,1,false]`,
			FunctionTemplate: `type MyQueue struct {
    // 请在此定义数据结构
}

func Constructor() MyQueue {
    // 请在此实现你的代码
}

func (q *MyQueue) Push(x int) {
    // 请在此实现你的代码
}

func (q *MyQueue) Pop() int {
    // 请在此实现你的代码
}

func (q *MyQueue) Peek() int {
    // 请在此实现你的代码
}

func (q *MyQueue) Empty() bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	q := Constructor()
	reader := bufio.NewReader(os.Stdin)
	results := []string{}
	for {
		line, err := reader.ReadString('\n')
		if line == "" || line == "\n" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "push":
			v, _ := strconv.Atoi(parts[1])
			q.Push(v)
		case "pop":
			results = append(results, strconv.Itoa(q.Pop()))
		case "peek":
			results = append(results, strconv.Itoa(q.Peek()))
		case "empty":
			if q.Empty() {
				results = append(results, "true")
			} else {
				results = append(results, "false")
			}
		}
		if err != nil {
			break
		}
	}
	fmt.Println(strings.Join(results, " "))
}`,
			TestCases: []TC{
				{"push 1\npush 2\npeek\npop\nempty", "1 1 false", true},
				{"push 1\npop\nempty", "1 true", true},
				{"push 3\npush 5\npop\npeek", "3 5", false},
			},
		},
		// ==================== 最小栈 ====================
		{
			Title:      "最小栈",
			Difficulty: "Medium",
			Tags:       "设计,栈",
			Description: `设计一个支持 push，pop，top 操作，并能在常数时间内检索到最小元素的栈。

实现 MinStack 类:
- void push(int val) 将元素val推入堆栈。
- void pop() 删除堆栈顶部的元素。
- int top() 获取堆栈顶部的元素。
- int getMin() 获取堆栈中的最小元素。

示例 1：
输入：ops = [push(-2),push(0),push(-3),getMin,pop,top,getMin]
输出：[-3,0,-2]`,
			FunctionTemplate: `type MinStack struct {
    // 请在此定义数据结构
}

func Constructor() MinStack {
    // 请在此实现你的代码
}

func (s *MinStack) Push(val int) {
    // 请在此实现你的代码
}

func (s *MinStack) Pop() {
    // 请在此实现你的代码
}

func (s *MinStack) Top() int {
    // 请在此实现你的代码
}

func (s *MinStack) GetMin() int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := Constructor()
	reader := bufio.NewReader(os.Stdin)
	results := []string{}
	for {
		line, err := reader.ReadString('\n')
		if line == "" || line == "\n" {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "push":
			v, _ := strconv.Atoi(parts[1])
			s.Push(v)
		case "pop":
			s.Pop()
		case "top":
			results = append(results, strconv.Itoa(s.Top()))
		case "getMin":
			results = append(results, strconv.Itoa(s.GetMin()))
		}
		if err != nil {
			break
		}
	}
	fmt.Println(strings.Join(results, " "))
}`,
			TestCases: []TC{
				{"push -2\npush 0\npush -3\ngetMin\npop\ntop\ngetMin", "-3 0 -2", true},
				{"push 1\npush 2\ngetMin\npush 0\ngetMin\npop\ngetMin", "1 0 1", true},
				{"push 5\ngetMin\npush 3\ngetMin\npop\ngetMin", "5 3 5", false},
			},
		},
		// ==================== 二叉搜索树的插入操作 ====================
		{
			Title:      "二叉搜索树的插入操作",
			Difficulty: "Medium",
			Tags:       "树,二叉搜索树",
			Description: `给定二叉搜索树（BST）的根节点 root 和要插入树中的值 val，将值插入二叉搜索树。返回插入后二叉搜索树的根节点。输入数据保证新值和原始二叉搜索树中的任意节点值都不同。

示例 1：
输入：root = [4,2,7,1,3], val = 5
输出：[4,2,7,1,3,5]`,
			FunctionTemplate: `func insertIntoBST(root *TreeNode, val int) *TreeNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func inorder(root *TreeNode) []string {
	if root == nil {
		return nil
	}
	res := inorder(root.Left)
	res = append(res, strconv.Itoa(root.Val))
	res = append(res, inorder(root.Right)...)
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var val int
	fmt.Scan(&val)
	root := buildTree(strings.Fields(line))
	result := insertIntoBST(root, val)
	fmt.Println(strings.Join(inorder(result), " "))
}`,
			TestCases: []TC{
				{"4 2 7 1 3\n5", "1 2 3 4 5 7", true},
				{"40 20 60 10 30 50 70\n25", "10 20 25 30 40 50 60 70", true},
				{"null\n5", "5", false},
				{"4 2 7 1 3\n0", "0 1 2 3 4 7", false},
			},
		},
		// ==================== 删除二叉搜索树中的节点 ====================
		{
			Title:      "删除二叉搜索树中的节点",
			Difficulty: "Medium",
			Tags:       "树,二叉搜索树",
			Description: `给定一个二叉搜索树的根节点 root 和一个值 key，删除二叉搜索树中的 key 对应的节点，并保证二叉搜索树的性质不变。返回二叉搜索树（有可能被更新）的根节点的引用。

示例 1：
输入：root = [5,3,6,2,4,null,7], key = 3
输出：[5,4,6,2,null,null,7]`,
			FunctionTemplate: `func deleteNode(root *TreeNode, key int) *TreeNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func inorder(root *TreeNode) []string {
	if root == nil {
		return nil
	}
	res := inorder(root.Left)
	res = append(res, strconv.Itoa(root.Val))
	res = append(res, inorder(root.Right)...)
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	var key int
	fmt.Scan(&key)
	root := buildTree(strings.Fields(line))
	result := deleteNode(root, key)
	io := inorder(result)
	fmt.Println(strings.Join(io, " "))
}`,
			TestCases: []TC{
				{"5 3 6 2 4 null 7\n3", "2 4 5 6 7", true},
				{"5 3 6 2 4 null 7\n0", "2 3 4 5 6 7", true},
				{"5 3 6 2 4 null 7\n5", "2 3 4 6 7", false},
				{"null\n0", "", false},
			},
		},
		// ==================== 二叉树中的最大路径和 ====================
		{
			Title:      "二叉树中的最大路径和",
			Difficulty: "Hard",
			Tags:       "树,DFS,动态规划",
			Description: `二叉树中的路径被定义为一条节点序列，序列中每对相邻节点之间都存在一条边。同一个节点在一条路径序列中至多出现一次。该路径至少包含一个节点，且不一定经过根节点。

路径和是路径中各节点值的总和。

给你一个二叉树的根节点 root，返回其最大路径和。

示例 1：
输入：root = [1,2,3]
输出：6

示例 2：
输入：root = [-3]
输出：-3`,
			FunctionTemplate: `func maxPathSum(root *TreeNode) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TreeNode struct {
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	root := buildTree(strings.Fields(line))
	fmt.Println(maxPathSum(root))
}`,
			TestCases: []TC{
				{"1 2 3", "6", true},
				{"-3", "-3", true},
				{"-10 9 20 null null 15 7", "42", true},
				{"1 -2 3", "4", false},
			},
		},
		// ==================== 整数转罗马数字 ====================
		{
			Title:      "整数转罗马数字",
			Difficulty: "Medium",
			Tags:       "哈希表,数学,字符串",
			Description: `罗马数字包含以下七种字符：I，V，X，L，C，D 和 M。

给你一个整数，将其转为罗马数字。

示例 1：
输入：num = 3
输出："III"

示例 2：
输入：num = 58
输出："LVIII"

示例 3：
输入：num = 1994
输出："MCMXCIV"`,
			FunctionTemplate: `func intToRoman(num int) string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var num int
	fmt.Scan(&num)
	fmt.Println(intToRoman(num))
}`,
			TestCases: []TC{
				{"3", "III", true},
				{"58", "LVIII", true},
				{"1994", "MCMXCIV", true},
				{"4", "IV", false},
				{"9", "IX", false},
				{"1", "I", false},
				{"3999", "MMMCMXCIX", false},
			},
		},
		// ==================== 罗马数字转整数 ====================
		{
			Title:      "罗马数字转整数",
			Difficulty: "Easy",
			Tags:       "哈希表,数学,字符串",
			Description: `给定一个罗马数字，将其转换成整数。

示例 1：
输入：s = "III"
输出：3

示例 2：
输入：s = "MCMXCIV"
输出：1994`,
			FunctionTemplate: `func romanToInt(s string) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var s string
	fmt.Scanln(&s)
	fmt.Println(romanToInt(s))
}`,
			TestCases: []TC{
				{"III", "3", true},
				{"LVIII", "58", true},
				{"MCMXCIV", "1994", true},
				{"IV", "4", false},
				{"IX", "9", false},
				{"XL", "40", false},
				{"XC", "90", false},
			},
		},
		// ==================== 编辑距离 ====================
		{
			Title:      "编辑距离",
			Difficulty: "Medium",
			Tags:       "字符串,动态规划",
			Description: `给你两个单词 word1 和 word2，请返回将 word1 转换成 word2 所使用的最少操作数（插入、删除、替换）。

示例 1：
输入：word1 = "horse", word2 = "ros"
输出：3

示例 2：
输入：word1 = "intention", word2 = "execution"
输出：5`,
			FunctionTemplate: `func minDistance(word1 string, word2 string) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	w1, _ := reader.ReadString('\n')
	w1 = strings.TrimSpace(w1)
	w2, _ := reader.ReadString('\n')
	w2 = strings.TrimSpace(w2)
	fmt.Println(minDistance(w1, w2))
}`,
			TestCases: []TC{
				{"horse\nros", "3", true},
				{"intention\nexecution", "5", true},
				{"\n", "0", false},
				{"a\n", "1", false},
				{"abc\nabc", "0", false},
				{"abc\nabd", "1", false},
			},
		},
		// ==================== 不同的二叉搜索树 ====================
		{
			Title:      "不同的二叉搜索树",
			Difficulty: "Medium",
			Tags:       "数学,动态规划,树",
			Description: `给你一个整数 n，求恰由 n 个节点组成且节点值从 1 到 n 互不相同的二叉搜索树有多少种？返回满足题意的二叉搜索树的种数。

示例 1：
输入：n = 3
输出：5

示例 2：
输入：n = 1
输出：1`,
			FunctionTemplate: `func numTrees(n int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	fmt.Println(numTrees(n))
}`,
			TestCases: []TC{
				{"3", "5", true},
				{"1", "1", true},
				{"2", "2", false},
				{"4", "14", false},
				{"5", "42", false},
				{"6", "132", false},
			},
		},
		// ==================== 分隔链表 ====================
		{
			Title:      "分隔链表",
			Difficulty: "Medium",
			Tags:       "链表,双指针",
			Description: `给你一个链表的头节点 head 和一个特定值 x，请你对链表进行分隔，使得所有小于 x 的节点都出现在大于或等于 x 的节点之前。

你应当保留两个分区中每个节点的初始相对位置。

示例 1：
输入：head = [1,4,3,2,5,2], x = 3
输出：[1,2,2,4,3,5]`,
			FunctionTemplate: `func partition(head *ListNode, x int) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(s string) *ListNode {
	parts := strings.Fields(s)
	dummy := &ListNode{}
	cur := dummy
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func printList(head *ListNode) string {
	var vals []string
	for head != nil {
		vals = append(vals, strconv.Itoa(head.Val))
		head = head.Next
	}
	return strings.Join(vals, " ")
}

func main() {
	var line string
	var x int
	fmt.Scanln(&line)
	fmt.Scan(&x)
	head := buildList(strings.ReplaceAll(line, ",", " "))
	fmt.Println(printList(partition(head, x)))
}`,
			TestCases: []TC{
				{"1 4 3 2 5 2\n3", "1 2 2 4 3 5", true},
				{"2 1\n2", "1 2", true},
				{"1\n0", "1", false},
				{"3 1 2\n2", "1 3 2", false},
			},
		},
		// ==================== 缺失的第一个正数 ====================
		{
			Title:      "缺失的第一个正数",
			Difficulty: "Hard",
			Tags:       "数组,哈希表",
			Description: `给你一个未排序的整数数组 nums，请你找出其中没有出现的最小的正整数。

请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。

示例 1：
输入：nums = [1,2,0]
输出：3

示例 2：
输入：nums = [3,4,-1,1]
输出：2`,
			FunctionTemplate: `func firstMissingPositive(nums []int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	fmt.Println(firstMissingPositive(nums))
}`,
			TestCases: []TC{
				{"1 2 0", "3", true},
				{"3 4 -1 1", "2", true},
				{"7 8 9 11 12", "1", true},
				{"1 2 3", "4", false},
				{"2 3 4", "1", false},
			},
		},
		// ==================== 第五批：最后14道Hot100 ====================
		// ==================== 和为K的子数组 ====================
		{
			Title:      "和为K的子数组",
			Difficulty: "Medium",
			Tags:       "数组,哈希表,前缀和",
			Description: `给你一个整数数组 nums 和一个整数 k ，请你统计并返回该数组中和为 k 的子数组的个数。
子数组是数组中元素的连续非空序列。

示例 1：
输入：nums = [1,1,1], k = 2
输出：2

示例 2：
输入：nums = [1,2,3], k = 3
输出：2`,
			FunctionTemplate: `func subarraySum(nums []int, k int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	parts := strings.Fields(line1)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	k, _ := strconv.Atoi(line2)
	fmt.Println(subarraySum(nums, k))
}`,
			TestCases: []TC{
				{"1 1 1\n2", "2", true},
				{"1 2 3\n3", "2", true},
				{"1\n0", "0", true},
				{"1 -1 0\n0", "3", false},
				{"-1 -1 1\n0", "1", false},
			},
		},
		// ==================== 滑动窗口最大值 ====================
		{
			Title:      "滑动窗口最大值",
			Difficulty: "Hard",
			Tags:       "数组,队列,滑动窗口,单调队列,堆",
			Description: `给你一个整数数组 nums，有一个大小为 k 的滑动窗口从数组的最左侧移动到数组的最右侧。你只可以看到在滑动窗口内的 k 个数字。滑动窗口每次只向右移动一位。返回滑动窗口中的最大值。

示例 1：
输入：nums = [1,3,-1,-3,5,3,6,7], k = 3
输出：[3,3,5,5,6,7]

示例 2：
输入：nums = [1], k = 1
输出：[1]`,
			FunctionTemplate: `func maxSlidingWindow(nums []int, k int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	parts := strings.Fields(line1)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	k, _ := strconv.Atoi(line2)
	res := maxSlidingWindow(nums, k)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				{"1 3 -1 -3 5 3 6 7\n3", "3 3 5 5 6 7", true},
				{"1\n1", "1", true},
				{"1 -1\n1", "1 -1", true},
				{"9 11\n2", "11", false},
				{"4 -2\n2", "4", false},
			},
		},
		// ==================== 轮转数组 ====================
		{
			Title:      "轮转数组",
			Difficulty: "Medium",
			Tags:       "数组,数学,双指针",
			Description: `给定一个整数数组 nums，将数组中的元素向右轮转 k 个位置，其中 k 是非负数。

示例 1：
输入：nums = [1,2,3,4,5,6,7], k = 3
输出：[5,6,7,1,2,3,4]

示例 2：
输入：nums = [-1,-100,3,99], k = 2
输出：[3,99,-1,-100]`,
			FunctionTemplate: `func rotate(nums []int, k int) {
    // 请在此实现你的代码（原地修改）
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	parts := strings.Fields(line1)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	k, _ := strconv.Atoi(line2)
	rotate(nums, k)
	strs := make([]string, len(nums))
	for i, v := range nums {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				{"1 2 3 4 5 6 7\n3", "5 6 7 1 2 3 4", true},
				{"-1 -100 3 99\n2", "3 99 -1 -100", true},
				{"1 2\n3", "2 1", true},
				{"1\n0", "1", false},
				{"1 2 3\n6", "1 2 3", false},
			},
		},
		// ==================== 除自身以外数组的乘积 ====================
		{
			Title:      "除自身以外数组的乘积",
			Difficulty: "Medium",
			Tags:       "数组,前缀和",
			Description: `给你一个整数数组 nums，返回数组 answer ，其中 answer[i] 等于 nums 中除 nums[i] 之外其余各元素的乘积 。
题目数据保证数组 nums 之中任意元素的全部前缀元素和后缀的乘积都在 32 位整数范围内。请不要使用除法。

示例 1：
输入：nums = [1,2,3,4]
输出：[24,12,8,6]

示例 2：
输入：nums = [-1,1,0,-3,3]
输出：[0,0,9,0,0]`,
			FunctionTemplate: `func productExceptSelf(nums []int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	res := productExceptSelf(nums)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				{"1 2 3 4", "24 12 8 6", true},
				{"-1 1 0 -3 3", "0 0 9 0 0", true},
				{"2 3", "3 2", true},
				{"1 1 1 1", "1 1 1 1", false},
				{"0 0", "0 0", false},
			},
		},
		// ==================== 搜索二维矩阵II ====================
		{
			Title:      "搜索二维矩阵II",
			Difficulty: "Medium",
			Tags:       "数组,二分查找,分治,矩阵",
			Description: `编写一个高效的算法来搜索 m x n 矩阵 matrix 中的一个目标值 target 。该矩阵具有以下特性：
每行的元素从左到右升序排列。
每列的元素从上到下升序排列。

示例 1：
输入：matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 5
输出：true

示例 2：
输入：matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 20
输出：false`,
			FunctionTemplate: `func searchMatrix(matrix [][]int, target int) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	dims := strings.Fields(line1)
	m, _ := strconv.Atoi(dims[0])
	n, _ := strconv.Atoi(dims[1])
	matrix := make([][]int, m)
	for i := 0; i < m; i++ {
		row, _ := reader.ReadString('\n')
		row = strings.TrimSpace(row)
		cols := strings.Fields(row)
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j], _ = strconv.Atoi(cols[j])
		}
	}
	line3, _ := reader.ReadString('\n')
	line3 = strings.TrimSpace(line3)
	target, _ := strconv.Atoi(line3)
	fmt.Println(searchMatrix(matrix, target))
}`,
			TestCases: []TC{
				{"5 5\n1 4 7 11 15\n2 5 8 12 19\n3 6 9 16 22\n10 13 14 17 24\n18 21 23 26 30\n5", "true", true},
				{"5 5\n1 4 7 11 15\n2 5 8 12 19\n3 6 9 16 22\n10 13 14 17 24\n18 21 23 26 30\n20", "false", true},
				{"1 1\n1\n1", "true", true},
				{"2 2\n1 3\n2 4\n4", "true", false},
				{"2 2\n1 3\n2 4\n5", "false", false},
			},
		},
		// ==================== 环形链表II ====================
		{
			Title:      "环形链表II",
			Difficulty: "Medium",
			Tags:       "哈希表,链表,双指针",
			Description: `给定一个链表的头节点 head ，返回链表开始入环的第一个节点。如果链表无环，则返回 null。
不允许修改链表。

示例 1：
输入：head = [3,2,0,-4], pos = 1
输出：返回索引为 1 的链表节点（值为2）

示例 2：
输入：head = [1,2], pos = 0
输出：返回索引为 0 的链表节点（值为1）

示例 3：
输入：head = [1], pos = -1
输出：null（无环）

输入格式：第一行为链表节点值（空格分隔），第二行为 pos（入环位置，-1表示无环）。
输出：入环节点的值，无环输出 null。`,
			FunctionTemplate: `type ListNode struct {
    Val  int
    Next *ListNode
}

func detectCycle(head *ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	parts := strings.Fields(line1)
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	pos, _ := strconv.Atoi(line2)

	if len(parts) == 0 {
		fmt.Println("null")
		return
	}
	nodes := make([]*ListNode, len(parts))
	for i, p := range parts {
		v, _ := strconv.Atoi(p)
		nodes[i] = &ListNode{Val: v}
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].Next = nodes[i+1]
	}
	if pos >= 0 && pos < len(nodes) {
		nodes[len(nodes)-1].Next = nodes[pos]
	}
	result := detectCycle(nodes[0])
	if result == nil {
		fmt.Println("null")
	} else {
		fmt.Println(result.Val)
	}
}`,
			TestCases: []TC{
				{"3 2 0 -4\n1", "2", true},
				{"1 2\n0", "1", true},
				{"1\n-1", "null", true},
				{"1 2 3 4 5\n2", "3", false},
				{"1 2 3\n-1", "null", false},
			},
		},
		// ==================== 合并K个升序链表 ====================
		{
			Title:      "合并K个升序链表",
			Difficulty: "Hard",
			Tags:       "链表,分治,堆,归并排序",
			Description: `给你一个链表数组，每个链表都已经按升序排列。请你将所有链表合并到一个升序链表中，返回合并后的链表。

示例 1：
输入：lists = [[1,4,5],[1,3,4],[2,6]]
输出：[1,1,2,3,4,4,5,6]

示例 2：
输入：lists = []
输出：[]

输入格式：第一行为链表个数k，接下来k行每行为一个链表的节点值（空格分隔，空行表示空链表）。
输出：合并后的链表节点值，空格分隔。`,
			FunctionTemplate: `type ListNode struct {
    Val  int
    Next *ListNode
}

func mergeKLists(lists []*ListNode) *ListNode {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	k, _ := strconv.Atoi(line1)
	lists := make([]*ListNode, k)
	for i := 0; i < k; i++ {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			lists[i] = nil
			continue
		}
		parts := strings.Fields(line)
		dummy := &ListNode{}
		cur := dummy
		for _, p := range parts {
			v, _ := strconv.Atoi(p)
			cur.Next = &ListNode{Val: v}
			cur = cur.Next
		}
		lists[i] = dummy.Next
	}
	head := mergeKLists(lists)
	strs := []string{}
	for head != nil {
		strs = append(strs, strconv.Itoa(head.Val))
		head = head.Next
	}
	if len(strs) == 0 {
		fmt.Println("empty")
	} else {
		fmt.Println(strings.Join(strs, " "))
	}
}`,
			TestCases: []TC{
				{"3\n1 4 5\n1 3 4\n2 6", "1 1 2 3 4 4 5 6", true},
				{"1\n\n", "empty", true},
				{"0\n", "empty", true},
				{"2\n1 2 3\n4 5 6", "1 2 3 4 5 6", false},
				{"3\n5\n3\n1", "1 3 5", false},
			},
		},
		// ==================== 腐烂的橘子 ====================
		{
			Title:      "腐烂的橘子",
			Difficulty: "Medium",
			Tags:       "广度优先搜索,数组,矩阵",
			Description: `在给定的 m x n 网格 grid 中，每个单元格可以有以下三个值之一：
值 0 代表空单元格；
值 1 代表新鲜橘子；
值 2 代表腐烂的橘子。
每分钟，腐烂的橘子周围 4 个方向上相邻的新鲜橘子都会腐烂。返回直到单元格中没有新鲜橘子为止所必须经过的最小分钟数。如果不可能，返回 -1。

示例 1：
输入：grid = [[2,1,1],[1,1,0],[0,1,1]]
输出：4

示例 2：
输入：grid = [[2,1,1],[0,1,1],[1,0,1]]
输出：-1`,
			FunctionTemplate: `func orangesRotting(grid [][]int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	dims := strings.Fields(line1)
	m, _ := strconv.Atoi(dims[0])
	n, _ := strconv.Atoi(dims[1])
	grid := make([][]int, m)
	for i := 0; i < m; i++ {
		row, _ := reader.ReadString('\n')
		row = strings.TrimSpace(row)
		cols := strings.Fields(row)
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			grid[i][j], _ = strconv.Atoi(cols[j])
		}
	}
	fmt.Println(orangesRotting(grid))
}`,
			TestCases: []TC{
				{"3 3\n2 1 1\n1 1 0\n0 1 1", "4", true},
				{"3 3\n2 1 1\n0 1 1\n1 0 1", "-1", true},
				{"1 1\n0", "0", true},
				{"1 1\n2", "0", false},
				{"2 2\n2 1\n1 1", "2", false},
			},
		},
		// ==================== 单词搜索 ====================
		{
			Title:      "单词搜索",
			Difficulty: "Medium",
			Tags:       "数组,回溯,矩阵",
			Description: `给定一个 m x n 二维字符网格 board 和一个字符串单词 word 。如果 word 存在于网格中，返回 true；否则返回 false。
单词必须按照字母顺序，通过相邻的单元格内的字母构成，其中"相邻"单元格是那些水平相邻或垂直相邻的单元格。同一个单元格内的字母不允许被重复使用。

示例 1：
输入：board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCCED"
输出：true

示例 2：
输入：board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "SEE"
输出：true

示例 3：
输入：board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCB"
输出：false`,
			FunctionTemplate: `func exist(board [][]byte, word string) bool {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	dims := strings.Fields(line1)
	m, _ := strconv.Atoi(dims[0])
	n, _ := strconv.Atoi(dims[1])
	board := make([][]byte, m)
	for i := 0; i < m; i++ {
		row, _ := reader.ReadString('\n')
		row = strings.TrimSpace(row)
		cols := strings.Fields(row)
		board[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			board[i][j] = cols[j][0]
		}
	}
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)
	fmt.Println(exist(board, word))
}`,
			TestCases: []TC{
				{"3 4\nA B C E\nS F C S\nA D E E\nABCCED", "true", true},
				{"3 4\nA B C E\nS F C S\nA D E E\nSEE", "true", true},
				{"3 4\nA B C E\nS F C S\nA D E E\nABCB", "false", true},
				{"1 1\nA\nA", "true", false},
				{"1 1\nA\nB", "false", false},
			},
		},
		// ==================== 搜索旋转排序数组 ====================
		{
			Title:      "搜索旋转排序数组",
			Difficulty: "Medium",
			Tags:       "数组,二分查找",
			Description: `整数数组 nums 按升序排列，数组中的值互不相同 。在传递给函数之前，nums 在预先未知的某个下标 k（0 <= k < nums.length）上进行了旋转。
给你旋转后的数组 nums 和一个整数 target ，如果 nums 中存在这个目标值，则返回它的下标，否则返回 -1 。

示例 1：
输入：nums = [4,5,6,7,0,1,2], target = 0
输出：4

示例 2：
输入：nums = [4,5,6,7,0,1,2], target = 3
输出：-1`,
			FunctionTemplate: `func search(nums []int, target int) int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	parts := strings.Fields(line1)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	target, _ := strconv.Atoi(line2)
	fmt.Println(search(nums, target))
}`,
			TestCases: []TC{
				{"4 5 6 7 0 1 2\n0", "4", true},
				{"4 5 6 7 0 1 2\n3", "-1", true},
				{"1\n0", "-1", true},
				{"1\n1", "0", false},
				{"3 1\n1", "1", false},
			},
		},
		// ==================== 在排序数组中查找元素的第一个和最后一个位置 ====================
		{
			Title:      "在排序数组中查找元素的第一个和最后一个位置",
			Difficulty: "Medium",
			Tags:       "数组,二分查找",
			Description: `给你一个按照非递减顺序排列的整数数组 nums，和一个目标值 target。请你找出给定目标值在数组中的开始位置和结束位置。
如果数组中不存在目标值 target ，返回 [-1, -1]。你必须设计并实现时间复杂度为 O(log n) 的算法。

示例 1：
输入：nums = [5,7,7,8,8,10], target = 8
输出：[3,4]

示例 2：
输入：nums = [5,7,7,8,8,10], target = 6
输出：[-1,-1]`,
			FunctionTemplate: `func searchRange(nums []int, target int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)
	parts := strings.Fields(line1)
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)
	target, _ := strconv.Atoi(line2)
	res := searchRange(nums, target)
	fmt.Printf("%d %d\n", res[0], res[1])
}`,
			TestCases: []TC{
				{"5 7 7 8 8 10\n8", "3 4", true},
				{"5 7 7 8 8 10\n6", "-1 -1", true},
				{"\n0", "-1 -1", true},
				{"1\n1", "0 0", false},
				{"2 2\n2", "0 1", false},
			},
		},
		// ==================== 每日温度 ====================
		{
			Title:      "每日温度",
			Difficulty: "Medium",
			Tags:       "栈,数组,单调栈",
			Description: `给定一个整数数组 temperatures ，表示每天的温度，返回一个数组 answer ，其中 answer[i] 是指对于第 i 天，下一个更高温度出现在几天后。如果气温在这之后都不会升高，请在该位置用 0 来代替。

示例 1：
输入：temperatures = [73,74,75,71,69,72,76,73]
输出：[1,1,4,2,1,1,0,0]

示例 2：
输入：temperatures = [30,40,50,60]
输出：[1,1,1,0]`,
			FunctionTemplate: `func dailyTemperatures(temperatures []int) []int {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	temps := make([]int, len(parts))
	for i, p := range parts {
		temps[i], _ = strconv.Atoi(p)
	}
	res := dailyTemperatures(temps)
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(strs, " "))
}`,
			TestCases: []TC{
				{"73 74 75 71 69 72 76 73", "1 1 4 2 1 1 0 0", true},
				{"30 40 50 60", "1 1 1 0", true},
				{"30 60 90", "1 1 0", true},
				{"90 80 70", "0 0 0", false},
				{"50 50 50", "0 0 0", false},
			},
		},
		// ==================== 字符串解码 ====================
		{
			Title:      "字符串解码",
			Difficulty: "Medium",
			Tags:       "栈,递归,字符串",
			Description: `给定一个经过编码的字符串，返回它解码后的字符串。
编码规则为: k[encoded_string]，表示其中方括号内部的 encoded_string 正好重复 k 次。

示例 1：
输入：s = "3[a]2[bc]"
输出："aaabcbc"

示例 2：
输入：s = "3[a2[c]]"
输出："accaccacc"

示例 3：
输入：s = "2[abc]3[cd]ef"
输出："abcabccdcdcdef"`,
			FunctionTemplate: `func decodeString(s string) string {
    // 请在此实现你的代码
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	fmt.Println(decodeString(line))
}`,
			TestCases: []TC{
				{"3[a]2[bc]", "aaabcbc", true},
				{"3[a2[c]]", "accaccacc", true},
				{"2[abc]3[cd]ef", "abcabccdcdcdef", true},
				{"abc", "abc", false},
				{"10[a]", "aaaaaaaaaa", false},
			},
		},
		// ==================== 数据流的中位数 ====================
		{
			Title:      "数据流的中位数",
			Difficulty: "Hard",
			Tags:       "设计,双指针,数据流,排序,堆",
			Description: `中位数是有序整数列表中间的数。如果列表的大小是偶数，则没有中间值，中位数是两个中间值的平均值。
实现 MedianFinder 类:
MedianFinder() 初始化 MedianFinder 对象。
void addNum(int num) 将数据流中的整数 num 添加到数据结构中。
double findMedian() 返回到目前为止所有元素的中位数。

示例：
输入：["MedianFinder","addNum","addNum","findMedian","addNum","findMedian"]
[[],[1],[2],[],[3],[]]
输出：[null,null,null,1.50,null,2.00]

输入格式：每行一个操作，addNum 后跟数字，findMedian 无参数。
输出：对每个 findMedian 输出中位数（保留两位小数），用空格分隔。`,
			FunctionTemplate: `import "container/heap"

type MedianFinder struct {
    // 请在此添加你的数据结构
}

func Constructor() MedianFinder {
    // 请在此实现你的代码
    return MedianFinder{}
}

func (this *MedianFinder) AddNum(num int) {
    // 请在此实现你的代码
}

func (this *MedianFinder) FindMedian() float64 {
    // 请在此实现你的代码
    return 0
}`,
			DriverCode: `
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	mf := Constructor()
	results := []string{}
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			if strings.HasPrefix(line, "addNum") {
				parts := strings.Fields(line)
				num, _ := strconv.Atoi(parts[1])
				mf.AddNum(num)
			} else if line == "findMedian" {
				med := mf.FindMedian()
				results = append(results, fmt.Sprintf("%.2f", med))
			}
		}
		if err != nil {
			break
		}
	}
	fmt.Println(strings.Join(results, " "))
}`,
			TestCases: []TC{
				{"addNum 1\naddNum 2\nfindMedian\naddNum 3\nfindMedian", "1.50 2.00", true},
				{"addNum 5\nfindMedian", "5.00", true},
				{"addNum 1\naddNum 1\nfindMedian", "1.00", true},
				{"addNum 3\naddNum 1\naddNum 2\nfindMedian", "2.00", false},
				{"addNum 10\naddNum 20\naddNum 30\naddNum 40\nfindMedian", "25.00", false},
			},
		},
	}

	for _, p := range problems {
		problem := models.Problem{
			Title:            p.Title,
			Description:      p.Description,
			Difficulty:       p.Difficulty,
			Tags:             p.Tags,
			TimeLimit:        5000,
			MemoryLimit:      256,
			FunctionTemplate: p.FunctionTemplate,
			DriverCode:       "package main\n" + p.DriverCode,
		}
		if err := database.DB.Create(&problem).Error; err != nil {
			log.Printf("Failed to create problem %s: %v", p.Title, err)
			continue
		}
		fmt.Printf("Created: %s (ID: %d, %d test cases)\n", p.Title, problem.ID, len(p.TestCases))
		for _, tc := range p.TestCases {
			testCase := models.TestCase{
				ProblemID: problem.ID,
				Input:     tc.Input,
				Output:    tc.Output,
				IsPublic:  tc.IsPublic,
			}
			if err := database.DB.Create(&testCase).Error; err != nil {
				log.Printf("Failed to create test case: %v", err)
			}
		}
	}
	fmt.Printf("\n导入完成，共 %d 道题目\n", len(problems))
}
