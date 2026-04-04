//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/your-org/oj-platform/internal/database"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/pkg/config"
)

// SolutionFunc 输入 -> 期望输出
type SolutionFunc func(input string) string

// solutions 注册表：题目名 -> 解法函数
var solutions = map[string]SolutionFunc{
	// ===== 哈希 =====
	"两数之和":       solveTwoSum,
	"字母异位词分组":    solveGroupAnagrams,
	"最长连续序列":     solveLongestConsecutive,

	// ===== 双指针 =====
	"移动零":        solveMoveZeroes,
	"盛最多水的容器":    solveMaxArea,
	"三数之和":       solveThreeSum,
	"接雨水":        solveTrap,

	// ===== 滑动窗口 =====
	"无重复字符的最长子串":          solveLengthOfLongestSubstring,
	"找到字符串中所有字母异位词": solveFindAnagrams,

	// ===== 子串 =====
	"和为K的子数组":    solveSubarraySum,
	"滑动窗口最大值":    solveMaxSlidingWindow,
	"最小覆盖子串":     solveMinWindow,

	// ===== 普通数组 =====
	"最大子数组和":          solveMaxSubArray,
	"合并区间":            solveMergeIntervals,
	"轮转数组":            solveRotateArray,
	"除自身以外数组的乘积": solveProductExceptSelf,
	"缺失的第一个正数":     solveFirstMissingPositive,

	// ===== 矩阵 =====
	"矩阵置零":       solveSetZeroes,
	"螺旋矩阵":       solveSpiralOrder,
	"旋转图像":       solveRotate,
	"搜索二维矩阵II": solveSearchMatrix2,

	// ===== 链表 =====
	"相交链表":       solveGetIntersectionNode,
	"反转链表":       solveReverseList,
	"回文链表":       solveIsPalindromeList,
	"环形链表":       solveHasCycle,
	"环形链表II":     solveDetectCycle,
	"合并两个有序链表": solveMergeTwoLists,
	"两数相加":       solveAddTwoNumbers,
	"删除链表的倒数第N个节点": solveRemoveNthFromEnd,
	"K个一组翻转链表":  solveReverseKGroup,
	"随机链表的复制":   solveCopyRandomList,
	"排序链表":       solveSortList,
	"合并K个升序链表": solveMergeKLists,
	"LRU缓存":      solveLRUCache,
	"分隔链表":       solvePartition,

	// ===== 二叉树 =====
	"二叉树的中序遍历":          solveInorderTraversal,
	"二叉树的最大深度":          solveMaxDepth,
	"翻转二叉树":              solveInvertTree,
	"对称二叉树":              solveIsSymmetric,
	"二叉树的直径":             solveDiameterOfBinaryTree,
	"二叉树的层序遍历":          solveLevelOrder,
	"验证二叉搜索树":           solveIsValidBST,
	"二叉搜索树中第K小的元素": solveKthSmallest,
	"二叉树的右视图":           solveRightSideView,
	"二叉树展开为链表":          solveFlatten,
	"从前序与中序遍历序列构造二叉树": solveBuildTree,
	"二叉树的最近公共祖先":     solveLowestCommonAncestor,
	"二叉树中的最大路径和":     solveMaxPathSum,
	"二叉树的前序遍历":         solvePreorderTraversal,
	"二叉树的后序遍历":         solvePostorderTraversal,
	"不同的二叉搜索树":         solveNumTrees,
	"二叉搜索树的插入操作":     solveInsertIntoBST,
	"删除二叉搜索树中的节点":   solveDeleteNode,
	"路径总和":               solveHasPathSum,

	// ===== 图论 =====
	"岛屿数量":       solveNumIslands,
	"腐烂的橘子":      solveOrangesRotting,
	"课程表":        solveCanFinish,
	"实现Trie前缀树": solveTrie,

	// ===== 回溯 =====
	"全排列":            solvePermute,
	"子集":              solveSubsets,
	"电话号码的字母组合": solveLetterCombinations,
	"组合总和":          solveCombinationSum,
	"括号生成":          solveGenerateParenthesis,
	"单词搜索":          solveWordExist,
	"全排列II":         solvePermuteUnique,
	"下一个排列":         solveNextPermutation,

	// ===== 二分查找 =====
	"搜索插入位置":                      solveSearchInsert,
	"搜索旋转排序数组":                solveSearchRotated,
	"在排序数组中查找元素的第一个和最后一个位置": solveSearchRange,

	// ===== 栈 =====
	"有效的括号":       solveIsValid,
	"最小栈":          solveMinStack,
	"字符串解码":       solveDecodeString,
	"每日温度":        solveDailyTemperatures,
	"柱状图中最大的矩形": solveLargestRectangle,
	"最长有效括号":     solveLongestValidParentheses,
	"用栈实现队列":     solveMyQueue,

	// ===== 堆 =====
	"数组中的第K个最大元素": solveFindKthLargest,
	"前K个高频元素":      solveTopKFrequent,
	"数据流的中位数":      solveMedianFinder,

	// ===== 贪心 =====
	"买卖股票的最佳时机": solveMaxProfit,
	"跳跃游戏":         solveCanJump,
	"颜色分类":         solveSortColors,

	// ===== 动态规划 =====
	"爬楼梯":          solveClimbStairs,
	"打家劫舍":         solveRob,
	"完全平方数":       solveNumSquares,
	"零钱兑换":         solveCoinChange,
	"单词拆分":         solveWordBreak,
	"最长递增子序列":   solveLengthOfLIS,
	"乘积最大子数组":   solveMaxProduct,
	"不同路径":         solveUniquePaths,
	"最长回文子串":     solveLongestPalindrome,
	"编辑距离":         solveMinDistance,

	// ===== 其他 =====
	"最长公共前缀":     solveLongestCommonPrefix,
	"验证回文串":       solveIsPalindrome,
	"只出现一次的数字": solveSingleNumber,
	"多数元素":         solveMajorityElement,
	"合并两个有序数组": solveMergeSortedArray,
	"寻找重复数":       solveFindDuplicate,
	"整数转罗马数字":   solveIntToRoman,
	"罗马数字转整数":   solveRomanToInt,
}

// ===== 解法实现 =====

// ----- 两数之和 -----
func solveTwoSum(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	target, _ := strconv.Atoi(lines[1])
	m := map[int]int{}
	for i, v := range nums {
		if j, ok := m[target-v]; ok {
			return fmt.Sprintf("%d %d", j, i)
		}
		m[v] = i
	}
	return "-1 -1"
}

// ----- 爬楼梯 -----
func solveClimbStairs(input string) string {
	n, _ := strconv.Atoi(strings.TrimSpace(input))
	a, b := 1, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return fmt.Sprintf("%d", b)
}

// ----- 最大子数组和 -----
func solveMaxSubArray(input string) string {
	nums := parseNums(input)
	maxSum, cur := nums[0], nums[0]
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
	return fmt.Sprintf("%d", maxSum)
}

// ----- 买卖股票的最佳时机 -----
func solveMaxProfit(input string) string {
	prices := parseNums(input)
	maxP, minP := 0, prices[0]
	for _, p := range prices[1:] {
		if p-minP > maxP {
			maxP = p - minP
		}
		if p < minP {
			minP = p
		}
	}
	return fmt.Sprintf("%d", maxP)
}

// ----- 只出现一次的数字 -----
func solveSingleNumber(input string) string {
	nums := parseNums(input)
	xor := 0
	for _, v := range nums {
		xor ^= v
	}
	return fmt.Sprintf("%d", xor)
}

// ----- 多数元素 -----
func solveMajorityElement(input string) string {
	nums := parseNums(input)
	count, candidate := 0, 0
	for _, v := range nums {
		if count == 0 {
			candidate = v
		}
		if v == candidate {
			count++
		} else {
			count--
		}
	}
	return fmt.Sprintf("%d", candidate)
}

