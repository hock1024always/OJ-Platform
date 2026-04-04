package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// 题目解题思路和代码映射
var problemSolutions = map[string]struct {
	Approach  string
	Code      string
	Tags      string
	Difficulty string
}{
	"两数之和": {
		Approach: "使用哈希表存储已遍历的数字及其索引，对于每个数字，检查 target - num 是否在哈希表中",
		Code: `func twoSum(nums []int, target int) []int {
    m := make(map[int]int)
    for i, num := range nums {
        if j, ok := m[target-num]; ok {
            return []int{j, i}
        }
        m[num] = i
    }
    return nil
}`,
		Tags: "数组,哈希表",
		Difficulty: "Easy",
	},
	"二叉树的中序遍历": {
		Approach: "递归：左子树 -> 根节点 -> 右子树。或使用栈迭代实现",
		Code: `func inorderTraversal(root *TreeNode) []int {
    var result []int
    var inorder func(*TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil {
            return
        }
        inorder(node.Left)
        result = append(result, node.Val)
        inorder(node.Right)
    }
    inorder(root)
    return result
}`,
		Tags: "二叉树,递归",
		Difficulty: "Easy",
	},
	"二叉树的最大深度": {
		Approach: "递归计算左右子树的最大深度，取较大值加1",
		Code: `func maxDepth(root *TreeNode) int {
    if root == nil {
        return 0
    }
    left := maxDepth(root.Left)
    right := maxDepth(root.Right)
    if left > right {
        return left + 1
    }
    return right + 1
}`,
		Tags: "二叉树,递归,深度优先搜索",
		Difficulty: "Easy",
	},
	"翻转二叉树": {
		Approach: "递归交换每个节点的左右子树",
		Code: `func invertTree(root *TreeNode) *TreeNode {
    if root == nil {
        return nil
    }
    root.Left, root.Right = invertTree(root.Right), invertTree(root.Left)
    return root
}`,
		Tags: "二叉树,递归",
		Difficulty: "Easy",
	},
	"最大子数组和": {
		Approach: "Kadane算法：维护当前子数组和，如果当前和为负数则重新开始",
		Code: `func maxSubArray(nums []int) int {
    maxSum := nums[0]
    currSum := nums[0]
    for i := 1; i < len(nums); i++ {
        if currSum < 0 {
            currSum = nums[i]
        } else {
            currSum += nums[i]
        }
        if currSum > maxSum {
            maxSum = currSum
        }
    }
    return maxSum
}`,
		Tags: "数组,动态规划,分治",
		Difficulty: "Medium",
	},
}

func main() {
	db, err := sql.Open("sqlite3", "oj_platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS problem_solutions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			problem_id INTEGER NOT NULL UNIQUE,
			title TEXT NOT NULL,
			difficulty TEXT,
			tags TEXT,
			solution_approach TEXT,
			solution_code TEXT,
			time_complexity TEXT,
			space_complexity TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("创建表失败:", err)
	}

	// 获取所有题目
	rows, err := db.Query("SELECT id, title FROM problems")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	stmt, err := db.Prepare(`
		INSERT OR REPLACE INTO problem_solutions 
		(problem_id, title, difficulty, tags, solution_approach, solution_code, time_complexity, space_complexity)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var id int
		var title string
		rows.Scan(&id, &title)

		if sol, ok := problemSolutions[title]; ok {
			_, err := stmt.Exec(id, title, sol.Difficulty, sol.Tags, sol.Approach, sol.Code, "O(n)", "O(n)")
			if err != nil {
				log.Printf("插入 %s 失败: %v", title, err)
			} else {
				count++
				fmt.Printf("✓ 已添加: %s\n", title)
			}
		}
	}

	fmt.Printf("\n共添加 %d 道题目的解题数据\n", count)
}
