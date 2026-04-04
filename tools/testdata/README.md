# OJ平台测试数据集

## 数据说明

本压缩包包含修复后的5000个测试用例（100道题目 × 50组测试数据），通过率100%。

## 数据来源与检测方法

### 数据来源
1. **原始数据**：从LeetCode Hot 100题目生成
2. **修复过程**：使用自动化验证脚本检测并修复了233个错误测试用例
3. **验证方法**：
   ```bash
   # 运行验证脚本
   go run scripts/validate_testcases.go
   ```

### 数据检测逻辑
验证脚本通过以下方式检测测试用例：
1. 读取数据库中的输入数据
2. 使用标准算法计算正确输出
3. 对比数据库存储的期望输出
4. 不匹配则标记为失败

## 压缩包内容

```
test_data/
├── test_cases.sql          # 测试用例SQL导出文件
├── test_cases_summary.txt  # 测试用例统计摘要
├── validation_report.md    # 详细验证报告
└── README.md               # 本说明文件
```

## 使用方法

### 1. 解压缩

```bash
# Linux/Mac
tar -xzvf oj_test_data_v1.0.tar.gz

# Windows (使用7-Zip或PowerShell)
# 右键解压或使用: Expand-Archive -Path oj_test_data_v1.0.zip -DestinationPath ./test_data
```

### 2. 导入数据库

```bash
# 进入项目目录
cd oj-platform

# 备份原数据库（可选）
cp oj_platform.db oj_platform.db.backup

# 导入测试数据
sqlite3 oj_platform.db < test_data/test_cases.sql
```

### 3. 验证导入结果

```bash
# 运行验证脚本
go run scripts/validate_testcases.go

# 预期输出：通过: 5000 (100.0%)
```

## 数据格式说明

### 测试用例表结构
```sql
CREATE TABLE test_cases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER NOT NULL,  -- 题目ID
    input TEXT NOT NULL,          -- 测试输入
    output TEXT NOT NULL,         -- 期望输出
    is_public BOOLEAN DEFAULT false,
    created_at DATETIME,
    deleted_at DATETIME
);
```

### 输入输出格式示例

**二叉树题目**：
- 输入: `3 9 20 null null 15 7` (BFS层序遍历)
- 输出: `3 9 20 null null 15 7` (根据题目要求)

**数组题目**：
- 输入: `2 7 11 15\n9` (数组+目标值，换行分隔)
- 输出: `0 1` (索引位置)

**字符串题目**：
- 输入: `abcabcbb`
- 输出: `3` (最长无重复子串长度)

## 版本信息

- **版本**: v1.0
- **生成日期**: 2026-03-19
- **测试用例总数**: 5000
- **题目数量**: 100
- **通过率**: 100%

## 注意事项

1. 导入前请确保数据库表结构已创建
2. 导入会覆盖现有test_cases表数据
3. 建议在导入前备份原数据库