// ----- 移动零 -----
func solveMoveZeroes(input string) string {
	nums := parseNums(input)
	res := make([]int, 0, len(nums))
	zeros := 0
	for _, v := range nums {
		if v == 0 {
			zeros++
		} else {
			res = append(res, v)
		}
	}
	for i := 0; i < zeros; i++ {
		res = append(res, 0)
	}
	return joinNums(res)
}

// ----- 合并两个有序数组 -----
func solveMergeSortedArray(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 3 {
		return ""
	}
	// Detect format: old format has 4 lines (nums1, m, nums2, n)
	// New format has 3 lines (m n, nums1, nums2)
	var a, b []int
	firstLineFields := strings.Fields(lines[0])
	if len(firstLineFields) == 2 {
		// New format: "m n" on first line
		m, _ := strconv.Atoi(firstLineFields[0])
		n, _ := strconv.Atoi(firstLineFields[1])
		nums1 := parseNums(lines[1])
		nums2 := parseNums(lines[2])
		if m > len(nums1) {
			m = len(nums1)
		}
		if n > len(nums2) {
			n = len(nums2)
		}
		a = nums1[:m]
		b = nums2[:n]
	} else if len(lines) >= 4 {
		// Old format: nums1, m, nums2, n
		nums1 := parseNums(lines[0])
		m, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
		nums2 := parseNums(lines[2])
		n, _ := strconv.Atoi(strings.TrimSpace(lines[3]))
		if m > len(nums1) {
			m = len(nums1)
		}
		if n > len(nums2) {
			n = len(nums2)
		}
		a = nums1[:m]
		b = nums2[:n]
	} else {
		return ""
	}
	merged := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			merged = append(merged, a[i])
			i++
		} else {
			merged = append(merged, b[j])
			j++
		}
	}
	merged = append(merged, a[i:]...)
	merged = append(merged, b[j:]...)
	return joinNums(merged)
}

// ----- 验证回文串 -----
func solveIsPalindrome(input string) string {
	s := strings.ToLower(input)
	i, j := 0, len(s)-1
	for i < j {
		for i < j && !isAlphaNum(s[i]) {
			i++
		}
		for i < j && !isAlphaNum(s[j]) {
			j--
		}
		if s[i] != s[j] {
			return "false"
		}
		i++
		j--
	}
	return "true"
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

// ----- 找到字符串中所有字母异位词 -----
func solveFindAnagrams(input string) string {
	lines := strings.Split(input, "\n")
	s, p := lines[0], lines[1]
	if len(s) < len(p) {
		return "none"
	}
	pCount := [26]int{}
	for _, c := range p {
		pCount[c-'a']++
	}
	sCount := [26]int{}
	res := []int{}
	for i := 0; i < len(s); i++ {
		sCount[s[i]-'a']++
		if i >= len(p) {
			sCount[s[i-len(p)]-'a']--
		}
		if i >= len(p)-1 && sCount == pCount {
			res = append(res, i-len(p)+1)
		}
	}
	if len(res) == 0 {
		return "none"
	}
	strs := make([]string, len(res))
	for i, v := range res {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, " ")
}

// ----- 二叉树的最大深度 -----
func solveMaxDepth(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	if root == nil {
		return "0"
	}
	var depth func(n *treeNode) int
	depth = func(n *treeNode) int {
		if n == nil {
			return 0
		}
		l, r := depth(n.left), depth(n.right)
		if l > r {
			return l + 1
		}
		return r + 1
	}
	return fmt.Sprintf("%d", depth(root))
}

// ----- 二叉树的中序遍历 -----
func solveInorderTraversal(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	result := []string{}
	var inorder func(n *treeNode)
	inorder = func(n *treeNode) {
		if n == nil {
			return
		}
		inorder(n.left)
		result = append(result, strconv.Itoa(n.val))
		inorder(n.right)
	}
	inorder(root)
	return strings.Join(result, " ")
}

// ----- 对称二叉树 -----
func solveIsSymmetric(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	if root == nil {
		return "true"
	}
	var isMirror func(l, r *treeNode) bool
	isMirror = func(l, r *treeNode) bool {
		if l == nil && r == nil {
			return true
		}
		if l == nil || r == nil {
			return false
		}
		return l.val == r.val && isMirror(l.left, r.right) && isMirror(l.right, r.left)
	}
	if isMirror(root.left, root.right) {
		return "true"
	}
	return "false"
}

// ----- 有效的括号 -----
func solveIsValid(input string) string {
	s := strings.TrimSpace(input)
	stack := []byte{}
	pairs := map[byte]byte{')': '(', ']': '[', '}': '{'}
	for _, c := range s {
		if c == '(' || c == '[' || c == '{' {
			stack = append(stack, byte(c))
		} else {
			if len(stack) == 0 || stack[len(stack)-1] != pairs[byte(c)] {
				return "false"
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) == 0 {
		return "true"
	}
	return "false"
}

// ----- 最长公共前缀 -----
func solveLongestCommonPrefix(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	// 输入是空格分隔的单词，单行
	strs := strings.Fields(input)
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

// ----- 反转链表 -----
func solveReverseList(input string) string {
	input = strings.TrimSpace(input)
	if input == "null" || input == "" {
		return "null"
	}
	nums := parseNums(input)
	// 反转
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return joinNums(nums)
}

// ----- 无重复字符的最长子串 -----
func solveLengthOfLongestSubstring(input string) string {
	s := strings.TrimSpace(input)
	charIdx := map[byte]int{}
	maxLen, start := 0, 0
	for i := 0; i < len(s); i++ {
		if idx, ok := charIdx[s[i]]; ok && idx >= start {
			start = idx + 1
		}
		charIdx[s[i]] = i
		if i-start+1 > maxLen {
			maxLen = i - start + 1
		}
	}
	return fmt.Sprintf("%d", maxLen)
}

// ----- 盛最多水的容器 -----
func solveMaxArea(input string) string {
	height := parseNums(input)
	i, j := 0, len(height)-1
	maxA := 0
	for i < j {
		h := height[i]
		if height[j] < h {
			h = height[j]
		}
		area := h * (j - i)
		if area > maxA {
			maxA = area
		}
		if height[i] < height[j] {
			i++
		} else {
			j--
		}
	}
	return fmt.Sprintf("%d", maxA)
}

// ----- 三数之和 -----
func solveThreeSum(input string) string {
	nums := parseNums(input)
	sort.Ints(nums)
	res := [][3]int{}
	for i := 0; i < len(nums)-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		j, k := i+1, len(nums)-1
		for j < k {
			sum := nums[i] + nums[j] + nums[k]
			if sum == 0 {
				res = append(res, [3]int{nums[i], nums[j], nums[k]})
				for j < k && nums[j] == nums[j+1] {
					j++
				}
				for j < k && nums[k] == nums[k-1] {
					k--
				}
				j++
				k--
			} else if sum < 0 {
				j++
			} else {
				k--
			}
		}
	}
	// 格式化输出
	lines := make([]string, len(res))
	for i, t := range res {
		lines[i] = fmt.Sprintf("%d %d %d", t[0], t[1], t[2])
	}
	return strings.Join(lines, "\n")
}

// ----- 翻转二叉树 -----
func solveInvertTree(input string) string {
	input = strings.TrimSpace(input)
	root := parseTreeBFS(input)
	var invert func(n *treeNode)
	invert = func(n *treeNode) {
		if n == nil {
			return
		}
		n.left, n.right = n.right, n.left
		invert(n.left)
		invert(n.right)
	}
	invert(root)
	return serializeTreeBFS(root)
}

// ----- 二叉树的层序遍历 -----
func solveLevelOrder(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	if root == nil {
		return ""
	}
	levels := []string{}
	queue := []*treeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		levelNodes := []string{}
		for i := 0; i < size; i++ {
			cur := queue[0]
			queue = queue[1:]
			levelNodes = append(levelNodes, strconv.Itoa(cur.val))
			if cur.left != nil {
				queue = append(queue, cur.left)
			}
			if cur.right != nil {
				queue = append(queue, cur.right)
			}
		}
		levels = append(levels, strings.Join(levelNodes, " "))
	}
	return strings.Join(levels, "\n")
}

// ----- 最长回文子串 -----
func solveLongestPalindrome(input string) string {
	s := strings.TrimSpace(input)
	if len(s) < 2 {
		return s
	}
	expand := func(left, right int) string {
		for left >= 0 && right < len(s) && s[left] == s[right] {
			left--
			right++
		}
		return s[left+1 : right]
	}
	result := ""
	for i := 0; i < len(s); i++ {
		s1 := expand(i, i)
		s2 := expand(i, i+1)
		if len(s1) > len(result) {
			result = s1
		}
		if len(s2) > len(result) {
			result = s2
		}
	}
	return result
}

// ----- 接雨水 -----
func solveTrap(input string) string {
	height := parseNums(input)
	if len(height) < 3 {
		return "0"
	}
	left, right := 0, len(height)-1
	leftMax, rightMax := height[left], height[right]
	water := 0
	for left < right {
		if leftMax < rightMax {
			left++
			if height[left] > leftMax {
				leftMax = height[left]
			} else {
				water += leftMax - height[left]
			}
		} else {
			right--
			if height[right] > rightMax {
				rightMax = height[right]
			} else {
				water += rightMax - height[right]
			}
		}
	}
	return fmt.Sprintf("%d", water)
}

// ----- 跳跃游戏 -----
func solveCanJump(input string) string {
	nums := parseNums(input)
	maxReach := 0
	for i, v := range nums {
		if i > maxReach {
			return "false"
		}
		if i+v > maxReach {
			maxReach = i + v
		}
		if maxReach >= len(nums)-1 {
			return "true"
		}
	}
	return "true"
}

// ----- 不同路径 -----
func solveUniquePaths(input string) string {
	parts := strings.Fields(input)
	m, _ := strconv.Atoi(parts[0])
	n, _ := strconv.Atoi(parts[1])
	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[j] += dp[j-1]
		}
	}
	return fmt.Sprintf("%d", dp[n-1])
}

