package codegen

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// GenerateRandomInput 根据函数签名和输入约束生成随机输入
func GenerateRandomInput(sig *FunctionSignature, constraints []InputConstraint, rng *rand.Rand) (string, error) {
	constraintMap := make(map[string]*InputConstraint)
	for i := range constraints {
		constraintMap[constraints[i].ParamName] = &constraints[i]
	}

	var lines []string
	for _, p := range sig.Params {
		c := constraintMap[p.Name]
		line, err := generateParamInput(p, c, rng)
		if err != nil {
			return "", fmt.Errorf("generate input for %s: %w", p.Name, err)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func generateParamInput(p Param, c *InputConstraint, rng *rand.Rand) (string, error) {
	minVal, maxVal := -1000, 1000
	minLen, maxLen := 1, 20
	minNodes, maxNodes := 0, 10

	if c != nil {
		if c.MinVal != 0 || c.MaxVal != 0 {
			minVal = c.MinVal
			maxVal = c.MaxVal
		}
		if c.MinLen != 0 || c.MaxLen != 0 {
			minLen = c.MinLen
			maxLen = c.MaxLen
		}
		if c.MinNodes != 0 || c.MaxNodes != 0 {
			minNodes = c.MinNodes
			maxNodes = c.MaxNodes
		}
	}

	switch p.Type {
	case TypeInt:
		return strconv.Itoa(randInt(rng, minVal, maxVal)), nil
	case TypeFloat:
		return fmt.Sprintf("%.2f", float64(randInt(rng, minVal*100, maxVal*100))/100.0), nil
	case TypeString:
		length := randInt(rng, minLen, maxLen)
		return randString(rng, length), nil
	case TypeBool:
		if rng.Intn(2) == 0 {
			return "true", nil
		}
		return "false", nil
	case TypeIntArray:
		length := randInt(rng, minLen, maxLen)
		nums := make([]string, length)
		for i := range nums {
			nums[i] = strconv.Itoa(randInt(rng, minVal, maxVal))
		}
		return strings.Join(nums, " "), nil
	case TypeStrArray:
		length := randInt(rng, minLen, maxLen)
		strs := make([]string, length)
		for i := range strs {
			strs[i] = randString(rng, randInt(rng, 1, 10))
		}
		return strings.Join(strs, " "), nil
	case TypeByteArray:
		length := randInt(rng, minLen, maxLen)
		return randString(rng, length), nil
	case TypeInt2D:
		rows := randInt(rng, minLen, maxLen)
		cols := randInt(rng, 1, 10)
		var lines []string
		lines = append(lines, strconv.Itoa(rows))
		for i := 0; i < rows; i++ {
			row := make([]string, cols)
			for j := range row {
				row[j] = strconv.Itoa(randInt(rng, minVal, maxVal))
			}
			lines = append(lines, strings.Join(row, " "))
		}
		return strings.Join(lines, "\n"), nil
	case TypeByte2D:
		rows := randInt(rng, minLen, maxLen)
		cols := randInt(rng, 1, 10)
		var lines []string
		lines = append(lines, strconv.Itoa(rows))
		for i := 0; i < rows; i++ {
			lines = append(lines, randString(rng, cols))
		}
		return strings.Join(lines, "\n"), nil
	case TypeListNode:
		n := randInt(rng, minNodes, maxNodes)
		if n == 0 {
			return "null", nil
		}
		vals := make([]string, n)
		for i := range vals {
			vals[i] = strconv.Itoa(randInt(rng, minVal, maxVal))
		}
		return strings.Join(vals, " "), nil
	case TypeTreeNode:
		n := randInt(rng, minNodes, maxNodes)
		if n == 0 {
			return "null", nil
		}
		return randTreeLevelOrder(rng, n, minVal, maxVal), nil
	default:
		return "", fmt.Errorf("unsupported type for input generation: %s", p.Type)
	}
}

func randInt(rng *rand.Rand, min, max int) int {
	if min >= max {
		return min
	}
	return min + rng.Intn(max-min+1)
}

func randString(rng *rand.Rand, length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b)
}

// randTreeLevelOrder 生成随机二叉树的层序遍历表示
func randTreeLevelOrder(rng *rand.Rand, n, minVal, maxVal int) string {
	if n == 0 {
		return "null"
	}
	vals := []string{strconv.Itoa(randInt(rng, minVal, maxVal))}
	count := 1
	slots := 2 // 可用的子节点槽位
	for count < n && slots > 0 {
		// 以较高概率生成节点
		if rng.Float64() < 0.7 && count < n {
			vals = append(vals, strconv.Itoa(randInt(rng, minVal, maxVal)))
			count++
			slots += 2
		} else {
			vals = append(vals, "null")
		}
		slots--
	}
	// 去掉尾部 null
	for len(vals) > 0 && vals[len(vals)-1] == "null" {
		vals = vals[:len(vals)-1]
	}
	return strings.Join(vals, " ")
}
