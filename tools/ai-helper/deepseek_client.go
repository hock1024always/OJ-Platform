package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const DeepSeekAPI = "https://api.deepseek.com/v1/chat/completions"

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message DeepSeekMessage `json:"message"`
	} `json:"choices"`
}

type MCPClient struct {
	baseURL string
}

func NewMCPClient(baseURL string) *MCPClient {
	return &MCPClient{baseURL: baseURL}
}

func (c *MCPClient) SearchProblems(query string) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"type":   "search",
		"params": map[string]interface{}{"query": query, "limit": 3},
	})
	
	resp, err := http.Post(c.baseURL+"/mcp", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func (c *MCPClient) GetSolutionCode(problemID int) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"type":   "get_solution_code",
		"params": map[string]interface{}{"problem_id": problemID},
	})
	
	resp, err := http.Post(c.baseURL+"/mcp", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func callDeepSeek(apiKey string, messages []DeepSeekMessage) (string, error) {
	req := DeepSeekRequest{
		Model:    "deepseek-chat",
		Messages: messages,
		Stream:   false,
	}
	
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", DeepSeekAPI, bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result DeepSeekResponse
	json.Unmarshal(body, &result)
	
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from DeepSeek")
}

// AI解题助手主函数
func SolveWithAI(problemDescription string) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Println("请设置 DEEPSEEK_API_KEY 环境变量")
		return
	}
	
	// 1. 使用MCP搜索相关题目
	mcpClient := NewMCPClient("http://localhost:8080")
	searchResult, _ := mcpClient.SearchProblems(problemDescription)
	
	// 2. 构建提示词
	prompt := fmt.Sprintf(`你是一个算法解题助手。用户遇到了以下问题：

%s

我为你搜索到了一些相关的题目和解答：

%s

请：
1. 分析这个问题属于哪种算法类型
2. 提供解题思路
3. 给出Go语言实现代码
4. 分析时间复杂度和空间复杂度`, problemDescription, searchResult)
	
	messages := []DeepSeekMessage{
		{Role: "system", Content: "你是专业的算法工程师，擅长解决LeetCode算法题。"},
		{Role: "user", Content: prompt},
	}
	
	// 3. 调用DeepSeek
	answer, err := callDeepSeek(apiKey, messages)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("=== AI解题助手回答 ===")
	fmt.Println(answer)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run deepseek_client.go \"<problem description>\"")
		return
	}
	
	problem := os.Args[1]
	SolveWithAI(problem)
}