// ----- 搜索插入位置 -----
func solveSearchInsert(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	target, _ := strconv.Atoi(lines[1])
	left, right := 0, len(nums)
	for left < right {
		mid := (left + right) / 2
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return fmt.Sprintf("%d", left)
}

// ----- 合并区间 -----
func solveMergeIntervals(input string) string {
	lines := strings.Split(input, "\n")
	intervals := make([][2]int, len(lines))
	for i, line := range lines {
		parts := strings.Fields(line)
		intervals[i][0], _ = strconv.Atoi(parts[0])
		intervals[i][1], _ = strconv.Atoi(parts[1])
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	merged := [][2]int{intervals[0]}
	for _, iv := range intervals[1:] {
		last := &merged[len(merged)-1]
		if iv[0] <= last[1] {
			if iv[1] > last[1] {
				last[1] = iv[1]
			}
		} else {
			merged = append(merged, iv)
		}
	}
	result := make([]string, len(merged))
	for i, iv := range merged {
		result[i] = fmt.Sprintf("%d %d", iv[0], iv[1])
	}
	return strings.Join(result, "\n")
}

// ----- 合并两个有序链表 -----
func solveMergeTwoLists(input string) string {
	lines := strings.Split(input, "\n")
	var l1, l2 []int
	if len(lines) > 0 && lines[0] != "" {
		l1 = parseNums(lines[0])
	}
	if len(lines) > 1 && lines[1] != "" {
		l2 = parseNums(lines[1])
	}
	merged := make([]int, 0, len(l1)+len(l2))
	i, j := 0, 0
	for i < len(l1) && j < len(l2) {
		if l1[i] <= l2[j] {
			merged = append(merged, l1[i])
			i++
		} else {
			merged = append(merged, l2[j])
			j++
		}
	}
	merged = append(merged, l1[i:]...)
	merged = append(merged, l2[j:]...)
	if len(merged) == 0 {
		return "null"
	}
	return joinNums(merged)
}

// ----- 环形链表 -----
func solveHasCycle(input string) string {
	lines := strings.Split(input, "\n")
	pos, _ := strconv.Atoi(lines[1])
	if pos >= 0 {
		return "true"
	}
	return "false"
}

// ----- 全排列 -----
func solvePermute(input string) string {
	nums := parseNums(input)
	result := [][]int{}
	var backtrack func(start int)
	backtrack = func(start int) {
		if start == len(nums) {
			tmp := make([]int, len(nums))
			copy(tmp, nums)
			result = append(result, tmp)
			return
		}
		for i := start; i < len(nums); i++ {
			nums[start], nums[i] = nums[i], nums[start]
			backtrack(start + 1)
			nums[start], nums[i] = nums[i], nums[start]
		}
	}
	backtrack(0)
	// 按字典序排序
	sort.Slice(result, func(i, j int) bool {
		for k := range result[i] {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return false
	})
	lines := make([]string, len(result))
	for i, perm := range result {
		lines[i] = joinNums(perm)
	}
	return strings.Join(lines, "\n")
}

// ----- 删除链表的倒数第N个节点 -----
func solveRemoveNthFromEnd(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	n, _ := strconv.Atoi(lines[1])
	// 删除倒数第 n 个
	idx := len(nums) - n
	if idx < 0 {
		idx = 0
	}
	result := append(nums[:idx], nums[idx+1:]...)
	if len(result) == 0 {
		return "null"
	}
	return joinNums(result)
}

// ----- 路径总和 -----
func solveHasPathSum(input string) string {
	lines := strings.Split(input, "\n")
	nodes := strings.Fields(lines[0])
	target, _ := strconv.Atoi(lines[1])
	// 空树：没有任何路径，返回 false
	if len(nodes) == 0 || nodes[0] == "null" {
		return "false"
	}
	var dfs func(idx, sum int) bool
	dfs = func(idx, sum int) bool {
		if idx >= len(nodes) || nodes[idx] == "null" {
			return false
		}
		v, _ := strconv.Atoi(nodes[idx])
		sum += v
		left := 2*idx + 1
		right := 2*idx + 2
		isLeaf := (left >= len(nodes) || nodes[left] == "null") && (right >= len(nodes) || nodes[right] == "null")
		if isLeaf {
			return sum == target
		}
		return dfs(left, sum) || dfs(right, sum)
	}
	if dfs(0, 0) {
		return "true"
	}
	return "false"
}

// ----- 子集 -----
func solveSubsets(input string) string {
	nums := parseNums(input)
	sort.Ints(nums)
	result := [][]int{{}}
	for _, v := range nums {
		size := len(result)
		for i := 0; i < size; i++ {
			tmp := make([]int, len(result[i]))
			copy(tmp, result[i])
			tmp = append(tmp, v)
			result = append(result, tmp)
		}
	}
	// 按长度排序，长度相同按字典序
	sort.Slice(result, func(i, j int) bool {
		if len(result[i]) != len(result[j]) {
			return len(result[i]) < len(result[j])
		}
		for k := range result[i] {
			if result[i][k] != result[j][k] {
				return result[i][k] < result[j][k]
			}
		}
		return false
	})
	lines := make([]string, len(result))
	for i, subset := range result {
		if len(subset) == 0 {
			lines[i] = ""
		} else {
			lines[i] = joinNums(subset)
		}
	}
	return strings.Join(lines, "\n")
}

// ----- 岛屿数量 -----
func solveNumIslands(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		return "0"
	}
	// Check if first line is "m n" header (two integers)
	startLine := 0
	firstFields := strings.Fields(lines[0])
	if len(firstFields) == 2 {
		_, e1 := strconv.Atoi(firstFields[0])
		_, e2 := strconv.Atoi(firstFields[1])
		if e1 == nil && e2 == nil {
			// Verify second line is not also "m n" style - check if it looks like grid data
			if len(lines) > 1 {
				secondFields := strings.Fields(lines[1])
				if len(secondFields) > 2 || (len(secondFields) == 1 && len(secondFields[0]) > 1) {
					startLine = 1
				}
			}
		}
	}
	grid := make([][]byte, 0)
	for i := startLine; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 1 {
			// Space separated
			row := make([]byte, len(fields))
			for j, f := range fields {
				row[j] = f[0]
			}
			grid = append(grid, row)
		} else {
			// Compact string
			grid = append(grid, []byte(line))
		}
	}
	if len(grid) == 0 {
		return "0"
	}
	m := len(grid)
	n := len(grid[0])
	count := 0
	var dfs func(i, j int)
	dfs = func(i, j int) {
		if i < 0 || i >= m || j < 0 || j >= n || grid[i][j] != '1' {
			return
		}
		grid[i][j] = '0'
		dfs(i+1, j)
		dfs(i-1, j)
		dfs(i, j+1)
		dfs(i, j-1)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				count++
				dfs(i, j)
			}
		}
	}
	return fmt.Sprintf("%d", count)
}

