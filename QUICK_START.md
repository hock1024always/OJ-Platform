# OJ Platform 快速使用指南

## 🎉 部署完成！

OJ刷题平台已经成功部署并导入力扣Hot100经典题目！

---

## 📍 访问地址

**前端页面**: http://172.20.19.106:8080

**API地址**: http://172.20.19.106:8080/api/v1/

**健康检查**: http://172.20.19.106:8080/health

---

## 📚 题库内容

已导入 **15道** 力扣Hot100经典题目：

### Easy (简单)
1. 两数之和
2. 反转链表
3. 合并两个有序链表
4. 爬楼梯
5. 二叉树的最大深度
6. 买卖股票的最佳时机
7. 只出现一次的数字
8. 多数元素
9. 移动零
10. 回文链表
11. 环形链表
12. 合并两个有序数组
13. 验证回文串

### Medium (中等)
1. 最大子数组和
2. 找到字符串中所有字母异位词

---

## 🚀 使用流程

### 1. 注册账号
访问: http://172.20.19.106:8080

- 点击"注册"标签
- 输入用户名、邮箱和密码
- 点击"注册"按钮

### 2. 登录系统
- 使用刚才注册的账号登录
- 系统会返回JWT Token并自动保存

### 3. 浏览题目
- 登录后自动跳转到题目列表页
- 点击任意题目查看详情

### 4. 提交代码
- 在题目详情页编写Go代码
- 点击"提交代码"按钮
- 等待判题结果（通常1-2秒）

### 5. 查看结果
系统会返回以下状态之一：
- ✅ **Accepted** - 代码通过所有测试用例
- ❌ **Wrong Answer** - 答案错误
- ⚠️ **Compile Error** - 编译错误
- ⏱️ **Time Limit Exceeded** - 运行超时
- 💥 **Runtime Error** - 运行时错误

---

## 💻 示例代码

### 两数之和（题目ID: 1）
```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// 读取数组
	scanner.Scan()
	numsStr := strings.Fields(scanner.Text())
	nums := make([]int, len(numsStr))
	for i, s := range numsStr {
		nums[i], _ = strconv.Atoi(s)
	}

	// 读取目标值
	scanner.Scan()
	target, _ := strconv.Atoi(scanner.Text())

	// 两数之和算法
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				fmt.Printf("%d %d", i, j)
				return
			}
		}
	}
}
```

---

## 🔧 管理命令

### 查看服务状态
```bash
ps aux | grep "bin/server"
```

### 查看日志
```bash
cd /home/haoqian.li/compile_dockers/oj-platform
tail -f server.log
```

### 重启服务
```bash
cd /home/haoqian.li/compile_dockers/oj-platform
lsof -ti:8080 | xargs kill -9
nohup ./bin/server > server.log 2>&1 &
```

### 停止服务
```bash
lsof -ti:8080 | xargs kill -9
```

---

## 📊 系统信息

- **服务器IP**: 172.20.19.106
- **端口**: 8080
- **数据库**: SQLite (oj_platform.db)
- **并发Worker**: 20个
- **超时时间**: 5秒
- **内存限制**: 256MB

---

## ⚠️ 注意事项

1. **仅支持Go语言** - 当前版本只支持Go代码提交
2. **输入格式** - 按照题目要求的格式输入
3. **输出格式** - 必须与期望输出完全一致（包括空格和换行）
4. **并发限制** - 系统支持最多20个并发判题任务
5. **超时机制** - 代码执行超过5秒会自动终止

---

## 🐛 故障排查

### 无法访问
1. 检查服务是否运行: `ps aux | grep bin/server`
2. 检查端口是否监听: `netstat -tlnp | grep 8080`
3. 检查防火墙设置

### 判题失败
1. 查看日志: `tail -f server.log`
2. 检查代码语法
3. 确认输入输出格式

### 登录失败
1. 确认用户名密码正确
2. 检查浏览器控制台错误信息
3. 清除浏览器缓存重试

---

## 📞 技术支持

- **项目路径**: `/home/haoqian.li/compile_dockers/oj-platform`
- **配置文件**: `config.yaml`
- **数据库**: `oj_platform.db`
- **日志文件**: `server.log`

---

## 🎯 下一步

1. 尝试提交第一道题目（两数之和）
2. 查看判题结果
3. 继续挑战其他题目

祝刷题愉快！🎉
