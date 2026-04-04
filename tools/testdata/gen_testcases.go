//go:build ignore
// +build ignore

// gen_testcases.go
// 为每道题自动生成50组测试用例（程序化生成，保留已有公开用例）
// 运行：go run scripts/gen_testcases.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"

	"github.com/your-org/oj-platform/internal/database"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/pkg/config"
)

const TARGET = 50 // 每题目标测试用例数

func main() {
	if err := config.Load("./config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := database.Init(&config.AppConfig.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	rng := rand.New(rand.NewSource(42))

	type generator struct {
		title string
		fn    func(rng *rand.Rand) (input, output string)
	}

	generators := []generator{
		{"两数之和", genTwoSum},
		{"爬楼梯", genClimbStairs},
		{"最大子数组和", genMaxSubArray},
		{"买卖股票的最佳时机", genMaxProfit},
		{"只出现一次的数字", genSingleNumber},
		{"多数元素", genMajorityElement},
		{"移动零", genMoveZeroes},
		{"合并两个有序数组", genMergeSortedArray},
		{"验证回文串", genIsPalindrome},
		{"找到字符串中所有字母异位词", genFindAnagrams},
		{"二叉树的最大深度", genMaxDepth},
		{"二叉树的中序遍历", genInorderTraversal},
		{"对称二叉树", genIsSymmetric},
		{"有效的括号", genIsValid},
		{"最长公共前缀", genLongestCommonPrefix},
		// 新增第一批
		{"反转链表", genReverseList},
		{"无重复字符的最长子串", genLengthOfLongestSubstring},
		{"盛最多水的容器", genMaxArea},
		{"三数之和", genThreeSum},
		{"翻转二叉树", genInvertTree},
		{"二叉树的层序遍历", genLevelOrder},
		{"最长回文子串", genLongestPalindrome},
		{"接雨水", genTrap},
		{"跳跃游戏", genCanJump},
		{"不同路径", genUniquePaths},
		{"搜索插入位置", genSearchInsert},
		{"合并区间", genMergeIntervals},
		// 补充跳过项
		{"合并两个有序链表", genMergeTwoLists},
		{"环形链表", genHasCycle},
		{"全排列", genPermute},
		{"删除链表的倒数第N个节点", genRemoveNthFromEnd},
		{"路径总和", genHasPathSum},
		{"子集", genSubsets},
		{"岛屿数量", genNumIslands},
		// 第二批新增
		{"字母异位词分组", genGroupAnagrams},
		{"最长连续序列", genLongestConsecutive},
		{"两数相加", genAddTwoNumbers},
		{"颜色分类", genSortColors},
		{"数组中的第K个最大元素", genFindKthLargest},
		{"前K个高频元素", genTopKFrequent},
		{"矩阵置零", genSetZeroes},
		{"螺旋矩阵", genSpiralOrder},
		{"旋转图像", genRotate},
		{"零钱兑换", genCoinChange},
		{"打家劫舍", genRob},
		{"完全平方数", genNumSquares},
		{"单词拆分", genWordBreak},
		{"最长递增子序列", genLengthOfLIS},
		{"乘积最大子数组", genMaxProduct},
		{"验证二叉搜索树", genIsValidBST},
		{"二叉搜索树中第K小的元素", genKthSmallest},
		{"二叉树的右视图", genRightSideView},
		{"二叉树的最近公共祖先", genLowestCommonAncestor},
		// 第三批新增
		{"二叉树的直径", genDiameterOfBinaryTree},
		{"二叉树展开为链表", genFlatten},
		{"从前序与中序遍历序列构造二叉树", genBuildTree},
		{"课程表", genCanFinish},
		{"实现Trie前缀树", genTrie},
		{"全排列II", genPermuteUnique},
		{"组合总和", genCombinationSum},
		{"电话号码的字母组合", genLetterCombinations},
		{"括号生成", genGenerateParenthesis},
		{"下一个排列", genNextPermutation},
		{"寻找重复数", genFindDuplicate},
		{"最小覆盖子串", genMinWindow},
		{"柱状图中最大的矩形", genLargestRectangle},
		{"最长有效括号", genLongestValidParentheses},
		// 第四批新增
		{"排序链表", genSortList},
		{"K个一组翻转链表", genReverseKGroup},
		{"随机链表的复制", genCopyRandomList},
		{"LRU缓存", genLRUCache},
		{"二叉树的后序遍历", genPostorderTraversal},
		{"二叉树的前序遍历", genPreorderTraversal},
		{"相交链表", genGetIntersectionNode},
		{"回文链表", genIsPalindromeList},
		{"用栈实现队列", genMyQueue},
		{"最小栈", genMinStack},
		{"二叉搜索树的插入操作", genInsertIntoBST},
		{"删除二叉搜索树中的节点", genDeleteNode},
		{"二叉树中的最大路径和", genMaxPathSum},
		{"整数转罗马数字", genIntToRoman},
		{"罗马数字转整数", genRomanToInt},
		{"编辑距离", genMinDistance},
		{"不同的二叉搜索树", genNumTrees},
		{"分隔链表", genPartition},
		{"缺失的第一个正数", genFirstMissingPositive},
		// 第五批新增（最后14道）
		{"和为K的子数组", genSubarraySum},
		{"滑动窗口最大值", genMaxSlidingWindow},
		{"轮转数组", genRotateArray},
		{"除自身以外数组的乘积", genProductExceptSelf},
		{"搜索二维矩阵II", genSearchMatrix2},
		{"环形链表II", genDetectCycle},
		{"合并K个升序链表", genMergeKLists},
		{"腐烂的橘子", genOrangesRotting},
		{"单词搜索", genWordExist},
		{"搜索旋转排序数组", genSearchRotated},
		{"在排序数组中查找元素的第一个和最后一个位置", genSearchRange},
		{"每日温度", genDailyTemperatures},
		{"字符串解码", genDecodeString},
		{"数据流的中位数", genMedianFinder},
	}

	var problems []models.Problem
	if err := database.DB.Order("id asc").Find(&problems).Error; err != nil {
		log.Fatalf("Failed to query problems: %v", err)
	}

	for _, p := range problems {
		var gen *generator
		for i := range generators {
			if generators[i].title == p.Title {
				gen = &generators[i]
				break
			}
		}
		if gen == nil {
			fmt.Printf("Skip (no generator): %s\n", p.Title)
			continue
		}

		var existing int64
		database.DB.Model(&models.TestCase{}).Where("problem_id = ?", p.ID).Count(&existing)
		need := TARGET - int(existing)
		if need <= 0 {
			fmt.Printf("Skip (already %d): %s\n", existing, p.Title)
			continue
		}

		added := 0
		attempts := 0
		for added < need && attempts < need*5 {
			attempts++
			input, output := gen.fn(rng)
			tc := models.TestCase{
				ProblemID: p.ID,
				Input:     input,
				Output:    output,
				IsPublic:  false,
			}
			if err := database.DB.Create(&tc).Error; err == nil {
				added++
			}
		}
		fmt.Printf("Added %d test cases to [%s] (total target: %d)\n", added, p.Title, TARGET)
	}
	fmt.Println("\n生成完成")
}

// ===== 两数之和 =====
func genTwoSum(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	used := map[int]bool{}
	for i := range nums {
		v := rng.Intn(201) - 100
		for used[v] {
			v = rng.Intn(201) - 100
		}
		nums[i] = v
		used[v] = true
	}
	// 随机选两个下标
	i := rng.Intn(n)
	j := (i + 1 + rng.Intn(n-1)) % n
	target := nums[i] + nums[j]
	parts := make([]string, n)
	for k, v := range nums {
		parts[k] = fmt.Sprintf("%d", v)
	}
	// 确定正确答案（可能有多对）
	a, b := -1, -1
	m := map[int]int{}
	for k, v := range nums {
		if idx, ok := m[target-v]; ok {
			a, b = idx, k
			break
		}
		m[v] = k
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", target),
		fmt.Sprintf("%d %d", a, b)
}

// ===== 爬楼梯 =====
func genClimbStairs(rng *rand.Rand) (string, string) {
	n := rng.Intn(44) + 1 // 1..44
	// fib
	a, b := 1, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return fmt.Sprintf("%d", n), fmt.Sprintf("%d", b)
}

// ===== 最大子数组和 =====
func genMaxSubArray(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(21) - 10
	}
	// Kadane
	maxSum := nums[0]
	cur := nums[0]
	for _, v := range nums[1:] {
		if cur < 0 {
			cur = v
		} else {
			cur += v
		}
		if cur > maxSum {
			maxSum = cur
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", maxSum)
}

// ===== 买卖股票 =====
func genMaxProfit(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	prices := make([]int, n)
	for i := range prices {
		prices[i] = rng.Intn(100) + 1
	}
	maxP := 0
	minP := prices[0]
	for _, p := range prices[1:] {
		if p-minP > maxP {
			maxP = p - minP
		}
		if p < minP {
			minP = p
		}
	}
	parts := make([]string, n)
	for i, v := range prices {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", maxP)
}

// ===== 只出现一次的数字 =====
func genSingleNumber(rng *rand.Rand) (string, string) {
	pairs := rng.Intn(5) + 1
	single := rng.Intn(201) - 100
	nums := []int{single}
	used := map[int]bool{single: true}
	for i := 0; i < pairs; i++ {
		v := rng.Intn(201) - 100
		for used[v] {
			v = rng.Intn(201) - 100
		}
		used[v] = true
		nums = append(nums, v, v)
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", single)
}

// ===== 多数元素 =====
func genMajorityElement(rng *rand.Rand) (string, string) {
	n := rng.Intn(9)*2 + 1 // 奇数，1~17
	majority := rng.Intn(21) - 10
	need := n/2 + 1
	nums := make([]int, 0, n)
	for i := 0; i < need; i++ {
		nums = append(nums, majority)
	}
	for len(nums) < n {
		v := rng.Intn(21) - 10
		if v != majority {
			nums = append(nums, v)
		}
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", majority)
}

// ===== 移动零 =====
func genMoveZeroes(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(5) // 0~4，增加零的比例
	}
	// in-place move zeros to end
	nonZero := []string{}
	zeroCount := 0
	for _, v := range nums {
		if v != 0 {
			nonZero = append(nonZero, fmt.Sprintf("%d", v))
		} else {
			zeroCount++
		}
	}
	result := append(nonZero, make([]string, zeroCount)...)
	for i := len(nonZero); i < len(result); i++ {
		result[i] = "0"
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), strings.Join(result, " ")
}

// ===== 合并两个有序数组 =====
func genMergeSortedArray(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 1
	n := rng.Intn(5) + 1
	a := make([]int, m)
	b := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(20)
	}
	for i := range b {
		b[i] = rng.Intn(20)
	}
	sort.Ints(a)
	sort.Ints(b)
	merged := append(a, b...)
	sort.Ints(merged)
	aParts := make([]string, m)
	for i, v := range a {
		aParts[i] = fmt.Sprintf("%d", v)
	}
	bParts := make([]string, n)
	for i, v := range b {
		bParts[i] = fmt.Sprintf("%d", v)
	}
	outParts := make([]string, m+n)
	for i, v := range merged {
		outParts[i] = fmt.Sprintf("%d", v)
	}
	input := fmt.Sprintf("%d %d\n%s\n%s", m, n, strings.Join(aParts, " "), strings.Join(bParts, " "))
	return input, strings.Join(outParts, " ")
}

// ===== 验证回文串 =====
var palindromeWords = []string{
	"A man a plan a canal Panama",
	"race a car",
	"Was it a car or a cat I saw",
	"No lemon no melon",
	"hello",
	"abba",
	"ab",
	"a",
	"",
	"0P",
}

func genIsPalindrome(rng *rand.Rand) (string, string) {
	s := palindromeWords[rng.Intn(len(palindromeWords))]
	// compute
	filtered := []byte{}
	for _, c := range []byte(strings.ToLower(s)) {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			filtered = append(filtered, c)
		}
	}
	isPalin := true
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		if filtered[i] != filtered[j] {
			isPalin = false
			break
		}
	}
	out := "false"
	if isPalin {
		out = "true"
	}
	return s, out
}

// ===== 找字母异位词 =====
func genFindAnagrams(rng *rand.Rand) (string, string) {
	pLen := rng.Intn(3) + 2   // 2~4
	sLen := rng.Intn(8) + pLen // s至少和p一样长
	letters := "abcde"
	pBytes := make([]byte, pLen)
	sBytes := make([]byte, sLen)
	for i := range pBytes {
		pBytes[i] = letters[rng.Intn(len(letters))]
	}
	for i := range sBytes {
		sBytes[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sBytes)
	p := string(pBytes)
	// find anagram positions
	freq := [26]int{}
	for _, c := range p {
		freq[c-'a']++
	}
	win := [26]int{}
	result := []string{}
	for i, c := range s {
		win[c-'a']++
		if i >= pLen {
			win[s[i-pLen]-'a']--
		}
		if i >= pLen-1 && win == freq {
			result = append(result, fmt.Sprintf("%d", i-pLen+1))
		}
	}
	out := strings.Join(result, " ")
	if out == "" {
		out = "none"
	}
	return s + "\n" + p, out
}

// ===== 二叉树最大深度 =====
func genMaxDepth(rng *rand.Rand) (string, string) {
	// 生成随机层数的完全/部分二叉树的BFS序列
	depth := rng.Intn(4) + 1 // 1~4
	maxNodes := (1 << depth) - 1
	n := rng.Intn(maxNodes) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(100)+1)
		}
	}
	// 计算深度
	h := calcDepth(vals, 0)
	return strings.Join(vals, " "), fmt.Sprintf("%d", h)
}