// ===== 辅助函数 =====

func parseNums(s string) []int {
	parts := strings.Fields(s)
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		if p == "" || p == "null" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

// treeNode is a proper binary tree node
type treeNode struct {
	val   int
	left  *treeNode
	right *treeNode
}

// parseTreeBFS parses a level-order BFS tree string like "1 null 2 3"
// using the proper LeetCode BFS approach (nulls don't have children)
func parseTreeBFS(s string) *treeNode {
	s = strings.TrimSpace(s)
	tokens := strings.Fields(s)
	if len(tokens) == 0 || tokens[0] == "null" {
		return nil
	}
	rootVal, _ := strconv.Atoi(tokens[0])
	root := &treeNode{val: rootVal}
	queue := []*treeNode{root}
	i := 1
	for len(queue) > 0 && i < len(tokens) {
		cur := queue[0]
		queue = queue[1:]
		// left child
		if i < len(tokens) {
			if tokens[i] != "null" {
				v, _ := strconv.Atoi(tokens[i])
				cur.left = &treeNode{val: v}
				queue = append(queue, cur.left)
			}
			i++
		}
		// right child
		if i < len(tokens) {
			if tokens[i] != "null" {
				v, _ := strconv.Atoi(tokens[i])
				cur.right = &treeNode{val: v}
				queue = append(queue, cur.right)
			}
			i++
		}
	}
	return root
}

// serializeTreeBFS converts tree to level-order string, removing trailing nulls
func serializeTreeBFS(root *treeNode) string {
	if root == nil {
		return "null"
	}
	result := []string{}
	queue := []*treeNode{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil {
			result = append(result, "null")
		} else {
			result = append(result, strconv.Itoa(cur.val))
			queue = append(queue, cur.left, cur.right)
		}
	}
	// Remove trailing nulls
	for len(result) > 0 && result[len(result)-1] == "null" {
		result = result[:len(result)-1]
	}
	return strings.Join(result, " ")
}

func joinNums(nums []int) string {
	strs := make([]string, len(nums))
	for i, v := range nums {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, " ")
}

// ===== 剩余题目占位实现 =====

func solveGroupAnagrams(input string) string {
	input = strings.TrimSpace(input)
	words := strings.Fields(input)
	groups := map[string][]string{}
	for _, w := range words {
		sorted := sortString(w)
		groups[sorted] = append(groups[sorted], w)
	}
	// Sort keys alphabetically for determinism
	keys := []string{}
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := []string{}
	for _, k := range keys {
		// Sort words within group alphabetically
		g := groups[k]
		sort.Strings(g)
		result = append(result, strings.Join(g, " "))
	}
	return strings.Join(result, "\n")
}

func sortString(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}

func solveLongestConsecutive(input string) string {
	nums := parseNums(input)
	if len(nums) == 0 {
		return "0"
	}
	numSet := make(map[int]bool)
	for _, n := range nums {
		numSet[n] = true
	}
	maxLen := 0
	for n := range numSet {
		if !numSet[n-1] {
			cur := n
			length := 1
			for numSet[cur+1] {
				cur++
				length++
			}
			if length > maxLen {
				maxLen = length
			}
		}
	}
	return fmt.Sprintf("%d", maxLen)
}

func solveSubarraySum(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	k, _ := strconv.Atoi(lines[1])
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
	return fmt.Sprintf("%d", count)
}

func solveMaxSlidingWindow(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	k, _ := strconv.Atoi(lines[1])
	if len(nums) == 0 || k == 0 {
		return ""
	}
	deque := []int{}
	result := []int{}
	for i := 0; i < len(nums); i++ {
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
	return joinNums(result)
}

func solveMinWindow(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	s := strings.TrimSpace(lines[0])
	t := strings.TrimSpace(lines[1])
	if len(s) == 0 || len(t) == 0 {
		return ""
	}
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	window := make(map[byte]int)
	left, right := 0, 0
	valid := 0
	start, minLen := 0, len(s)+1
	for right < len(s) {
		c := s[right]
		right++
		if need[c] > 0 {
			window[c]++
			if window[c] == need[c] {
				valid++
			}
		}
		for valid == len(need) {
			if right-left < minLen {
				start = left
				minLen = right - left
			}
			d := s[left]
			left++
			if need[d] > 0 {
				if window[d] == need[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	if minLen == len(s)+1 {
		return ""
	}
	return s[start : start+minLen]
}
func solveRotateArray(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	k, _ := strconv.Atoi(lines[1])
	n := len(nums)
	if n == 0 {
		return ""
	}
	k = k % n
	if k == 0 {
		return joinNums(nums)
	}
	rotated := make([]int, n)
	for i := 0; i < n; i++ {
		rotated[(i+k)%n] = nums[i]
	}
	return joinNums(rotated)
}
func solveProductExceptSelf(input string) string {
	nums := parseNums(input)
	n := len(nums)
	if n == 0 {
		return ""
	}
	result := make([]int, n)
	left := 1
	for i := 0; i < n; i++ {
		result[i] = left
		left *= nums[i]
	}
	right := 1
	for i := n - 1; i >= 0; i-- {
		result[i] *= right
		right *= nums[i]
	}
	return joinNums(result)
}
func solveFirstMissingPositive(input string) string {
	nums := parseNums(input)
	n := len(nums)
	for i := 0; i < n; i++ {
		for nums[i] > 0 && nums[i] <= n && nums[nums[i]-1] != nums[i] {
			nums[i], nums[nums[i]-1] = nums[nums[i]-1], nums[i]
		}
	}
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return fmt.Sprintf("%d", i+1)
		}
	}
	return fmt.Sprintf("%d", n+1)
}
func solveSetZeroes(input string) string {
	lines := strings.Split(input, "\n")
	m := len(lines)
	if m == 0 {
		return ""
	}
	matrix := make([][]int, m)
	for i := 0; i < m; i++ {
		matrix[i] = parseNums(lines[i])
	}
	n := len(matrix[0])
	rows := make(map[int]bool)
	cols := make(map[int]bool)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == 0 {
				rows[i] = true
				cols[j] = true
			}
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if rows[i] || cols[j] {
				matrix[i][j] = 0
			}
		}
	}
	result := make([]string, m)
	for i := 0; i < m; i++ {
		result[i] = joinNums(matrix[i])
	}
	return strings.Join(result, "\n")
}
func solveSpiralOrder(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	m := len(lines)
	if m == 0 {
		return ""
	}
	matrix := make([][]int, m)
	for i := 0; i < m; i++ {
		matrix[i] = parseNums(lines[i])
	}
	n := len(matrix[0])
	result := []int{}
	top, bottom, left, right := 0, m-1, 0, n-1
	for top <= bottom && left <= right {
		for j := left; j <= right; j++ {
			result = append(result, matrix[top][j])
		}
		top++
		for i := top; i <= bottom; i++ {
			result = append(result, matrix[i][right])
		}
		right--
		if top <= bottom {
			for j := right; j >= left; j-- {
				result = append(result, matrix[bottom][j])
			}
			bottom--
		}
		if left <= right {
			for i := bottom; i >= top; i-- {
				result = append(result, matrix[i][left])
			}
			left++
		}
	}
	return joinNums(result)
}
func solveRotate(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	n := len(lines)
	if n == 0 {
		return ""
	}
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = parseNums(lines[i])
	}
	// 先转置
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
	// 再水平翻转
	for i := 0; i < n; i++ {
		for j := 0; j < n/2; j++ {
			matrix[i][j], matrix[i][n-1-j] = matrix[i][n-1-j], matrix[i][j]
		}
	}
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = joinNums(matrix[i])
	}
	return strings.Join(result, "\n")
}
func solveSearchMatrix2(input string) string {
	lines := strings.Split(input, "\n")
	dims := strings.Fields(lines[0])
	m, _ := strconv.Atoi(dims[0])
	n, _ := strconv.Atoi(dims[1])
	matrix := make([][]int, m)
	for i := 0; i < m; i++ {
		matrix[i] = parseNums(lines[i+1])
	}
	target, _ := strconv.Atoi(lines[m+1])
	i, j := 0, n-1
	for i < m && j >= 0 {
		if matrix[i][j] == target {
			return "true"
		} else if matrix[i][j] > target {
			j--
		} else {
			i++
		}
	}
	return "false"
}
func solveGetIntersectionNode(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 3 {
		return "null"
	}
	numsA := parseNums(lines[0])
	offsets := strings.Fields(lines[2])
	if len(offsets) < 1 {
		return "null"
	}
	posA, _ := strconv.Atoi(offsets[0])
	if posA < 0 || posA >= len(numsA) {
		return "null"
	}
	return fmt.Sprintf("%d", numsA[posA])
}
func solveIsPalindromeList(input string) string {
	nums := parseNums(input)
	n := len(nums)
	for i := 0; i < n/2; i++ {
		if nums[i] != nums[n-1-i] {
			return "false"
		}
	}
	return "true"
}
func solveDetectCycle(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return "null"
	}
	pos, _ := strconv.Atoi(lines[1])
	if pos < 0 {
		return "null"
	}
	nums := parseNums(lines[0])
	if pos < len(nums) {
		return fmt.Sprintf("%d", nums[pos])
	}
	return "null"
}
func solveAddTwoNumbers(input string) string {
	lines := strings.Split(input, "\n")
	l1 := parseNums(lines[0])
	l2 := parseNums(lines[1])
	result := []int{}
	carry := 0
	i, j := 0, 0
	for i < len(l1) || j < len(l2) || carry > 0 {
		sum := carry
		if i < len(l1) {
			sum += l1[i]
			i++
		}
		if j < len(l2) {
			sum += l2[j]
			j++
		}
		result = append(result, sum%10)
		carry = sum / 10
	}
	return joinNums(result)
}
func solveReverseKGroup(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	nums := parseNums(lines[0])
	k, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	if k <= 1 {
		return joinNums(nums)
	}
	result := []int{}
	for i := 0; i+k <= len(nums); i += k {
		group := nums[i : i+k]
		for j := len(group) - 1; j >= 0; j-- {
			result = append(result, group[j])
		}
	}
	rem := len(nums) % k
	if rem > 0 {
		result = append(result, nums[len(nums)-rem:]...)
	}
	return joinNums(result)
}
func solveCopyRandomList(input string) string {
	input = strings.TrimSpace(input)
	// Output format uses "|" as separator between nodes
	// Input format: "val random_idx val random_idx ..."
	// Output format: "val random_idx | val random_idx | ..."
	tokens := strings.Fields(input)
	if len(tokens) == 0 || input == "null" {
		return "null"
	}
	nodes := []string{}
	for i := 0; i+1 < len(tokens); i += 2 {
		nodes = append(nodes, tokens[i]+" "+tokens[i+1])
	}
	return strings.Join(nodes, " | ")
}
func solveSortList(input string) string {
	nums := parseNums(input)
	sort.Ints(nums)
	return joinNums(nums)
}
func solveMergeKLists(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return "empty"
	}
	k, _ := strconv.Atoi(lines[0])
	allNums := []int{}
	for i := 1; i <= k; i++ {
		if i < len(lines) {
			nums := parseNums(lines[i])
			allNums = append(allNums, nums...)
		}
	}
	if len(allNums) == 0 {
		return "empty"
	}
	sort.Ints(allNums)
	return joinNums(allNums)
}
func solveLRUCache(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return ""
	}
	capacity, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	cache := make(map[int]int)
	order := []int{} // front = most recent
	evict := func(key int) {
		for i, k := range order {
			if k == key {
				order = append(order[:i], order[i+1:]...)
				return
			}
		}
	}
	results := []string{}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if parts[0] == "get" {
			key, _ := strconv.Atoi(parts[1])
			if val, ok := cache[key]; ok {
				evict(key)
				order = append([]int{key}, order...)
				results = append(results, strconv.Itoa(val))
			} else {
				results = append(results, "-1")
			}
		} else if parts[0] == "put" {
			key, _ := strconv.Atoi(parts[1])
			val, _ := strconv.Atoi(parts[2])
			if _, ok := cache[key]; ok {
				evict(key)
			} else if len(cache) >= capacity {
				// evict least recently used (end of order)
				lru := order[len(order)-1]
				order = order[:len(order)-1]
				delete(cache, lru)
			}
			cache[key] = val
			order = append([]int{key}, order...)
		}
	}
	return strings.Join(results, " ")
}
func solvePartition(input string) string {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	x, _ := strconv.Atoi(lines[1])
	less, geq := []int{}, []int{}
	for _, v := range nums {
		if v < x {
			less = append(less, v)
		} else {
			geq = append(geq, v)
		}
	}
	return joinNums(append(less, geq...))
}
func solveDiameterOfBinaryTree(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	if root == nil {
		return "0"
	}
	maxDiam := 0
	var depth func(n *treeNode) int
	depth = func(n *treeNode) int {
		if n == nil {
			return 0
		}
		l := depth(n.left)
		r := depth(n.right)
		if l+r > maxDiam {
			maxDiam = l + r
		}
		if l > r {
			return l + 1
		}
		return r + 1
	}
	depth(root)
	return fmt.Sprintf("%d", maxDiam)
}
func solveIsValidBST(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	var validate func(n *treeNode, min, max int) bool
	validate = func(n *treeNode, min, max int) bool {
		if n == nil {
			return true
		}
		if n.val <= min || n.val >= max {
			return false
		}
		return validate(n.left, min, n.val) && validate(n.right, n.val, max)
	}
	if validate(root, -1<<62, 1<<62) {
		return "true"
	}
	return "false"
}
func solveKthSmallest(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	root := parseTreeBFS(lines[0])
	k, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	result := []int{}
	var inorder func(n *treeNode)
	inorder = func(n *treeNode) {
		if n == nil {
			return
		}
		inorder(n.left)
		result = append(result, n.val)
		inorder(n.right)
	}
	inorder(root)
	if k <= len(result) {
		return fmt.Sprintf("%d", result[k-1])
	}
	return ""
}
func solveRightSideView(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	if root == nil {
		return ""
	}
	result := []string{}
	queue := []*treeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			cur := queue[0]
			queue = queue[1:]
			if i == size-1 {
				result = append(result, strconv.Itoa(cur.val))
			}
			if cur.left != nil {
				queue = append(queue, cur.left)
			}
			if cur.right != nil {
				queue = append(queue, cur.right)
			}
		}
	}
	return strings.Join(result, " ")
}
func solveFlatten(input string) string {
	input = strings.TrimSpace(input)
	root := parseTreeBFS(input)
	result := []string{}
	var preorder func(n *treeNode)
	preorder = func(n *treeNode) {
		if n == nil {
			return
		}
		result = append(result, strconv.Itoa(n.val))
		preorder(n.left)
		preorder(n.right)
	}
	preorder(root)
	return strings.Join(result, " ")
}
func solveBuildTree(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	preorder := parseNums(lines[0])
	inorder := parseNums(lines[1])
	if len(preorder) == 0 {
		return "null"
	}
	var build func(pre, ino []int) *treeNode
	build = func(pre, ino []int) *treeNode {
		if len(pre) == 0 {
			return nil
		}
		root := &treeNode{val: pre[0]}
		mid := 0
		for i, v := range ino {
			if v == pre[0] {
				mid = i
				break
			}
		}
		root.left = build(pre[1:mid+1], ino[:mid])
		root.right = build(pre[mid+1:], ino[mid+1:])
		return root
	}
	root := build(preorder, inorder)
	return serializeTreeBFS(root)
}
func solveLowestCommonAncestor(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	root := parseTreeBFS(lines[0])
	pq := strings.Fields(lines[1])
	p, _ := strconv.Atoi(pq[0])
	q, _ := strconv.Atoi(pq[1])
	var lca func(n *treeNode) *treeNode
	lca = func(n *treeNode) *treeNode {
		if n == nil || n.val == p || n.val == q {
			return n
		}
		l := lca(n.left)
		r := lca(n.right)
		if l != nil && r != nil {
			return n
		}
		if l != nil {
			return l
		}
		return r
	}
	result := lca(root)
	if result == nil {
		return ""
	}
	return strconv.Itoa(result.val)
}
func solveMaxPathSum(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	if root == nil {
		return "0"
	}
	maxSum := -1 << 62
	var gain func(n *treeNode) int
	gain = func(n *treeNode) int {
		if n == nil {
			return 0
		}
		l := gain(n.left)
		r := gain(n.right)
		if l < 0 {
			l = 0
		}
		if r < 0 {
			r = 0
		}
		if n.val+l+r > maxSum {
			maxSum = n.val + l + r
		}
		if l > r {
			return n.val + l
		}
		return n.val + r
	}
	gain(root)
	return fmt.Sprintf("%d", maxSum)
}
func solvePreorderTraversal(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	result := []string{}
	var preorder func(n *treeNode)
	preorder = func(n *treeNode) {
		if n == nil {
			return
		}
		result = append(result, strconv.Itoa(n.val))
		preorder(n.left)
		preorder(n.right)
	}
	preorder(root)
	return strings.Join(result, " ")
}
func solvePostorderTraversal(input string) string {
	root := parseTreeBFS(strings.TrimSpace(input))
	result := []string{}
	var postorder func(n *treeNode)
	postorder = func(n *treeNode) {
		if n == nil {
			return
		}
		postorder(n.left)
		postorder(n.right)
		result = append(result, strconv.Itoa(n.val))
	}
	postorder(root)
	return strings.Join(result, " ")
}
func solveInsertIntoBST(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	root := parseTreeBFS(lines[0])
	val, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	result := []int{}
	var inorder func(n *treeNode)
	inorder = func(n *treeNode) {
		if n == nil {
			return
		}
		inorder(n.left)
		result = append(result, n.val)
		inorder(n.right)
	}
	inorder(root)
	result = append(result, val)
	sort.Ints(result)
	return joinNums(result)
}
func solveDeleteNode(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	root := parseTreeBFS(lines[0])
	key, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	result := []int{}
	var inorder func(n *treeNode)
	inorder = func(n *treeNode) {
		if n == nil {
			return
		}
		inorder(n.left)
		if n.val != key {
			result = append(result, n.val)
		}
		inorder(n.right)
	}
	inorder(root)
	return joinNums(result)
}
func solveNumTrees(input string) string            {
	n, _ := strconv.Atoi(strings.TrimSpace(input))
	dp := make([]int, n+1)
	dp[0], dp[1] = 1, 1
	for i := 2; i <= n; i++ {
		for j := 0; j < i; j++ {
			dp[i] += dp[j] * dp[i-1-j]
		}
	}
	return fmt.Sprintf("%d", dp[n])
}
func solveOrangesRotting(input string) string      {
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return "0"
	}
	dims := strings.Fields(lines[0])
	m, _ := strconv.Atoi(dims[0])
	n, _ := strconv.Atoi(dims[1])
	grid := make([][]int, m)
	for i := 0; i < m; i++ {
		grid[i] = parseNums(lines[i+1])
	}
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
	if fresh == 0 {
		return "0"
	}
	dirs := [][2]int{{0,1},{0,-1},{1,0},{-1,0}}
	minutes := 0
	for len(queue) > 0 && fresh > 0 {
		minutes++
		size := len(queue)
		for i := 0; i < size; i++ {
			p := queue[0]
			queue = queue[1:]
			for _, d := range dirs {
				nr, nc := p.r+d[0], p.c+d[1]
				if nr >= 0 && nr < m && nc >= 0 && nc < n && grid[nr][nc] == 1 {
					grid[nr][nc] = 2
					fresh--
					queue = append(queue, pos{nr, nc})
				}
			}
		}
	}
	if fresh > 0 {
		return "-1"
	}
	return fmt.Sprintf("%d", minutes)
}
func solveCanFinish(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 1 {
		return "true"
	}
	numCourses, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	graph := make([][]int, numCourses)
	inDegree := make([]int, numCourses)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		graph[b] = append(graph[b], a)
		inDegree[a]++
	}
	queue := []int{}
	for i := 0; i < numCourses; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		count++
		for _, next := range graph[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if count == numCourses {
		return "true"
	}
	return "false"
}
func solveTrie(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	trie := map[string]bool{}
	results := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		op, word := parts[0], parts[1]
		switch op {
		case "insert":
			trie[word] = true
		case "search":
			if trie[word] {
				results = append(results, "true")
			} else {
				results = append(results, "false")
			}
		case "startsWith":
			found := false
			for k := range trie {
				if strings.HasPrefix(k, word) {
					found = true
					break
				}
			}
			if found {
				results = append(results, "true")
			} else {
				results = append(results, "false")
			}
		}
	}
	return strings.Join(results, " ")
}
func solveLetterCombinations(input string) string  {
	digits := strings.TrimSpace(input)
	if len(digits) == 0 {
		return ""
	}
	mapping := []string{"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"}
	result := []string{""}
	for _, d := range digits {
		num := int(d - '0')
		letters := mapping[num]
		newResult := []string{}
		for _, prefix := range result {
			for _, l := range letters {
				newResult = append(newResult, prefix+string(l))
			}
		}
		result = newResult
	}
	return strings.Join(result, " ")
}
func solveCombinationSum(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	candidates := parseNums(lines[0])
	target, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	sort.Ints(candidates)
	result := [][]int{}
	var backtrack func(start, remain int, path []int)
	backtrack = func(start, remain int, path []int) {
		if remain == 0 {
			cp := make([]int, len(path))
			copy(cp, path)
			result = append(result, cp)
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remain {
				break
			}
			backtrack(i, remain-candidates[i], append(path, candidates[i]))
		}
	}
	backtrack(0, target, []int{})
	rows := []string{}
	for _, r := range result {
		rows = append(rows, joinNums(r))
	}
	return strings.Join(rows, "\n")
}
func solveGenerateParenthesis(input string) string {
	n, _ := strconv.Atoi(strings.TrimSpace(input))
	result := []string{}
	var gen func(cur string, open, close int)
	gen = func(cur string, open, close int) {
		if len(cur) == 2*n {
			result = append(result, cur)
			return
		}
		if open < n {
			gen(cur+"(", open+1, close)
		}
		if close < open {
			gen(cur+")", open, close+1)
		}
	}
	gen("", 0, 0)
	return strings.Join(result, " ")
}
func solveWordExist(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return "false"
	}
	// First line may be "m n" header
	startLine := 0
	firstFields := strings.Fields(lines[0])
	if len(firstFields) == 2 {
		_, e1 := strconv.Atoi(firstFields[0])
		_, e2 := strconv.Atoi(firstFields[1])
		if e1 == nil && e2 == nil {
			startLine = 1
		}
	}
	grid := [][]byte{}
	for i := startLine; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		row := make([]byte, len(fields))
		for j, f := range fields {
			row[j] = f[0]
		}
		grid = append(grid, row)
	}
	word := strings.TrimSpace(lines[len(lines)-1])
	if len(grid) == 0 || word == "" {
		return "false"
	}
	m := len(grid)
	n := len(grid[0])
	var dfs func(i, j, k int) bool
	dfs = func(i, j, k int) bool {
		if k == len(word) {
			return true
		}
		if i < 0 || i >= m || j < 0 || j >= n || grid[i][j] != word[k] {
			return false
		}
		tmp := grid[i][j]
		grid[i][j] = '#'
		found := dfs(i+1, j, k+1) || dfs(i-1, j, k+1) || dfs(i, j+1, k+1) || dfs(i, j-1, k+1)
		grid[i][j] = tmp
		return found
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if dfs(i, j, 0) {
				return "true"
			}
		}
	}
	return "false"
}
func solvePermuteUnique(input string) string {
	nums := parseNums(input)
	sort.Ints(nums)
	result := [][]int{}
	used := make([]bool, len(nums))
	var backtrack func(path []int)
	backtrack = func(path []int) {
		if len(path) == len(nums) {
			cp := make([]int, len(path))
			copy(cp, path)
			result = append(result, cp)
			return
		}
		for i := 0; i < len(nums); i++ {
			if used[i] {
				continue
			}
			if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			backtrack(append(path, nums[i]))
			used[i] = false
		}
	}
	backtrack([]int{})
	rows := []string{}
	for _, r := range result {
		rows = append(rows, joinNums(r))
	}
	return strings.Join(rows, "\n")
}
func solveNextPermutation(input string) string     {
	nums := parseNums(input)
	n := len(nums)
	if n <= 1 {
		return joinNums(nums)
	}
	// 找到第一个递减的位置
	i := n - 2
	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	}
	if i >= 0 {
		j := n - 1
		for j > i && nums[j] <= nums[i] {
			j--
		}
		nums[i], nums[j] = nums[j], nums[i]
	}
	// 反转 i+1 到末尾
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		nums[l], nums[r] = nums[r], nums[l]
	}
	return joinNums(nums)
}
func solveSearchRotated(input string) string       {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	target, _ := strconv.Atoi(lines[1])
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == target {
			return fmt.Sprintf("%d", mid)
		}
		if nums[left] <= nums[mid] {
			if target >= nums[left] && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			if target > nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	return "-1"
}
func solveSearchRange(input string) string         {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	target, _ := strconv.Atoi(lines[1])
	findFirst := func() int {
		left, right := 0, len(nums)-1
		result := -1
		for left <= right {
			mid := (left + right) / 2
			if nums[mid] == target {
				result = mid
				right = mid - 1
			} else if nums[mid] < target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return result
	}
	findLast := func() int {
		left, right := 0, len(nums)-1
		result := -1
		for left <= right {
			mid := (left + right) / 2
			if nums[mid] == target {
				result = mid
				left = mid + 1
			} else if nums[mid] < target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return result
	}
	return fmt.Sprintf("%d %d", findFirst(), findLast())
}
func solveMinStack(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	stack := []int{}
	minStack := []int{}
	results := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "push":
			val, _ := strconv.Atoi(parts[1])
			stack = append(stack, val)
			if len(minStack) == 0 || val <= minStack[len(minStack)-1] {
				minStack = append(minStack, val)
			} else {
				minStack = append(minStack, minStack[len(minStack)-1])
			}
		case "pop":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
				minStack = minStack[:len(minStack)-1]
			}
		case "top":
			if len(stack) > 0 {
				results = append(results, strconv.Itoa(stack[len(stack)-1]))
			}
		case "getMin":
			if len(minStack) > 0 {
				results = append(results, strconv.Itoa(minStack[len(minStack)-1]))
			}
		}
	}
	return strings.Join(results, " ")
}
func solveDecodeString(input string) string {
	input = strings.TrimSpace(input)
	stack := []string{}
	countStack := []int{}
	current := ""
	k := 0
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			k = k*10 + int(ch-'0')
		} else if ch == '[' {
			countStack = append(countStack, k)
			stack = append(stack, current)
			current = ""
			k = 0
		} else if ch == ']' {
			cnt := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]
			prev := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			repeated := ""
			for i := 0; i < cnt; i++ {
				repeated += current
			}
			current = prev + repeated
		} else {
			current += string(ch)
		}
	}
	return current
}
func solveDailyTemperatures(input string) string   {
	temps := parseNums(input)
	n := len(temps)
	result := make([]int, n)
	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && temps[i] > temps[stack[len(stack)-1]] {
			idx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[idx] = i - idx
		}
		stack = append(stack, i)
	}
	return joinNums(result)
}
func solveLargestRectangle(input string) string {
	heights := parseNums(input)
	n := len(heights)
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
	return fmt.Sprintf("%d", maxArea)
}
func solveLongestValidParentheses(input string) string {
	s := strings.TrimSpace(input)
	stack := []int{-1}
	maxLen := 0
	for i, ch := range s {
		if ch == '(' {
			stack = append(stack, i)
		} else {
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				stack = append(stack, i)
			} else {
				l := i - stack[len(stack)-1]
				if l > maxLen {
					maxLen = l
				}
			}
		}
	}
	return fmt.Sprintf("%d", maxLen)
}
func solveMyQueue(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	in := []int{}
	out := []int{}
	results := []string{}
	move := func() {
		if len(out) == 0 {
			for len(in) > 0 {
				out = append(out, in[len(in)-1])
				in = in[:len(in)-1]
			}
		}
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "push":
			val, _ := strconv.Atoi(parts[1])
			in = append(in, val)
		case "pop":
			move()
			if len(out) > 0 {
				results = append(results, strconv.Itoa(out[len(out)-1]))
				out = out[:len(out)-1]
			}
		case "peek":
			move()
			if len(out) > 0 {
				results = append(results, strconv.Itoa(out[len(out)-1]))
			}
		case "empty":
			if len(in)+len(out) == 0 {
				results = append(results, "true")
			} else {
				results = append(results, "false")
			}
		}
	}
	return strings.Join(results, " ")
}
func solveFindKthLargest(input string) string      {
	lines := strings.Split(input, "\n")
	nums := parseNums(lines[0])
	k, _ := strconv.Atoi(lines[1])
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	return fmt.Sprintf("%d", nums[k-1])
}
func solveTopKFrequent(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return ""
	}
	nums := parseNums(lines[0])
	k, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	freq := make(map[int]int)
	for _, n := range nums {
		freq[n]++
	}
	type pair struct{ val, cnt int }
	pairs := []pair{}
	for v, c := range freq {
		pairs = append(pairs, pair{v, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].cnt != pairs[j].cnt {
			return pairs[i].cnt > pairs[j].cnt
		}
		return pairs[i].val < pairs[j].val
	})
	result := []int{}
	for i := 0; i < k && i < len(pairs); i++ {
		result = append(result, pairs[i].val)
	}
	sort.Ints(result)
	return joinNums(result)
}
func solveMedianFinder(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	nums := []int{}
	results := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if parts[0] == "addNum" {
			val, _ := strconv.Atoi(parts[1])
			// Insert in sorted position
			pos := sort.SearchInts(nums, val)
			nums = append(nums, 0)
			copy(nums[pos+1:], nums[pos:])
			nums[pos] = val
		} else if parts[0] == "findMedian" {
			n := len(nums)
			if n == 0 {
				results = append(results, "0.00")
			} else if n%2 == 1 {
				results = append(results, fmt.Sprintf("%.2f", float64(nums[n/2])))
			} else {
				median := float64(nums[n/2-1]+nums[n/2]) / 2.0
				results = append(results, fmt.Sprintf("%.2f", median))
			}
		}
	}
	return strings.Join(results, " ")
}
func solveSortColors(input string) string          {
	nums := parseNums(input)
	low, mid, high := 0, 0, len(nums)-1
	for mid <= high {
		if nums[mid] == 0 {
			nums[low], nums[mid] = nums[mid], nums[low]
			low++
			mid++
		} else if nums[mid] == 1 {
			mid++
		} else {
			nums[mid], nums[high] = nums[high], nums[mid]
			high--
		}
	}
	return joinNums(nums)
}
func solveRob(input string) string                 {
	nums := parseNums(input)
	n := len(nums)
	if n == 0 {
		return "0"
	}
	if n == 1 {
		return fmt.Sprintf("%d", nums[0])
	}
	prev2, prev1 := 0, nums[0]
	for i := 1; i < n; i++ {
		cur := prev1
		if prev2+nums[i] > cur {
			cur = prev2 + nums[i]
		}
		prev2, prev1 = prev1, cur
	}
	return fmt.Sprintf("%d", prev1)
}
func solveNumSquares(input string) string          {
	n, _ := strconv.Atoi(strings.TrimSpace(input))
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = i
		for j := 1; j*j <= i; j++ {
			if dp[i-j*j]+1 < dp[i] {
				dp[i] = dp[i-j*j] + 1
			}
		}
	}
	return fmt.Sprintf("%d", dp[n])
}
func solveCoinChange(input string) string          {
	lines := strings.Split(input, "\n")
	coins := parseNums(lines[0])
	amount, _ := strconv.Atoi(lines[1])
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = amount + 1
	}
	dp[0] = 0
	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	if dp[amount] > amount {
		return "-1"
	}
	return fmt.Sprintf("%d", dp[amount])
}
func solveWordBreak(input string) string {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return "false"
	}
	s := strings.TrimSpace(lines[0])
	wordList := strings.Fields(lines[1])
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}
	dp := make([]bool, len(s)+1)
	dp[0] = true
	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			if dp[j] && wordSet[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	if dp[len(s)] {
		return "true"
	}
	return "false"
}
func solveLengthOfLIS(input string) string         {
	nums := parseNums(input)
	if len(nums) == 0 {
		return "0"
	}
	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
	}
	maxLen := 1
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
		}
	}
	return fmt.Sprintf("%d", maxLen)
}
func solveMaxProduct(input string) string          {
	nums := parseNums(input)
	if len(nums) == 0 {
		return "0"
	}
	maxProd, minProd, result := nums[0], nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < 0 {
			maxProd, minProd = minProd, maxProd
		}
		if nums[i] > maxProd*nums[i] {
			maxProd = nums[i]
		} else {
			maxProd = maxProd * nums[i]
		}
		if nums[i] < minProd*nums[i] {
			minProd = nums[i]
		} else {
			minProd = minProd * nums[i]
		}
		if maxProd > result {
			result = maxProd
		}
	}
	return fmt.Sprintf("%d", result)
}
func solveMinDistance(input string) string {
	input = strings.TrimRight(input, "\n\r")
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return "0"
	}
	word1 := strings.TrimSpace(lines[0])
	word2 := strings.TrimSpace(lines[1])
	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return fmt.Sprintf("%d", dp[m][n])
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
func solveFindDuplicate(input string) string {
	nums := parseNums(input)
	slow, fast := nums[0], nums[nums[0]]
	for slow != fast {
		slow = nums[slow]
		fast = nums[nums[fast]]
	}
	slow = 0
	for slow != fast {
		slow = nums[slow]
		fast = nums[fast]
	}
	return fmt.Sprintf("%d", slow)
}
func solveIntToRoman(input string) string {
	num, _ := strconv.Atoi(strings.TrimSpace(input))
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	result := ""
	for i, v := range vals {
		for num >= v {
			result += syms[i]
			num -= v
		}
	}
	return result
}
func solveRomanToInt(input string) string {
	s := strings.TrimSpace(input)
	romanMap := map[byte]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}
	result := 0
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && romanMap[s[i]] < romanMap[s[i+1]] {
			result -= romanMap[s[i]]
		} else {
			result += romanMap[s[i]]
		}
	}
	return fmt.Sprintf("%d", result)
}

