-- AI解题助手数据库表结构

-- 题目详细内容表（包含解题思路和源代码）
CREATE TABLE IF NOT EXISTS problem_solutions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,                          -- 题目标题
    description TEXT,                             -- 题目描述
    difficulty TEXT,                              -- 难度：easy/medium/hard
    tags TEXT,                                    -- 标签，JSON数组格式
    solution_approach TEXT,                       -- 解题思路/方法
    solution_code TEXT,                           -- Go语言标准解答代码
    time_complexity TEXT,                         -- 时间复杂度分析
    space_complexity TEXT,                        -- 空间复杂度分析
    key_points TEXT,                              -- 关键点/易错点
    similar_problems TEXT,                        -- 相似题目ID列表
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id)
);

-- 倒排索引表（用于全文搜索）
CREATE TABLE IF NOT EXISTS inverted_index (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word TEXT NOT NULL,                           -- 关键词
    problem_id INTEGER NOT NULL,                  -- 题目ID
    field_type TEXT NOT NULL,                     -- 字段类型：title/description/solution/code
    position INTEGER,                             -- 词在文档中的位置
    weight REAL DEFAULT 1.0,                      -- 权重
    UNIQUE(word, problem_id, field_type, position)
);

-- 创建全文搜索索引
CREATE INDEX IF NOT EXISTS idx_inverted_word ON inverted_index(word);
CREATE INDEX IF NOT EXISTS idx_inverted_problem ON inverted_index(problem_id);
CREATE INDEX IF NOT EXISTS idx_solution_problem ON problem_solutions(problem_id);

-- 搜索历史表（用于优化和统计）
CREATE TABLE IF NOT EXISTS search_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,                          -- 搜索查询
    results_count INTEGER,                        -- 返回结果数
    clicked_problem_id INTEGER,                   -- 用户点击的题目
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- MCP调用日志表
CREATE TABLE IF NOT EXISTS mcp_call_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT,                              -- 会话ID
    request_type TEXT,                            -- 请求类型：search/get_problem/compare
    request_params TEXT,                          -- 请求参数JSON
    response_data TEXT,                           -- 响应数据JSON
    latency_ms INTEGER,                           -- 响应延迟
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