func calcDepth(vals []string, i int) int {
	if i >= len(vals) || vals[i] == "null" {
		return 0
	}
	l := calcDepth(vals, 2*i+1)
	r := calcDepth(vals, 2*i+2)
	if l > r {
		return l + 1
	}
	return r + 1
}

// ===== 二叉树中序遍历 =====
func genInorderTraversal(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(100)+1)
		}
	}
	result := []string{}
	inorder(vals, 0, &result)
	return strings.Join(vals, " "), strings.Join(result, " ")
}

func inorder(vals []string, i int, result *[]string) {
	if i >= len(vals) || vals[i] == "null" {
		return
	}
	inorder(vals, 2*i+1, result)
	*result = append(*result, vals[i])
	inorder(vals, 2*i+2, result)
}

// ===== 对称二叉树 =====
func genIsSymmetric(rng *rand.Rand) (string, string) {
	// 有时生成对称树，有时随机
	makeSymmetric := rng.Intn(2) == 0
	if makeSymmetric {
		depth := rng.Intn(3) + 1
		vals := buildSymmetric(depth)
		return strings.Join(vals, " "), "true"
	}
	n := rng.Intn(6) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(3) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(3)+1)
		}
	}
	sym := checkSymmetric(vals, 0, 0)
	out := "false"
	if sym {
		out = "true"
	}
	return strings.Join(vals, " "), out
}

func buildSymmetric(depth int) []string {
	if depth == 0 {
		return []string{}
	}
	size := (1 << depth) - 1
	vals := make([]string, size)
	vals[0] = "1"
	buildSym(vals, 1, 0, depth-1)
	return vals
}

func buildSym(vals []string, v, idx, rem int) {
	if rem == 0 || idx >= len(vals) {
		return
	}
	l, r := 2*idx+1, 2*idx+2
	if l < len(vals) {
		vals[l] = fmt.Sprintf("%d", v)
	}
	if r < len(vals) {
		vals[r] = fmt.Sprintf("%d", v)
	}
	buildSym(vals, v+1, l, rem-1)
	buildSym(vals, v+1, r, rem-1)
}

func checkSymmetric(vals []string, i, j int) bool {
	iNull := i >= len(vals) || vals[i] == "null"
	jNull := j >= len(vals) || vals[j] == "null"
	if iNull && jNull {
		return true
	}
	if iNull != jNull {
		return false
	}
	if vals[i] != vals[j] {
		return false
	}
	return checkSymmetric(vals, 2*i+1, 2*j+2) && checkSymmetric(vals, 2*i+2, 2*j+1)
}

// ===== 有效的括号 =====
func genIsValid(rng *rand.Rand) (string, string) {
	// 随机决定是否生成有效括号串
	valid := rng.Intn(2) == 0
	var s string
	if valid {
		s = genValidBrackets(rng)
	} else {
		s = genInvalidBrackets(rng)
	}
	// 计算
	stack := []byte{}
	ok := true
	for _, c := range []byte(s) {
		switch c {
		case '(', '[', '{':
			stack = append(stack, c)
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				ok = false
			} else {
				stack = stack[:len(stack)-1]
			}
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				ok = false
			} else {
				stack = stack[:len(stack)-1]
			}
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				ok = false
			} else {
				stack = stack[:len(stack)-1]
			}
		}
	}
	if len(stack) > 0 {
		ok = false
	}
	out := "false"
	if ok {
		out = "true"
	}
	return s, out
}

func genValidBrackets(rng *rand.Rand) string {
	pairs := rng.Intn(4) + 1
	types := []string{"()", "[]", "{}"}
	result := ""
	for i := 0; i < pairs; i++ {
		t := types[rng.Intn(3)]
		result += string(t[0])
		if rng.Intn(2) == 0 && i < pairs-1 {
			inner := types[rng.Intn(3)]
			result += string(inner[0]) + string(inner[1])
		}
		result += string(t[1])
	}
	return result
}

func genInvalidBrackets(rng *rand.Rand) string {
	samples := []string{"(]", "([)]", "{[]", "}", "((", ")(", "]"}
	return samples[rng.Intn(len(samples))]
}

// ===== 最长公共前缀 =====
func genLongestCommonPrefix(rng *rand.Rand) (string, string) {
	words := []string{"flower", "flow", "flight", "flock", "float",
		"dog", "racecar", "car", "abc", "ab", "a", "ab", "apple", "app", "application"}
	n := rng.Intn(4) + 2
	chosen := make([]string, n)
	for i := range chosen {
		chosen[i] = words[rng.Intn(len(words))]
	}
	// lcp
	prefix := chosen[0]
	for _, w := range chosen[1:] {
		for !strings.HasPrefix(w, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				break
			}
		}
	}
	return strings.Join(chosen, "\n"), prefix
}

// ===== 新增题目的生成器 =====

// 反转链表
func genReverseList(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	vals := make([]string, n)
	for i := range vals {
		vals[i] = fmt.Sprintf("%d", rng.Intn(100)+1)
	}
	rev := make([]string, n)
	for i := range rev {
		rev[i] = vals[n-1-i]
	}
	return strings.Join(vals, " "), strings.Join(rev, " ")
}

// 无重复字符的最长子串
func genLengthOfLongestSubstring(rng *rand.Rand) (string, string) {
	letters := "abcdefghij"
	n := rng.Intn(10) + 1
	bs := make([]byte, n)
	for i := range bs {
		bs[i] = letters[rng.Intn(len(letters))]
	}
	s := string(bs)
	// 计算
	maxLen, start := 0, 0
	pos := map[byte]int{}
	for i := 0; i < len(s); i++ {
		if p, ok := pos[s[i]]; ok && p >= start {
			start = p + 1
		}
		pos[s[i]] = i
		if i-start+1 > maxLen {
			maxLen = i - start + 1
		}
	}
	return s, fmt.Sprintf("%d", maxLen)
}

// 盛最多水的容器
func genMaxArea(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	height := make([]int, n)
	for i := range height {
		height[i] = rng.Intn(20) + 1
	}
	maxA := 0
	l, r := 0, n-1
	for l < r {
		h := height[l]
		if height[r] < h {
			h = height[r]
		}
		a := h * (r - l)
		if a > maxA {
			maxA = a
		}
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
	}
	parts := make([]string, n)
	for i, v := range height {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", maxA)
}

// 三数之和
func genThreeSum(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 3
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(11) - 5
	}
	// 计算结果
	sort.Ints(nums)
	seen := map[[3]int]bool{}
	var results [][3]int
	for i := 0; i < n-2; i++ {
		l, r := i+1, n-1
		for l < r {
			s := nums[i] + nums[l] + nums[r]
			if s == 0 {
				t := [3]int{nums[i], nums[l], nums[r]}
				if !seen[t] {
					seen[t] = true
					results = append(results, t)
				}
				l++
				r--
			} else if s < 0 {
				l++
			} else {
				r--
			}
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	if len(results) == 0 {
		return strings.Join(parts, " "), ""
	}
	outLines := make([]string, len(results))
	for i, t := range results {
		outLines[i] = fmt.Sprintf("%d %d %d", t[0], t[1], t[2])
	}
	return strings.Join(parts, " "), strings.Join(outLines, "\n")
}

// 翻转二叉树
func genInvertTree(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(10)+1)
		}
	}
	// 计算翻转后的BFS序列
	invertBFS(vals, 0)
	// 去掉末尾null
	end := len(vals)
	for end > 0 && vals[end-1] == "null" {
		end--
	}
	return strings.Join(vals, " "), strings.Join(vals[:end], " ")
}

func invertBFS(vals []string, i int) {
	if i >= len(vals) || vals[i] == "null" {
		return
	}
	l, r := 2*i+1, 2*i+2
	if l < len(vals) && r < len(vals) {
		vals[l], vals[r] = vals[r], vals[l]
	}
	invertBFS(vals, l)
	invertBFS(vals, r)
}

// 二叉树的层序遍历
func genLevelOrder(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(3) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(20)+1)
		}
	}
	// 计算层序输出
	levels := [][]string{}
	queue := []int{0}
	for len(queue) > 0 {
		nextQ := []int{}
		level := []string{}
		for _, idx := range queue {
			if idx >= len(vals) || vals[idx] == "null" {
				continue
			}
			level = append(level, vals[idx])
			nextQ = append(nextQ, 2*idx+1, 2*idx+2)
		}
		if len(level) > 0 {
			levels = append(levels, level)
		}
		queue = nextQ
	}
	out := []string{}
	for _, l := range levels {
		out = append(out, strings.Join(l, " "))
	}
	return strings.Join(vals, " "), strings.Join(out, "\n")
}