// ===== 主程序 =====

func main() {
	if err := config.Load("./config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := database.Init(&config.AppConfig.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	var problems []models.Problem
	if err := database.DB.Order("id asc").Find(&problems).Error; err != nil {
		log.Fatalf("Failed to query problems: %v", err)
	}

	fmt.Println("========== 测试数据验证报告 ==========")
	fmt.Printf("总题目数: %d\n\n", len(problems))

	totalCases := 0
	passedCases := 0
	failedProblems := []string{}

	for _, p := range problems {
		solve, ok := solutions[p.Title]
		if !ok {
			fmt.Printf("[跳过] %s - 无解法实现\n", p.Title)
			continue
		}

		var testCases []models.TestCase
		if err := database.DB.Where("problem_id = ?", p.ID).Order("id asc").Find(&testCases).Error; err != nil {
			log.Printf("Failed to query test cases for %s: %v", p.Title, err)
			continue
		}

		problemPassed := 0
		problemFailed := 0

		for _, tc := range testCases {
			totalCases++
			expected := solve(tc.Input)
			// 标准化比较：去除首尾空白
			expected = strings.TrimSpace(expected)
			actual := strings.TrimSpace(tc.Output)

			if expected == actual {
				passedCases++
				problemPassed++
			} else {
				problemFailed++
				if problemFailed <= 2 { // 只打印前2个失败用例
					fmt.Printf("[失败] %s - 用例 #%d\n", p.Title, tc.ID)
					fmt.Printf("  输入: %s\n", truncate(tc.Input, 100))
					fmt.Printf("  期望(DB): %s\n", truncate(actual, 100))
					fmt.Printf("  实际(代码): %s\n", truncate(expected, 100))
				}
			}
		}

		if problemFailed > 0 {
			failedProblems = append(failedProblems, p.Title)
			fmt.Printf("[统计] %s: 通过 %d/%d\n\n", p.Title, problemPassed, len(testCases))
		}
	}

	fmt.Println("\n========== 汇总 ==========")
	fmt.Printf("总测试用例: %d\n", totalCases)
	fmt.Printf("通过: %d (%.1f%%)\n", passedCases, float64(passedCases)/float64(totalCases)*100)
	fmt.Printf("失败: %d\n", totalCases-passedCases)
	fmt.Printf("有问题题目: %d\n", len(failedProblems))
	if len(failedProblems) > 0 {
		fmt.Println("\n需要修正的题目:")
		for _, t := range failedProblems {
			fmt.Printf("  - %s\n", t)
		}
	}
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