// 最长回文子串
func genLongestPalindrome(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	bs := make([]byte, n)
	letters := "abcde"
	for i := range bs {
		bs[i] = letters[rng.Intn(len(letters))]
	}
	s := string(bs)
	// 中心扩展法
	best := s[0:1]
	expand := func(l, r int) {
		for l >= 0 && r < len(s) && s[l] == s[r] {
			if r-l+1 > len(best) {
				best = s[l : r+1]
			}
			l--
			r++
		}
	}
	for i := range s {
		expand(i, i)
		expand(i, i+1)
	}
	return s, best
}

// 接雨水
func genTrap(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	height := make([]int, n)
	for i := range height {
		height[i] = rng.Intn(7)
	}
	// 双指针计算
	water := 0
	l, r := 0, n-1
	lMax, rMax := 0, 0
	for l < r {
		if height[l] < height[r] {
			if height[l] >= lMax {
				lMax = height[l]
			} else {
				water += lMax - height[l]
			}
			l++
		} else {
			if height[r] >= rMax {
				rMax = height[r]
			} else {
				water += rMax - height[r]
			}
			r--
		}
	}
	parts := make([]string, n)
	for i, v := range height {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", water)
}

// 跳跃游戏
func genCanJump(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(5)
	}
	// 贪心计算
	maxReach := 0
	can := true
	for i, v := range nums {
		if i > maxReach {
			can = false
			break
		}
		if i+v > maxReach {
			maxReach = i + v
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	out := "false"
	if can {
		out = "true"
	}
	return strings.Join(parts, " "), out
}

// 不同路径
func genUniquePaths(rng *rand.Rand) (string, string) {
	m := rng.Intn(8) + 1
	n := rng.Intn(8) + 1
	// DP 计算
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][0] = 1
	}
	for j := 0; j < n; j++ {
		dp[0][j] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = dp[i-1][j] + dp[i][j-1]
		}
	}
	return fmt.Sprintf("%d %d", m, n), fmt.Sprintf("%d", dp[m-1][n-1])
}

// 搜索插入位置
func genSearchInsert(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	// 生成不重复有序数组
	seen := map[int]bool{}
	nums := []int{}
	for len(nums) < n {
		v := rng.Intn(30)
		if !seen[v] {
			seen[v] = true
			nums = append(nums, v)
		}
	}
	sort.Ints(nums)
	target := rng.Intn(35) // 可能在范围外
	// 二分
	lo, hi := 0, len(nums)
	for lo < hi {
		mid := (lo + hi) / 2
		if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", target), fmt.Sprintf("%d", lo)
}

// 合并两个有序链表
func genMergeTwoLists(rng *rand.Rand) (string, string) {
	n1 := rng.Intn(5) + 1
	n2 := rng.Intn(5) + 1
	a := make([]int, n1)
	b := make([]int, n2)
	for i := range a {
		a[i] = rng.Intn(20)
	}
	for i := range b {
		b[i] = rng.Intn(20)
	}
	sort.Ints(a)
	sort.Ints(b)
	// 合并
	merged := append(append([]int{}, a...), b...)
	sort.Ints(merged)
	aParts := make([]string, n1)
	for i, v := range a {
		aParts[i] = fmt.Sprintf("%d", v)
	}
	bParts := make([]string, n2)
	for i, v := range b {
		bParts[i] = fmt.Sprintf("%d", v)
	}
	mParts := make([]string, len(merged))
	for i, v := range merged {
		mParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(aParts, " ") + "\n" + strings.Join(bParts, " "), strings.Join(mParts, " ")
}

// 环形链表（判断是否有环，使用特殊输入格式）
func genHasCycle(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	vals := make([]string, n)
	for i := range vals {
		vals[i] = fmt.Sprintf("%d", rng.Intn(100))
	}
	hasCycle := rng.Intn(2) == 0
	pos := -1
	if hasCycle {
		pos = rng.Intn(n)
	}
	out := "false"
	if hasCycle {
		out = "true"
	}
	return strings.Join(vals, " ") + "\n" + fmt.Sprintf("%d", pos), out
}

// 全排列
func genPermute(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2 // 2~5
	// 生成不重复数字
	seen := map[int]bool{}
	nums := []int{}
	for len(nums) < n {
		v := rng.Intn(11) - 5
		if !seen[v] {
			seen[v] = true
			nums = append(nums, v)
		}
	}
	// 生成所有全排列
	var perms [][]int
	var permHelper func([]int, []int)
	permHelper = func(cur, rem []int) {
		if len(rem) == 0 {
			tmp := make([]int, len(cur))
			copy(tmp, cur)
			perms = append(perms, tmp)
			return
		}
		for i, v := range rem {
			next := append([]int{}, rem[:i]...)
			next = append(next, rem[i+1:]...)
			permHelper(append(cur, v), next)
		}
	}
	permHelper([]int{}, nums)
	// 对排列排序（字典序）
	sort.Slice(perms, func(i, j int) bool {
		for k := range perms[i] {
			if perms[i][k] != perms[j][k] {
				return perms[i][k] < perms[j][k]
			}
		}
		return false
	})
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	outLines := make([]string, len(perms))
	for i, p := range perms {
		pp := make([]string, len(p))
		for j, v := range p {
			pp[j] = fmt.Sprintf("%d", v)
		}
		outLines[i] = strings.Join(pp, " ")
	}
	return strings.Join(parts, " "), strings.Join(outLines, "\n")
}

// 删除链表的倒数第N个节点
func genRemoveNthFromEnd(rng *rand.Rand) (string, string) {
	size := rng.Intn(7) + 2
	vals := make([]int, size)
	for i := range vals {
		vals[i] = rng.Intn(100) + 1
	}
	n := rng.Intn(size) + 1 // 1..size
	// 删除倒数第n个
	removeIdx := size - n
	result := append(append([]int{}, vals[:removeIdx]...), vals[removeIdx+1:]...)
	parts := make([]string, size)
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	resParts := make([]string, len(result))
	for i, v := range result {
		resParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", n), strings.Join(resParts, " ")
}

// 路径总和
func genHasPathSum(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	vals := make([]string, n)
	intVals := make([]int, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
			intVals[i] = 0
		} else {
			v := rng.Intn(21) - 10
			vals[i] = fmt.Sprintf("%d", v)
			intVals[i] = v
		}
	}
	// 计算所有根到叶子路径的和
	var pathSums []int
	var dfs func(idx, curSum int)
	dfs = func(idx, curSum int) {
		if idx >= len(vals) || vals[idx] == "null" {
			return
		}
		curSum += intVals[idx]
		l, r := 2*idx+1, 2*idx+2
		lNull := l >= len(vals) || vals[l] == "null"
		rNull := r >= len(vals) || vals[r] == "null"
		if lNull && rNull {
			pathSums = append(pathSums, curSum)
			return
		}
		dfs(l, curSum)
		dfs(r, curSum)
	}
	dfs(0, 0)
	// 随机选一个目标（有时选存在的，有时选不存在的）
	var target int
	hasPath := false
	if len(pathSums) > 0 && rng.Intn(2) == 0 {
		target = pathSums[rng.Intn(len(pathSums))]
		hasPath = true
	} else {
		target = rng.Intn(41) - 20
		for _, s := range pathSums {
			if s == target {
				hasPath = true
				break
			}
		}
	}
	out := "false"
	if hasPath {
		out = "true"
	}
	return strings.Join(vals, " ") + "\n" + fmt.Sprintf("%d", target), out
}

// 子集
func genSubsets(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1 // 1~4
	seen := map[int]bool{}
	nums := []int{}
	for len(nums) < n {
		v := rng.Intn(11) - 5
		if !seen[v] {
			seen[v] = true
			nums = append(nums, v)
		}
	}
	sort.Ints(nums)
	// 生成所有子集
	total := 1 << n
	subsets := make([][]int, 0, total)
	for mask := 0; mask < total; mask++ {
		sub := []int{}
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sub = append(sub, nums[i])
			}
		}
		subsets = append(subsets, sub)
	}
	// 按长度排序，长度相同按字典序
	sort.Slice(subsets, func(i, j int) bool {
		if len(subsets[i]) != len(subsets[j]) {
			return len(subsets[i]) < len(subsets[j])
		}
		for k := range subsets[i] {
			if subsets[i][k] != subsets[j][k] {
				return subsets[i][k] < subsets[j][k]
			}
		}
		return false
	})
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	outLines := make([]string, len(subsets))
	for i, sub := range subsets {
		if len(sub) == 0 {
			outLines[i] = ""
		} else {
			pp := make([]string, len(sub))
			for j, v := range sub {
				pp[j] = fmt.Sprintf("%d", v)
			}
			outLines[i] = strings.Join(pp, " ")
		}
	}
	return strings.Join(parts, " "), strings.Join(outLines, "\n")
}

// 岛屿数量
func genNumIslands(rng *rand.Rand) (string, string) {
	rows := rng.Intn(4) + 2
	cols := rng.Intn(4) + 2
	grid := make([][]byte, rows)
	for i := range grid {
		grid[i] = make([]byte, cols)
		for j := range grid[i] {
			if rng.Intn(3) == 0 {
				grid[i][j] = '0'
			} else {
				grid[i][j] = '1'
			}
		}
	}
	// BFS计算岛屿数量
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}
	count := 0
	var bfs func(r, c int)
	bfs = func(r, c int) {
		if r < 0 || r >= rows || c < 0 || c >= cols || visited[r][c] || grid[r][c] == '0' {
			return
		}
		visited[r][c] = true
		bfs(r+1, c)
		bfs(r-1, c)
		bfs(r, c+1)
		bfs(r, c-1)
	}
	for i := range grid {
		for j := range grid[i] {
			if !visited[i][j] && grid[i][j] == '1' {
				bfs(i, j)
				count++
			}
		}
	}
	lines := make([]string, rows)
	for i, row := range grid {
		lines[i] = string(row)
	}
	return fmt.Sprintf("%d %d\n%s", rows, cols, strings.Join(lines, "\n")), fmt.Sprintf("%d", count)
}

// 合并区间
func genMergeIntervals(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	type interval struct{ s, e int }
	intervals := make([]interval, n)
	for i := range intervals {
		s := rng.Intn(20)
		e := s + rng.Intn(8)
		intervals[i] = interval{s, e}
	}
	sort.Slice(intervals, func(i, j int) bool { return intervals[i].s < intervals[j].s })
	// 合并
	merged := []interval{intervals[0]}
	for _, iv := range intervals[1:] {
		last := &merged[len(merged)-1]
		if iv.s <= last.e {
			if iv.e > last.e {
				last.e = iv.e
			}
		} else {
			merged = append(merged, iv)
		}
	}
	inLines := make([]string, n)
	for i, iv := range intervals {
		inLines[i] = fmt.Sprintf("%d %d", iv.s, iv.e)
	}
	outLines := make([]string, len(merged))
	for i, iv := range merged {
		outLines[i] = fmt.Sprintf("%d %d", iv.s, iv.e)
	}
	return strings.Join(inLines, "\n"), strings.Join(outLines, "\n")
}

// ===== 第二批生成器 =====

// 字母异位词分组
func genGroupAnagrams(rng *rand.Rand) (string, string) {
	allWords := [][]string{
		{"eat", "tea", "ate"},
		{"tan", "nat"},
		{"bat"},
		{"abc", "bca", "cab"},
		{"xyz"},
		{"dog", "god"},
	}
	groupCount := rng.Intn(3) + 2
	chosen := []string{}
	expectedGroups := map[string][]string{}
	for i := 0; i < groupCount && i < len(allWords); i++ {
		group := allWords[i]
		pick := rng.Intn(len(group)) + 1
		if pick > len(group) {
			pick = len(group)
		}
		for j := 0; j < pick; j++ {
			w := group[j]
			chosen = append(chosen, w)
			key := sortStr(w)
			expectedGroups[key] = append(expectedGroups[key], w)
		}
	}
	type groupT struct {
		key  string
		vals []string
	}
	var groups []groupT
	for k, v := range expectedGroups {
		sort.Strings(v)
		groups = append(groups, groupT{k, v})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].vals[0] < groups[j].vals[0]
	})
	outLines := make([]string, len(groups))
	for i, g := range groups {
		outLines[i] = strings.Join(g.vals, " ")
	}
	return strings.Join(chosen, " "), strings.Join(outLines, "\n")
}

func sortStr(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}

// 最长连续序列
func genLongestConsecutive(rng *rand.Rand) (string, string) {
	nums := []int{}
	seenSet := map[int]bool{}
	start := rng.Intn(50)
	seqLen := rng.Intn(6) + 2
	for i := 0; i < seqLen; i++ {
		v := start + i
		if !seenSet[v] {
			seenSet[v] = true
			nums = append(nums, v)
		}
	}
	extra := rng.Intn(5) + 2
	for k := 0; k < extra; k++ {
		v := rng.Intn(100)
		if !seenSet[v] {
			seenSet[v] = true
			nums = append(nums, v)
		}
	}
	numSet := map[int]bool{}
	for _, v := range nums {
		numSet[v] = true
	}
	best := 0
	for v := range numSet {
		if !numSet[v-1] {
			cur := v
			length := 1
			for numSet[cur+1] {
				cur++
				length++
			}
			if length > best {
				best = length
			}
		}
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", best)
}

// 两数相加（链表加法）
func genAddTwoNumbers(rng *rand.Rand) (string, string) {
	a := rng.Intn(999) + 1
	b := rng.Intn(999) + 1
	toDigits := func(n int) []string {
		if n == 0 {
			return []string{"0"}
		}
		ds := []string{}
		for n > 0 {
			ds = append(ds, fmt.Sprintf("%d", n%10))
			n /= 10
		}
		return ds
	}
	return strings.Join(toDigits(a), " ") + "\n" + strings.Join(toDigits(b), " "),
		strings.Join(toDigits(a+b), " ")
}

// 颜色分类（荷兰国旗）
func genSortColors(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(3)
	}
	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Ints(sorted)
	parts := make([]string, n)
	sParts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range sorted {
		sParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), strings.Join(sParts, " ")
}

// 数组中第K大元素
func genFindKthLargest(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(20) + 1
	}
	k := rng.Intn(n) + 1
	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", k), fmt.Sprintf("%d", sorted[k-1])
}

// 前K个高频元素
func genTopKFrequent(rng *rand.Rand) (string, string) {
	k := rng.Intn(3) + 1
	freq := map[int]int{}
	nums := []int{}
	for i := 0; i < k+2; i++ {
		v := (i + 1) * 10
		cnt := rng.Intn(4) + 2
		for j := 0; j < cnt; j++ {
			nums = append(nums, v)
		}
		freq[v] = cnt
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	type kv struct{ val, cnt int }
	var kvs []kv
	for v, c := range freq {
		kvs = append(kvs, kv{v, c})
	}
	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].cnt != kvs[j].cnt {
			return kvs[i].cnt > kvs[j].cnt
		}
		return kvs[i].val < kvs[j].val
	})
	result := []int{}
	for i := 0; i < k && i < len(kvs); i++ {
		result = append(result, kvs[i].val)
	}
	sort.Ints(result)
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	resParts := make([]string, len(result))
	for i, v := range result {
		resParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", k), strings.Join(resParts, " ")
}

// 矩阵置零
func genSetZeroes(rng *rand.Rand) (string, string) {
	rows := rng.Intn(3) + 2
	cols := rng.Intn(3) + 2
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			if rng.Intn(5) == 0 {
				matrix[i][j] = 0
			} else {
				matrix[i][j] = rng.Intn(9) + 1
			}
		}
	}
	zeroRows := map[int]bool{}
	zeroCols := map[int]bool{}
	for i, row := range matrix {
		for j, v := range row {
			if v == 0 {
				zeroRows[i] = true
				zeroCols[j] = true
			}
		}
	}
	result := make([][]int, rows)
	for i := range result {
		result[i] = make([]int, cols)
		copy(result[i], matrix[i])
		for j := range result[i] {
			if zeroRows[i] || zeroCols[j] {
				result[i][j] = 0
			}
		}
	}
	inLines := make([]string, rows)
	outLines := make([]string, rows)
	for i := range matrix {
		rp := make([]string, cols)
		op := make([]string, cols)
		for j := range matrix[i] {
			rp[j] = fmt.Sprintf("%d", matrix[i][j])
			op[j] = fmt.Sprintf("%d", result[i][j])
		}
		inLines[i] = strings.Join(rp, " ")
		outLines[i] = strings.Join(op, " ")
	}
	return strings.Join(inLines, "\n"), strings.Join(outLines, "\n")
}

// 螺旋矩阵
func genSpiralOrder(rng *rand.Rand) (string, string) {
	m := rng.Intn(3) + 2
	n := rng.Intn(3) + 2
	matrix := make([][]int, m)
	val := 1
	for i := range matrix {
		matrix[i] = make([]int, n)
		for j := range matrix[i] {
			matrix[i][j] = val
			val++
		}
	}
	result := []string{}
	top, bottom, left, right := 0, m-1, 0, n-1
	for top <= bottom && left <= right {
		for i := left; i <= right; i++ {
			result = append(result, fmt.Sprintf("%d", matrix[top][i]))
		}
		top++
		for i := top; i <= bottom; i++ {
			result = append(result, fmt.Sprintf("%d", matrix[i][right]))
		}
		right--
		if top <= bottom {
			for i := right; i >= left; i-- {
				result = append(result, fmt.Sprintf("%d", matrix[bottom][i]))
			}
			bottom--
		}
		if left <= right {
			for i := bottom; i >= top; i-- {
				result = append(result, fmt.Sprintf("%d", matrix[i][left]))
			}
			left++
		}
	}
	inLines := make([]string, m)
	for i, row := range matrix {
		pp := make([]string, n)
		for j, v := range row {
			pp[j] = fmt.Sprintf("%d", v)
		}
		inLines[i] = strings.Join(pp, " ")
	}
	return strings.Join(inLines, "\n"), strings.Join(result, " ")
}

// 旋转图像
func genRotate(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 2
	matrix := make([][]int, n)
	val := 1
	for i := range matrix {
		matrix[i] = make([]int, n)
		for j := range matrix[i] {
			matrix[i][j] = val
			val++
		}
	}
	rotated := make([][]int, n)
	for i := range rotated {
		rotated[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			rotated[j][n-1-i] = matrix[i][j]
		}
	}
	inLines := make([]string, n)
	outLines := make([]string, n)
	for i := range matrix {
		ip := make([]string, n)
		op := make([]string, n)
		for j := range matrix[i] {
			ip[j] = fmt.Sprintf("%d", matrix[i][j])
			op[j] = fmt.Sprintf("%d", rotated[i][j])
		}
		inLines[i] = strings.Join(ip, " ")
		outLines[i] = strings.Join(op, " ")
	}
	return strings.Join(inLines, "\n"), strings.Join(outLines, "\n")
}

// 零钱兑换
func genCoinChange(rng *rand.Rand) (string, string) {
	coinSets := [][]int{
		{1, 2, 5},
		{1, 5, 10, 25},
		{2, 5, 10},
		{1, 3, 4},
		{2},
		{1},
	}
	coins := coinSets[rng.Intn(len(coinSets))]
	amount := rng.Intn(15) + 1
	dp := make([]int, amount+1)
	const bigN = 1<<30 - 1
	for i := 1; i <= amount; i++ {
		dp[i] = bigN
	}
	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c] != bigN && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	ans := dp[amount]
	if ans == bigN {
		ans = -1
	}
	parts := make([]string, len(coins))
	for i, c := range coins {
		parts[i] = fmt.Sprintf("%d", c)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", amount), fmt.Sprintf("%d", ans)
}

// 打家劫舍
func genRob(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(100) + 1
	}
	if n == 1 {
		return fmt.Sprintf("%d", nums[0]), fmt.Sprintf("%d", nums[0])
	}
	a, b := nums[0], nums[1]
	if nums[0] > b {
		b = nums[0]
	}
	for i := 2; i < n; i++ {
		nb := b
		if a+nums[i] > nb {
			nb = a + nums[i]
		}
		a, b = b, nb
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", b)
}

// 完全平方数
func genNumSquares(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 1
	dp := make([]int, n+1)
	const bigN = 1<<30 - 1
	for i := 1; i <= n; i++ {
		dp[i] = bigN
	}
	for i := 1; i <= n; i++ {
		for j := 1; j*j <= i; j++ {
			if dp[i-j*j]+1 < dp[i] {
				dp[i] = dp[i-j*j] + 1
			}
		}
	}
	return fmt.Sprintf("%d", n), fmt.Sprintf("%d", dp[n])
}

// 单词拆分
func genWordBreak(rng *rand.Rand) (string, string) {
	type pair struct {
		s    string
		dict []string
		ans  string
	}
	cases := []pair{
		{"leetcode", []string{"leet", "code"}, "true"},
		{"applepenapple", []string{"apple", "pen"}, "true"},
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}, "false"},
		{"aaaa", []string{"aa", "aaa"}, "true"},
		{"abcd", []string{"ab", "cd"}, "true"},
		{"ab", []string{"a", "b"}, "true"},
		{"hello", []string{"he", "world"}, "false"},
		{"dogs", []string{"dog", "s"}, "true"},
	}
	c := cases[rng.Intn(len(cases))]
	return c.s + "\n" + strings.Join(c.dict, " "), c.ans
}

// 最长递增子序列
func genLengthOfLIS(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(20) + 1
	}
	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1
	}
	best := 1
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", best)
}

// 乘积最大子数组
func genMaxProduct(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	nums := make([]int, n)
	for i := range nums {
		v := rng.Intn(7) - 3
		if v == 0 {
			v = 1
		}
		nums[i] = v
	}
	maxP, minP, best := nums[0], nums[0], nums[0]
	for i := 1; i < n; i++ {
		cands := [3]int{nums[i], maxP * nums[i], minP * nums[i]}
		sort.Ints(cands[:])
		minP, maxP = cands[0], cands[2]
		if maxP > best {
			best = maxP
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", best)
}

// 验证二叉搜索树
func genIsValidBST(rng *rand.Rand) (string, string) {
	// 生成有效BST
	if rng.Intn(2) == 0 {
		n := rng.Intn(5) + 2
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i*5 + rng.Intn(4) + 1
		}
		treeVals := buildBSTLO(vals)
		strs := make([]string, len(treeVals))
		for i, v := range treeVals {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, " "), "true"
	}
	// 随机树，计算是否为BST
	n := rng.Intn(4) + 3
	vals := make([]string, n)
	for i := range vals {
		vals[i] = fmt.Sprintf("%d", rng.Intn(10)+1)
	}
	ok := checkBSTArr(vals, 0, -(1 << 30), 1<<30)
	out := "false"
	if ok {
		out = "true"
	}
	return strings.Join(vals, " "), out
}

func buildBSTLO(sorted []int) []int {
	result := make([]int, len(sorted))
	var build func(lo, hi, idx int)
	build = func(lo, hi, idx int) {
		if lo > hi || idx >= len(sorted) {
			return
		}
		mid := (lo + hi) / 2
		result[idx] = sorted[mid]
		build(lo, mid-1, 2*idx+1)
		build(mid+1, hi, 2*idx+2)
	}
	build(0, len(sorted)-1, 0)
	return result
}

func checkBSTArr(vals []string, i, lo, hi int) bool {
	if i >= len(vals) || vals[i] == "null" {
		return true
	}
	var cur int
	fmt.Sscanf(vals[i], "%d", &cur)
	if cur <= lo || cur >= hi {
		return false
	}
	return checkBSTArr(vals, 2*i+1, lo, cur) && checkBSTArr(vals, 2*i+2, cur, hi)
}

// 二叉搜索树中第K小的元素
func genKthSmallest(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i*3 + rng.Intn(3) + 1
	}
	sort.Ints(vals)
	k := rng.Intn(n) + 1
	kth := vals[k-1]
	treeVals := buildBSTLO(vals)
	strs := make([]string, len(treeVals))
	for i, v := range treeVals {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, " ") + "\n" + fmt.Sprintf("%d", k), fmt.Sprintf("%d", kth)
}

// 二叉树的右视图
func genRightSideView(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(20)+1)
		}
	}
	rightView := []string{}
	queue := []int{0}
	for len(queue) > 0 {
		nextQ := []int{}
		lastVal := ""
		for _, idx := range queue {
			if idx >= len(vals) || vals[idx] == "null" {
				continue
			}
			lastVal = vals[idx]
			nextQ = append(nextQ, 2*idx+1, 2*idx+2)
		}
		if lastVal != "" {
			rightView = append(rightView, lastVal)
		}
		queue = nextQ
	}
	return strings.Join(vals, " "), strings.Join(rightView, " ")
}

// 二叉树的最近公共祖先
func genLowestCommonAncestor(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 3
	vals := make([]string, n)
	intVals := make([]int, n)
	for i := range vals {
		v := (i + 1) * 2
		vals[i] = fmt.Sprintf("%d", v)
		intVals[i] = v
	}
	nonNull := []int{}
	for i := range vals {
		nonNull = append(nonNull, i)
	}
	rng.Shuffle(len(nonNull), func(i, j int) { nonNull[i], nonNull[j] = nonNull[j], nonNull[i] })
	pIdx, qIdx := nonNull[0], nonNull[1]
	// find LCA index
	lcaIdx := lcaIndex(pIdx, qIdx)
	return strings.Join(vals, " ") + "\n" + fmt.Sprintf("%d %d", intVals[pIdx], intVals[qIdx]),
		fmt.Sprintf("%d", intVals[lcaIdx])
}

func lcaIndex(p, q int) int {
	for p != q {
		if p > q {
			p = (p - 1) / 2
		} else {
			q = (q - 1) / 2
		}
	}
	return p
}
// ===== 第三批生成器 =====

// 二叉树的直径
func genDiameterOfBinaryTree(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(100)+1)
		}
	}
	// 计算直径 = 左深度+右深度的最大值
	best := 0
	var dfs func(i int) int
	dfs = func(i int) int {
		if i >= len(vals) || vals[i] == "null" {
			return 0
		}
		l := dfs(2*i + 1)
		r := dfs(2*i + 2)
		if l+r > best {
			best = l + r
		}
		if l > r {
			return l + 1
		}
		return r + 1
	}
	dfs(0)
	return strings.Join(vals, " "), fmt.Sprintf("%d", best)
}

// 二叉树展开为链表（前序遍历顺序）
func genFlatten(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	vals := make([]string, n)
	intVals := []int{}
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			v := rng.Intn(20) + 1
			vals[i] = fmt.Sprintf("%d", v)
		}
	}
	// 前序遍历
	var preorder func(i int)
	preorder = func(i int) {
		if i >= len(vals) || vals[i] == "null" {
			return
		}
		v, _ := fmt.Sscanf(vals[i], "%d", new(int))
		_ = v
		var cur int
		fmt.Sscanf(vals[i], "%d", &cur)
		intVals = append(intVals, cur)
		preorder(2*i + 1)
		preorder(2*i + 2)
	}
	preorder(0)
	parts := make([]string, len(intVals))
	for i, v := range intVals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(vals, " "), strings.Join(parts, " ")
}

// 从前序与中序构造二叉树
func genBuildTree(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i*3 + rng.Intn(3) + 1
	}
	// 构建一棵随机树，然后计算其前序和中序
	// 为简单起见，构建满足BST的层序树，然后计算遍历
	treeArr := buildBSTLO(vals)
	var preArr, inArr []int
	var dfsTree func(i int)
	dfsTree = func(i int) {
		if i >= len(treeArr) || treeArr[i] == 0 {
			return
		}
		preArr = append(preArr, treeArr[i])
		dfsTree(2*i + 1)
		inArr = append(inArr, treeArr[i])
		dfsTree(2*i + 2)
	}
	dfsTree(0)
	// 输出层序（重建验证用原始vals顺序）
	expected := make([]string, len(treeArr))
	for i, v := range treeArr {
		expected[i] = fmt.Sprintf("%d", v)
	}
	preParts := make([]string, len(preArr))
	inParts := make([]string, len(inArr))
	for i, v := range preArr {
		preParts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range inArr {
		inParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(preParts, " ") + "\n" + strings.Join(inParts, " "),
		strings.Join(expected, " ")
}

// 课程表（拓扑排序）
func genCanFinish(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	// 随机决定有无环
	hasCycle := rng.Intn(3) == 0
	edges := []string{}
	if !hasCycle {
		// 生成DAG：只允许i->j where j<i
		for i := 1; i < n; i++ {
			if rng.Intn(2) == 0 {
				j := rng.Intn(i)
				edges = append(edges, fmt.Sprintf("%d %d", i, j))
			}
		}
		return fmt.Sprintf("%d\n%s", n, strings.Join(edges, "\n")), "true"
	}
	// 生成有环：建一个环
	cycleLen := rng.Intn(n-1) + 2
	for i := 0; i < cycleLen; i++ {
		edges = append(edges, fmt.Sprintf("%d %d", i, (i+1)%cycleLen))
	}
	return fmt.Sprintf("%d\n%s", n, strings.Join(edges, "\n")), "false"
}

// 实现Trie前缀树
func genTrie(rng *rand.Rand) (string, string) {
	words := []string{"apple", "app", "application", "bat", "ball", "band", "cat", "car", "card"}
	// 随机选择几个词插入，然后查询
	n := rng.Intn(4) + 2
	inserted := []string{}
	for i := 0; i < n; i++ {
		inserted = append(inserted, words[rng.Intn(len(words))])
	}
	ops := []string{}
	results := []string{}
	insertedSet := map[string]bool{}
	for _, w := range inserted {
		ops = append(ops, "insert "+w)
		insertedSet[w] = true
	}
	// 查询一些词
	queries := append([]string{}, inserted...)
	queries = append(queries, words[rng.Intn(len(words))])
	for _, q := range queries {
		ops = append(ops, "search "+q)
		if insertedSet[q] {
			results = append(results, "true")
		} else {
			results = append(results, "false")
		}
	}
	// 前缀查询
	if len(inserted) > 0 {
		prefix := inserted[0]
		if len(prefix) > 2 {
			prefix = prefix[:2]
		}
		ops = append(ops, "startsWith "+prefix)
		hasPrefix := false
		for w := range insertedSet {
			if strings.HasPrefix(w, prefix) {
				hasPrefix = true
				break
			}
		}
		if hasPrefix {
			results = append(results, "true")
		} else {
			results = append(results, "false")
		}
	}
	return strings.Join(ops, "\n"), strings.Join(results, " ")
}

// 全排列II（有重复元素）
func genPermuteUnique(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(3) + 1 // 1~3，容易产生重复
	}
	sort.Ints(nums)
	// 生成所有不重复全排列
	var perms [][]int
	used := make([]bool, n)
	var bt func(cur []int)
	bt = func(cur []int) {
		if len(cur) == n {
			tmp := make([]int, n)
			copy(tmp, cur)
			perms = append(perms, tmp)
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			bt(append(cur, nums[i]))
			used[i] = false
		}
	}
	bt([]int{})
	sort.Slice(perms, func(i, j int) bool {
		for k := range perms[i] {
			if perms[i][k] != perms[j][k] {
				return perms[i][k] < perms[j][k]
			}
		}
		return false
	})
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	outLines := make([]string, len(perms))
	for i, p := range perms {
		pp := make([]string, len(p))
		for j, v := range p {
			pp[j] = fmt.Sprintf("%d", v)
		}
		outLines[i] = strings.Join(pp, " ")
	}
	return strings.Join(parts, " "), strings.Join(outLines, "\n")
}

// 组合总和
func genCombinationSum(rng *rand.Rand) (string, string) {
	// 使用固定候选集+随机目标
	candidateSets := [][]int{
		{2, 3, 6, 7},
		{2, 3, 5},
		{1, 2},
		{2, 5, 10},
	}
	cs := candidateSets[rng.Intn(len(candidateSets))]
	target := rng.Intn(8) + 2
	// 回溯找所有组合
	sort.Ints(cs)
	var combos [][]int
	var bt func(start, remain int, cur []int)
	bt = func(start, remain int, cur []int) {
		if remain == 0 {
			tmp := make([]int, len(cur))
			copy(tmp, cur)
			combos = append(combos, tmp)
			return
		}
		for i := start; i < len(cs); i++ {
			if cs[i] > remain {
				break
			}
			bt(i, remain-cs[i], append(cur, cs[i]))
		}
	}
	bt(0, target, []int{})
	sort.Slice(combos, func(i, j int) bool {
		for k := 0; k < len(combos[i]) && k < len(combos[j]); k++ {
			if combos[i][k] != combos[j][k] {
				return combos[i][k] < combos[j][k]
			}
		}
		return len(combos[i]) < len(combos[j])
	})
	csParts := make([]string, len(cs))
	for i, v := range cs {
		csParts[i] = fmt.Sprintf("%d", v)
	}
	outLines := make([]string, len(combos))
	for i, combo := range combos {
		pp := make([]string, len(combo))
		for j, v := range combo {
			pp[j] = fmt.Sprintf("%d", v)
		}
		outLines[i] = strings.Join(pp, " ")
	}
	return strings.Join(csParts, " ") + "\n" + fmt.Sprintf("%d", target),
		strings.Join(outLines, "\n")
}

// 电话号码的字母组合
func genLetterCombinations(rng *rand.Rand) (string, string) {
	phoneMap := map[byte]string{
		'2': "abc", '3': "def", '4': "ghi", '5': "jkl",
		'6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz",
	}
	digits := "23456789"
	n := rng.Intn(3) + 1
	chosen := make([]byte, n)
	for i := range chosen {
		chosen[i] = digits[rng.Intn(len(digits))]
	}
	input := string(chosen)
	// 生成所有组合
	var result []string
	var bt func(idx int, cur string)
	bt = func(idx int, cur string) {
		if idx == len(input) {
			result = append(result, cur)
			return
		}
		for _, c := range phoneMap[input[idx]] {
			bt(idx+1, cur+string(c))
		}
	}
	bt(0, "")
	sort.Strings(result)
	return input, strings.Join(result, " ")
}

// 括号生成
func genGenerateParenthesis(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	var result []string
	var bt func(open, close int, cur string)
	bt = func(open, close int, cur string) {
		if open == n && close == n {
			result = append(result, cur)
			return
		}
		if open < n {
			bt(open+1, close, cur+"(")
		}
		if close < open {
			bt(open, close+1, cur+")")
		}
	}
	bt(0, 0, "")
	sort.Strings(result)
	return fmt.Sprintf("%d", n), strings.Join(result, " ")
}

// 下一个排列
func genNextPermutation(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(5) + 1
	}
	// 计算下一个排列
	next := make([]int, n)
	copy(next, nums)
	// find i from right where next[i] < next[i+1]
	i := n - 2
	for i >= 0 && next[i] >= next[i+1] {
		i--
	}
	if i >= 0 {
		j := n - 1
		for j >= 0 && next[j] <= next[i] {
			j--
		}
		next[i], next[j] = next[j], next[i]
	}
	// reverse from i+1
	lo, hi := i+1, n-1
	for lo < hi {
		next[lo], next[hi] = next[hi], next[lo]
		lo++
		hi--
	}
	parts := make([]string, n)
	nParts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range next {
		nParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), strings.Join(nParts, " ")
}

// 寻找重复数
func genFindDuplicate(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	dup := rng.Intn(n) + 1
	nums := make([]int, n+1)
	for i := range nums {
		nums[i] = i + 1
		if nums[i] > n {
			nums[i] = n
		}
	}
	// 将一个随机位置的值替换为dup
	pos := rng.Intn(n)
	nums[pos] = dup
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	// 验证确实有重复且只有一个
	freq := map[int]int{}
	for _, v := range nums {
		freq[v]++
	}
	actualDup := dup
	for v, c := range freq {
		if c > 1 {
			actualDup = v
			break
		}
	}
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", actualDup)
}

// 最小覆盖子串
func genMinWindow(rng *rand.Rand) (string, string) {
	type pair struct {
		s, t, ans string
	}
	cases := []pair{
		{"ADOBECODEBANC", "ABC", "BANC"},
		{"a", "a", "a"},
		{"aa", "aa", "aa"},
		{"ab", "a", "a"},
		{"ab", "b", "b"},
		{"AABC", "ABC", "ABC"},
		{"cabwefgewcwaefgcf", "cae", "cwae"},
		{"bba", "ab", "ba"},
	}
	c := cases[rng.Intn(len(cases))]
	return c.s + "\n" + c.t, c.ans
}

// 柱状图中最大的矩形
func genLargestRectangle(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	heights := make([]int, n)
	for i := range heights {
		heights[i] = rng.Intn(8) + 1
	}
	// 单调栈计算
	stack := []int{}
	maxArea := 0
	for i := 0; i <= n; i++ {
		h := 0
		if i < n {
			h = heights[i]
		}
		for len(stack) > 0 && h < heights[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			width := i
			if len(stack) > 0 {
				width = i - stack[len(stack)-1] - 1
			}
			area := heights[top] * width
			if area > maxArea {
				maxArea = area
			}
		}
		stack = append(stack, i)
	}
	parts := make([]string, n)
	for i, v := range heights {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", maxArea)
}

// 最长有效括号
func genLongestValidParentheses(rng *rand.Rand) (string, string) {
	// 随机括号串
	n := rng.Intn(8) + 2
	bs := make([]byte, n)
	for i := range bs {
		if rng.Intn(2) == 0 {
			bs[i] = '('
		} else {
			bs[i] = ')'
		}
	}
	s := string(bs)
	// 计算最长有效括号
	dp := make([]int, n)
	best := 0
	for i := 1; i < n; i++ {
		if s[i] == ')' {
			if s[i-1] == '(' {
				if i >= 2 {
					dp[i] = dp[i-2] + 2
				} else {
					dp[i] = 2
				}
			} else if dp[i-1] > 0 {
				j := i - dp[i-1] - 1
				if j >= 0 && s[j] == '(' {
					dp[i] = dp[i-1] + 2
					if j-1 >= 0 {
						dp[i] += dp[j-1]
					}
				}
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return s, fmt.Sprintf("%d", best)
}
// ===== 第四批生成器 =====

// 排序链表
func genSortList(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(20) - 5
	}
	sorted := make([]int, n)
	copy(sorted, nums)
	sort.Ints(sorted)
	parts := make([]string, n)
	sParts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range sorted {
		sParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), strings.Join(sParts, " ")
}

// K个一组翻转链表
func genReverseKGroup(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	k := rng.Intn(n) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	// 模拟翻转
	result := make([]int, n)
	copy(result, nums)
	for start := 0; start+k <= n; start += k {
		for i, j := start, start+k-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
	}
	parts := make([]string, n)
	rParts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range result {
		rParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", k), strings.Join(rParts, " ")
}

// 随机链表的复制（输出原样）
func genCopyRandomList(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(100)
	}
	// 随机指针索引（-1表示null）
	randoms := make([]int, n)
	for i := range randoms {
		if rng.Intn(3) == 0 {
			randoms[i] = -1
		} else {
			randoms[i] = rng.Intn(n)
		}
	}
	// 构建输入和期望输出（相同格式）
	inputParts := make([]string, n)
	outputParts := make([]string, n)
	for i, v := range vals {
		rStr := "null"
		if randoms[i] >= 0 {
			rStr = fmt.Sprintf("%d", randoms[i])
		}
		inputParts[i] = fmt.Sprintf("%d %s", v, rStr)
		outputParts[i] = fmt.Sprintf("%d %s", v, rStr)
	}
	return strings.Join(inputParts, " "), strings.Join(outputParts, " | ")
}

// LRU缓存
func genLRUCache(rng *rand.Rand) (string, string) {
	cap := rng.Intn(3) + 1
	// 模拟LRU
	type entry struct{ k, v int }
	lru := []entry{}
	ops := []string{}
	results := []string{}
	for i := 0; i < 6+rng.Intn(4); i++ {
		if rng.Intn(2) == 0 || len(lru) == 0 {
			// put
			k := rng.Intn(5) + 1
			v := rng.Intn(10) + 1
			ops = append(ops, fmt.Sprintf("put %d %d", k, v))
			// update lru
			for j, e := range lru {
				if e.k == k {
					lru = append(lru[:j], lru[j+1:]...)
					break
				}
			}
			lru = append(lru, entry{k, v})
			if len(lru) > cap {
				lru = lru[1:]
			}
		} else {
			// get
			k := rng.Intn(5) + 1
			ops = append(ops, fmt.Sprintf("get %d", k))
			found := false
			for j, e := range lru {
				if e.k == k {
					results = append(results, fmt.Sprintf("%d", e.v))
					lru = append(lru[:j], lru[j+1:]...)
					lru = append(lru, e)
					found = true
					break
				}
			}
			if !found {
				results = append(results, "-1")
			}
		}
	}
	return fmt.Sprintf("%d\n%s", cap, strings.Join(ops, "\n")), strings.Join(results, " ")
}

// 二叉树后序遍历
func genPostorderTraversal(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(100)+1)
		}
	}
	result := []string{}
	var postorder func(i int)
	postorder = func(i int) {
		if i >= len(vals) || vals[i] == "null" {
			return
		}
		postorder(2*i + 1)
		postorder(2*i + 2)
		result = append(result, vals[i])
	}
	postorder(0)
	return strings.Join(vals, " "), strings.Join(result, " ")
}

// 二叉树前序遍历
func genPreorderTraversal(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	vals := make([]string, n)
	for i := range vals {
		if rng.Intn(4) == 0 && i > 0 {
			vals[i] = "null"
		} else {
			vals[i] = fmt.Sprintf("%d", rng.Intn(100)+1)
		}
	}
	result := []string{}
	var preorder func(i int)
	preorder = func(i int) {
		if i >= len(vals) || vals[i] == "null" {
			return
		}
		result = append(result, vals[i])
		preorder(2*i + 1)
		preorder(2*i + 2)
	}
	preorder(0)
	return strings.Join(vals, " "), strings.Join(result, " ")
}

// 相交链表（输出交点值或null）
func genGetIntersectionNode(rng *rand.Rand) (string, string) {
	// 生成两条链表，随机决定是否相交
	hasIntersect := rng.Intn(2) == 0
	nA := rng.Intn(4) + 2
	nB := rng.Intn(4) + 2
	aVals := make([]int, nA)
	bVals := make([]int, nB)
	for i := range aVals {
		aVals[i] = rng.Intn(50) + 1
	}
	for i := range bVals {
		bVals[i] = rng.Intn(50) + 1
	}
	if hasIntersect && nA >= 2 && nB >= 2 {
		skipA := rng.Intn(nA - 1)
		skipB := rng.Intn(nB - 1)
		// 公共部分长度 = nA - skipA，确保 B 也能容纳
		commonLen := nA - skipA
		if skipB+commonLen > nB {
			commonLen = nB - skipB
			skipA = nA - commonLen
		}
		if skipA < 0 {
			skipA = 0
			commonLen = nA
			skipB = nB - commonLen
			if skipB < 0 {
				skipB = 0
				commonLen = nA
			}
		}
		aParts := make([]string, nA)
		bParts := make([]string, nB)
		for i, v := range aVals {
			aParts[i] = fmt.Sprintf("%d", v)
		}
		for i := 0; i < skipB; i++ {
			bParts[i] = fmt.Sprintf("%d", bVals[i])
		}
		for i := skipB; i < nB; i++ {
			idx := skipA + (i - skipB)
			if idx >= nA {
				idx = nA - 1
			}
			bParts[i] = fmt.Sprintf("%d", aVals[idx])
		}
		return strings.Join(aParts, " ") + "\n" + strings.Join(bParts, " ") +
			"\n" + fmt.Sprintf("%d %d", skipA, skipB),
			fmt.Sprintf("%d", aVals[skipA])
	}
	aParts := make([]string, nA)
	bParts := make([]string, nB)
	for i, v := range aVals {
		aParts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range bVals {
		bParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(aParts, " ") + "\n" + strings.Join(bParts, " ") +
		"\n-1 -1", "null"
}

// 回文链表
func genIsPalindromeList(rng *rand.Rand) (string, string) {
	isPalin := rng.Intn(2) == 0
	if isPalin {
		n := rng.Intn(4) + 1
		half := make([]int, n)
		for i := range half {
			half[i] = rng.Intn(9) + 1
		}
		full := append(append([]int{}, half...), reverseInts(half)...)
		if rng.Intn(2) == 0 && n > 0 {
			full = append(append([]int{}, half...), append([]int{half[n-1]}, reverseInts(half)...)...)
		}
		parts := make([]string, len(full))
		for i, v := range full {
			parts[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(parts, " "), "true"
	}
	n := rng.Intn(5) + 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(5) + 1
	}
	// check if accidentally palindrome
	isPal := true
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		if vals[i] != vals[j] {
			isPal = false
			break
		}
	}
	if isPal && n > 1 {
		vals[0] = vals[n-1] + 1
		if vals[0] > 9 {
			vals[0] = 1
		}
	}
	parts := make([]string, n)
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	isPal = true
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		if vals[i] != vals[j] {
			isPal = false
			break
		}
	}
	out := "false"
	if isPal {
		out = "true"
	}
	return strings.Join(parts, " "), out
}

func reverseInts(a []int) []int {
	r := make([]int, len(a))
	for i, v := range a {
		r[len(a)-1-i] = v
	}
	return r
}

// 用栈实现队列
func genMyQueue(rng *rand.Rand) (string, string) {
	ops := []string{}
	results := []string{}
	queue := []int{}
	for i := 0; i < 6+rng.Intn(4); i++ {
		if len(queue) == 0 || rng.Intn(2) == 0 {
			v := rng.Intn(10) + 1
			ops = append(ops, fmt.Sprintf("push %d", v))
			queue = append(queue, v)
		} else {
			switch rng.Intn(3) {
			case 0:
				ops = append(ops, "pop")
				results = append(results, fmt.Sprintf("%d", queue[0]))
				queue = queue[1:]
			case 1:
				ops = append(ops, "peek")
				results = append(results, fmt.Sprintf("%d", queue[0]))
			case 2:
				ops = append(ops, "empty")
				if len(queue) == 0 {
					results = append(results, "true")
				} else {
					results = append(results, "false")
				}
			}
		}
	}
	return strings.Join(ops, "\n"), strings.Join(results, " ")
}

// 最小栈
func genMinStack(rng *rand.Rand) (string, string) {
	ops := []string{}
	results := []string{}
	stack := []int{}
	minS := []int{}
	for i := 0; i < 6+rng.Intn(4); i++ {
		if len(stack) == 0 || rng.Intn(2) == 0 {
			v := rng.Intn(20) - 5
			ops = append(ops, fmt.Sprintf("push %d", v))
			stack = append(stack, v)
			if len(minS) == 0 || v <= minS[len(minS)-1] {
				minS = append(minS, v)
			} else {
				minS = append(minS, minS[len(minS)-1])
			}
		} else {
			switch rng.Intn(3) {
			case 0:
				ops = append(ops, "pop")
				stack = stack[:len(stack)-1]
				minS = minS[:len(minS)-1]
			case 1:
				ops = append(ops, "top")
				results = append(results, fmt.Sprintf("%d", stack[len(stack)-1]))
			case 2:
				ops = append(ops, "getMin")
				results = append(results, fmt.Sprintf("%d", minS[len(minS)-1]))
			}
		}
	}
	return strings.Join(ops, "\n"), strings.Join(results, " ")
}

// 二叉搜索树插入操作（中序遍历验证）
func genInsertIntoBST(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i*4 + rng.Intn(3) + 1
	}
	newVal := rng.Intn(n*4+5) + 1
	// ensure unique
	for _, v := range vals {
		if v == newVal {
			newVal += n*4 + 5
		}
	}
	all := append(append([]int{}, vals...), newVal)
	sort.Ints(all)
	treeVals := buildBSTLO(vals)
	strs := make([]string, len(treeVals))
	for i, v := range treeVals {
		strs[i] = fmt.Sprintf("%d", v)
	}
	outParts := make([]string, len(all))
	for i, v := range all {
		outParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, " ") + "\n" + fmt.Sprintf("%d", newVal), strings.Join(outParts, " ")
}

// 删除BST中的节点（中序遍历验证）
func genDeleteNode(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 3
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i*4 + rng.Intn(3) + 1
	}
	keyIdx := rng.Intn(n)
	key := vals[keyIdx]
	remaining := append(append([]int{}, vals[:keyIdx]...), vals[keyIdx+1:]...)
	treeVals := buildBSTLO(vals)
	strs := make([]string, len(treeVals))
	for i, v := range treeVals {
		strs[i] = fmt.Sprintf("%d", v)
	}
	outParts := make([]string, len(remaining))
	for i, v := range remaining {
		outParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, " ") + "\n" + fmt.Sprintf("%d", key), strings.Join(outParts, " ")
}

// 二叉树中的最大路径和
func genMaxPathSum(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	vals := make([]string, n)
	intVals := make([]int, n)
	for i := range vals {
		v := rng.Intn(21) - 10
		vals[i] = fmt.Sprintf("%d", v)
		intVals[i] = v
	}
	best := intVals[0]
	var dfs func(i int) int
	dfs = func(i int) int {
		if i >= len(vals) || vals[i] == "null" {
			return 0
		}
		l := dfs(2*i + 1)
		r := dfs(2*i + 2)
		if l < 0 {
			l = 0
		}
		if r < 0 {
			r = 0
		}
		pathSum := intVals[i] + l + r
		if pathSum > best {
			best = pathSum
		}
		gain := intVals[i] + l
		if r > l {
			gain = intVals[i] + r
		}
		return gain
	}
	dfs(0)
	return strings.Join(vals, " "), fmt.Sprintf("%d", best)
}

// 整数转罗马数字
func genIntToRoman(rng *rand.Rand) (string, string) {
	num := rng.Intn(3999) + 1
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	result := ""
	n := num
	for i, v := range vals {
		for n >= v {
			result += syms[i]
			n -= v
		}
	}
	return fmt.Sprintf("%d", num), result
}

// 罗马数字转整数
func genRomanToInt(rng *rand.Rand) (string, string) {
	num := rng.Intn(3999) + 1
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	roman := ""
	n := num
	for i, v := range vals {
		for n >= v {
			roman += syms[i]
			n -= v
		}
	}
	return roman, fmt.Sprintf("%d", num)
}

// 编辑距离
func genMinDistance(rng *rand.Rand) (string, string) {
	words := []string{"horse", "ros", "intention", "execution", "abc", "abd", "kitten", "sitting", "a", "b", "ab", "ba"}
	w1 := words[rng.Intn(len(words))]
	w2 := words[rng.Intn(len(words))]
	// DP
	m, n2 := len(w1), len(w2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n2+1)
		dp[i][0] = i
	}
	for j := 0; j <= n2; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n2; j++ {
			if w1[i-1] == w2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return w1 + "\n" + w2, fmt.Sprintf("%d", dp[m][n2])
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// 不同的二叉搜索树（卡特兰数）
func genNumTrees(rng *rand.Rand) (string, string) {
	n := rng.Intn(12) + 1
	// 卡特兰数 dp
	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1
	for i := 2; i <= n; i++ {
		for j := 0; j < i; j++ {
			dp[i] += dp[j] * dp[i-1-j]
		}
	}
	return fmt.Sprintf("%d", n), fmt.Sprintf("%d", dp[n])
}

// 分隔链表
func genPartition(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(10) + 1
	}
	x := rng.Intn(8) + 1
	// 模拟分隔
	less := []int{}
	geq := []int{}
	for _, v := range nums {
		if v < x {
			less = append(less, v)
		} else {
			geq = append(geq, v)
		}
	}
	result := append(less, geq...)
	parts := make([]string, n)
	rParts := make([]string, len(result))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	for i, v := range result {
		rParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", x), strings.Join(rParts, " ")
}

// 缺失的第一个正数
func genFirstMissingPositive(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	// 随机决定缺失哪个正数
	missing := rng.Intn(n) + 1
	nums := []int{}
	for i := 1; i <= n+2; i++ {
		if i != missing {
			nums = append(nums, i)
		}
		if len(nums) >= n {
			break
		}
	}
	// 混入一些负数和零
	for i := range nums {
		if rng.Intn(4) == 0 {
			nums[i] = rng.Intn(5) - 3
		}
	}
	// 重新计算真正的缺失值
	posSet := map[int]bool{}
	for _, v := range nums {
		if v > 0 {
			posSet[v] = true
		}
	}
	actual := 1
	for posSet[actual] {
		actual++
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " "), fmt.Sprintf("%d", actual)
}

// ==================== 第五批：最后14道Hot100 ====================

// 和为K的子数组（前缀和 + 哈希表）
func genSubarraySum(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(21) - 10
	}
	k := rng.Intn(21) - 10
	// 前缀和计算
	count := 0
	prefixSum := 0
	prefixMap := map[int]int{0: 1}
	for _, v := range nums {
		prefixSum += v
		if c, ok := prefixMap[prefixSum-k]; ok {
			count += c
		}
		prefixMap[prefixSum]++
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", k), fmt.Sprintf("%d", count)
}

// 滑动窗口最大值（单调队列）
func genMaxSlidingWindow(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := rng.Intn(n) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(41) - 20
	}
	// 单调队列模拟
	deque := []int{}
	result := []int{}
	for i := 0; i < n; i++ {
		for len(deque) > 0 && nums[deque[len(deque)-1]] <= nums[i] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
		if deque[0] <= i-k {
			deque = deque[1:]
		}
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}
	numParts := make([]string, n)
	for i, v := range nums {
		numParts[i] = fmt.Sprintf("%d", v)
	}
	resParts := make([]string, len(result))
	for i, v := range result {
		resParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(numParts, " ") + "\n" + fmt.Sprintf("%d", k), strings.Join(resParts, " ")
}

// 轮转数组
func genRotateArray(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(41) - 20
	}
	k := rng.Intn(n*2) + 1
	kk := k % n
	// 旋转结果
	rotated := make([]int, n)
	for i := 0; i < n; i++ {
		rotated[(i+kk)%n] = nums[i]
	}
	inParts := make([]string, n)
	for i, v := range nums {
		inParts[i] = fmt.Sprintf("%d", v)
	}
	outParts := make([]string, n)
	for i, v := range rotated {
		outParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(inParts, " ") + "\n" + fmt.Sprintf("%d", k), strings.Join(outParts, " ")
}

// 除自身以外数组的乘积
func genProductExceptSelf(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(11) - 3 // [-3, 7]
	}
	result := make([]int, n)
	for i := 0; i < n; i++ {
		prod := 1
		for j := 0; j < n; j++ {
			if j != i {
				prod *= nums[j]
			}
		}
		result[i] = prod
	}
	inParts := make([]string, n)
	for i, v := range nums {
		inParts[i] = fmt.Sprintf("%d", v)
	}
	outParts := make([]string, n)
	for i, v := range result {
		outParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(inParts, " "), strings.Join(outParts, " ")
}

// 搜索二维矩阵II
func genSearchMatrix2(rng *rand.Rand) (string, string) {
	m := rng.Intn(4) + 2
	n := rng.Intn(4) + 2
	matrix := make([][]int, m)
	// 构造每行每列递增的矩阵
	matrix[0] = make([]int, n)
	matrix[0][0] = rng.Intn(5) + 1
	for j := 1; j < n; j++ {
		matrix[0][j] = matrix[0][j-1] + rng.Intn(5) + 1
	}
	for i := 1; i < m; i++ {
		matrix[i] = make([]int, n)
		matrix[i][0] = matrix[i-1][0] + rng.Intn(5) + 1
		for j := 1; j < n; j++ {
			maxPrev := matrix[i][j-1]
			if matrix[i-1][j] > maxPrev {
				maxPrev = matrix[i-1][j]
			}
			matrix[i][j] = maxPrev + rng.Intn(3) + 1
		}
	}
	// 50% 概率搜索存在的值
	var target int
	found := "false"
	if rng.Intn(2) == 0 {
		ri := rng.Intn(m)
		ci := rng.Intn(n)
		target = matrix[ri][ci]
		found = "true"
	} else {
		target = matrix[m-1][n-1] + rng.Intn(10) + 1
	}
	lines := []string{fmt.Sprintf("%d %d", m, n)}
	for i := 0; i < m; i++ {
		rowParts := make([]string, n)
		for j := 0; j < n; j++ {
			rowParts[j] = fmt.Sprintf("%d", matrix[i][j])
		}
		lines = append(lines, strings.Join(rowParts, " "))
	}
	lines = append(lines, fmt.Sprintf("%d", target))
	return strings.Join(lines, "\n"), found
}

// 环形链表II
func genDetectCycle(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(50) + 1
	}
	hasCycle := rng.Intn(2) == 0
	parts := make([]string, n)
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	if hasCycle {
		pos := rng.Intn(n)
		return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", pos), fmt.Sprintf("%d", vals[pos])
	}
	return strings.Join(parts, " ") + "\n-1", "null"
}

// 合并K个升序链表
func genMergeKLists(rng *rand.Rand) (string, string) {
	k := rng.Intn(4) + 1
	allVals := []int{}
	lines := []string{fmt.Sprintf("%d", k)}
	for i := 0; i < k; i++ {
		listLen := rng.Intn(5) + 1
		vals := make([]int, listLen)
		for j := range vals {
			vals[j] = rng.Intn(50) + 1
		}
		sort.Ints(vals)
		allVals = append(allVals, vals...)
		ps := make([]string, listLen)
		for j, v := range vals {
			ps[j] = fmt.Sprintf("%d", v)
		}
		lines = append(lines, strings.Join(ps, " "))
	}
	sort.Ints(allVals)
	if len(allVals) == 0 {
		return strings.Join(lines, "\n"), "empty"
	}
	resParts := make([]string, len(allVals))
	for i, v := range allVals {
		resParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(lines, "\n"), strings.Join(resParts, " ")
}

// 腐烂的橘子（BFS）
func genOrangesRotting(rng *rand.Rand) (string, string) {
	m := rng.Intn(3) + 2
	n := rng.Intn(3) + 2
	grid := make([][]int, m)
	for i := range grid {
		grid[i] = make([]int, n)
		for j := range grid[i] {
			grid[i][j] = rng.Intn(3) // 0, 1, or 2
		}
	}
	// BFS 模拟
	type pos struct{ r, c int }
	queue := []pos{}
	fresh := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 2 {
				queue = append(queue, pos{i, j})
			} else if grid[i][j] == 1 {
				fresh++
			}
		}
	}
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	minutes := 0
	gridCopy := make([][]int, m)
	for i := range grid {
		gridCopy[i] = make([]int, n)
		copy(gridCopy[i], grid[i])
	}
	for len(queue) > 0 && fresh > 0 {
		nextQueue := []pos{}
		for _, p := range queue {
			for _, d := range dirs {
				nr, nc := p.r+d[0], p.c+d[1]
				if nr >= 0 && nr < m && nc >= 0 && nc < n && gridCopy[nr][nc] == 1 {
					gridCopy[nr][nc] = 2
					fresh--
					nextQueue = append(nextQueue, pos{nr, nc})
				}
			}
		}
		if len(nextQueue) > 0 {
			minutes++
		}
		queue = nextQueue
	}
	result := minutes
	if fresh > 0 {
		result = -1
	}
	lines := []string{fmt.Sprintf("%d %d", m, n)}
	for i := 0; i < m; i++ {
		rowParts := make([]string, n)
		for j := 0; j < n; j++ {
			rowParts[j] = fmt.Sprintf("%d", grid[i][j])
		}
		lines = append(lines, strings.Join(rowParts, " "))
	}
	return strings.Join(lines, "\n"), fmt.Sprintf("%d", result)
}

// 单词搜索（回溯）
func genWordExist(rng *rand.Rand) (string, string) {
	m := rng.Intn(3) + 2
	n := rng.Intn(3) + 2
	board := make([][]byte, m)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = byte('A' + rng.Intn(6)) // A-F
		}
	}
	// 50% 概率生成存在的路径
	var word string
	if rng.Intn(2) == 0 {
		// 随机走一条路径
		wLen := rng.Intn(4) + 2
		if wLen > m*n {
			wLen = m*n
		}
		visited := make([][]bool, m)
		for i := range visited {
			visited[i] = make([]bool, n)
		}
		r, c := rng.Intn(m), rng.Intn(n)
		visited[r][c] = true
		path := []byte{board[r][c]}
		dirs := [][2]int{{0,1},{0,-1},{1,0},{-1,0}}
		for len(path) < wLen {
			neighbors := [][2]int{}
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < m && nc >= 0 && nc < n && !visited[nr][nc] {
					neighbors = append(neighbors, [2]int{nr, nc})
				}
			}
			if len(neighbors) == 0 {
				break
			}
			next := neighbors[rng.Intn(len(neighbors))]
			r, c = next[0], next[1]
			visited[r][c] = true
			path = append(path, board[r][c])
		}
		word = string(path)
	} else {
		wLen := rng.Intn(3) + 2
		wb := make([]byte, wLen)
		for i := range wb {
			wb[i] = byte('A' + rng.Intn(8)) // A-H，可能包含不在board上的
		}
		word = string(wb)
	}
	// 回溯验证
	var dfs func(i, j, idx int, vis [][]bool) bool
	dfs = func(i, j, idx int, vis [][]bool) bool {
		if idx == len(word) {
			return true
		}
		if i < 0 || i >= m || j < 0 || j >= n || vis[i][j] || board[i][j] != word[idx] {
			return false
		}
		vis[i][j] = true
		dirs := [][2]int{{0,1},{0,-1},{1,0},{-1,0}}
		for _, d := range dirs {
			if dfs(i+d[0], j+d[1], idx+1, vis) {
				vis[i][j] = false
				return true
			}
		}
		vis[i][j] = false
		return false
	}
	vis := make([][]bool, m)
	for i := range vis {
		vis[i] = make([]bool, n)
	}
	found := false
	for i := 0; i < m && !found; i++ {
		for j := 0; j < n && !found; j++ {
			if dfs(i, j, 0, vis) {
				found = true
			}
		}
	}
	lines := []string{fmt.Sprintf("%d %d", m, n)}
	for i := 0; i < m; i++ {
		rowParts := make([]string, n)
		for j := 0; j < n; j++ {
			rowParts[j] = string(board[i][j])
		}
		lines = append(lines, strings.Join(rowParts, " "))
	}
	lines = append(lines, word)
	return strings.Join(lines, "\n"), fmt.Sprintf("%v", found)
}

// 搜索旋转排序数组
func genSearchRotated(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	sorted := make([]int, n)
	sorted[0] = rng.Intn(10)
	for i := 1; i < n; i++ {
		sorted[i] = sorted[i-1] + rng.Intn(5) + 1
	}
	k := rng.Intn(n)
	rotated := make([]int, n)
	for i := 0; i < n; i++ {
		rotated[i] = sorted[(i+k)%n]
	}
	// 50% 概率搜索存在的值
	var target int
	expected := -1
	if rng.Intn(2) == 0 {
		idx := rng.Intn(n)
		target = rotated[idx]
		expected = idx
	} else {
		target = sorted[n-1] + rng.Intn(10) + 1
	}
	parts := make([]string, n)
	for i, v := range rotated {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", target), fmt.Sprintf("%d", expected)
}

// 在排序数组中查找元素的第一个和最后一个位置
func genSearchRange(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	nums[0] = rng.Intn(5)
	for i := 1; i < n; i++ {
		nums[i] = nums[i-1] + rng.Intn(3)
	}
	// 50% 概率搜索存在的值
	var target int
	first, last := -1, -1
	if rng.Intn(2) == 0 {
		target = nums[rng.Intn(n)]
	} else {
		target = nums[n-1] + rng.Intn(5) + 1
	}
	for i, v := range nums {
		if v == target {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ") + "\n" + fmt.Sprintf("%d", target), fmt.Sprintf("%d %d", first, last)
}

// 每日温度（单调栈）
func genDailyTemperatures(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	temps := make([]int, n)
	for i := range temps {
		temps[i] = rng.Intn(71) + 30 // [30, 100]
	}
	// 单调栈解法
	answer := make([]int, n)
	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && temps[i] > temps[stack[len(stack)-1]] {
			idx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			answer[idx] = i - idx
		}
		stack = append(stack, i)
	}
	inParts := make([]string, n)
	for i, v := range temps {
		inParts[i] = fmt.Sprintf("%d", v)
	}
	outParts := make([]string, n)
	for i, v := range answer {
		outParts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(inParts, " "), strings.Join(outParts, " ")
}

// 字符串解码（递归/栈）
func genDecodeString(rng *rand.Rand) (string, string) {
	// 随机生成编码字符串并计算答案
	type genResult struct {
		encoded string
		decoded string
	}
	var genExpr func(depth int) genResult
	genExpr = func(depth int) genResult {
		result := genResult{}
		parts := rng.Intn(3) + 1
		for p := 0; p < parts; p++ {
			if depth < 2 && rng.Intn(3) > 0 {
				k := rng.Intn(3) + 1
				inner := genExpr(depth + 1)
				result.encoded += fmt.Sprintf("%d[%s]", k, inner.encoded)
				for i := 0; i < k; i++ {
					result.decoded += inner.decoded
				}
			} else {
				cLen := rng.Intn(3) + 1
				s := make([]byte, cLen)
				for i := range s {
					s[i] = byte('a' + rng.Intn(4))
				}
				result.encoded += string(s)
				result.decoded += string(s)
			}
		}
		return result
	}
	r := genExpr(0)
	return r.encoded, r.decoded
}

// 数据流的中位数
func genMedianFinder(rng *rand.Rand) (string, string) {
	ops := rng.Intn(6) + 4
	nums := []int{}
	lines := []string{}
	results := []string{}
	for i := 0; i < ops; i++ {
		if len(nums) == 0 || rng.Intn(3) > 0 {
			v := rng.Intn(201) - 100
			nums = append(nums, v)
			lines = append(lines, fmt.Sprintf("addNum %d", v))
		} else {
			lines = append(lines, "findMedian")
			sorted := make([]int, len(nums))
			copy(sorted, nums)
			sort.Ints(sorted)
			n2 := len(sorted)
			var med float64
			if n2%2 == 1 {
				med = float64(sorted[n2/2])
			} else {
				med = float64(sorted[n2/2-1]+sorted[n2/2]) / 2.0
			}
			results = append(results, fmt.Sprintf("%.2f", med))
		}
	}
	// 保证至少一个findMedian
	if len(results) == 0 {
		lines = append(lines, "findMedian")
		sorted := make([]int, len(nums))
		copy(sorted, nums)
		sort.Ints(sorted)
		n2 := len(sorted)
		var med float64
		if n2%2 == 1 {
			med = float64(sorted[n2/2])
		} else {
			med = float64(sorted[n2/2-1]+sorted[n2/2]) / 2.0
		}
		results = append(results, fmt.Sprintf("%.2f", med))
	}
	return strings.Join(lines, "\n"), strings.Join(results, " ")
}
