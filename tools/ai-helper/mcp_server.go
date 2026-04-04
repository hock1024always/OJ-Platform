package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// MCPRequest MCP请求结构
type MCPRequest struct {
	SessionID string                 `json:"session_id"`
	Type      string                 `json:"type"`
	Params    map[string]interface{} `json:"params"`
}

// MCPResponse MCP响应结构
type MCPResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ProblemResult 搜索结果
type ProblemResult struct {
	ProblemID         int     `json:"problem_id"`
	Title             string  `json:"title"`
	Difficulty        string  `json:"difficulty"`
	Tags              string  `json:"tags"`
	SolutionApproach  string  `json:"solution_approach"`
	SolutionCode      string  `json:"solution_code"`
	Score             float64 `json:"score"`
}

type MCPServer struct {
	db *sql.DB
}

// 全文搜索
func (s *MCPServer) searchProblems(query string, limit int) ([]ProblemResult, error) {
	// 分词
	words := tokenize(query)
	if len(words) == 0 {
		return nil, nil
	}

	// 构建查询
	placeholders := make([]string, len(words))
	args := make([]interface{}, len(words))
	for i, w := range words {
		placeholders[i] = "?"
		args[i] = w
	}

	sqlQuery := fmt.Sprintf(`
		SELECT 
			ps.problem_id,
			ps.title,
			ps.difficulty,
			ps.tags,
			ps.solution_approach,
			ps.solution_code,
			SUM(ii.weight) as score
		FROM problem_solutions ps
		JOIN inverted_index ii ON ps.problem_id = ii.problem_id
		WHERE ii.word IN (%s)
		GROUP BY ps.problem_id
		ORDER BY score DESC
		LIMIT ?
	`, strings.Join(placeholders, ","))

	args = append(args, limit)

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ProblemResult
	for rows.Next() {
		var r ProblemResult
		err := rows.Scan(&r.ProblemID, &r.Title, &r.Difficulty, &r.Tags, 
			&r.SolutionApproach, &r.SolutionCode, &r.Score)
		if err != nil {
			continue
		}
		results = append(results, r)
	}

	return results, nil
}

// 获取单个题目详情
func (s *MCPServer) getProblem(problemID int) (*ProblemResult, error) {
	var r ProblemResult
	err := s.db.QueryRow(`
		SELECT problem_id, title, difficulty, tags, solution_approach, solution_code
		FROM problem_solutions WHERE problem_id = ?
	`, problemID).Scan(&r.ProblemID, &r.Title, &r.Difficulty, &r.Tags, 
		&r.SolutionApproach, &r.SolutionCode)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// MCP处理函数
func (s *MCPServer) handleMCP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond(w, MCPResponse{Success: false, Error: err.Error()})
		return
	}

	var resp MCPResponse
	switch req.Type {
	case "search":
		query, _ := req.Params["query"].(string)
		limit := 5
		if l, ok := req.Params["limit"].(float64); ok {
			limit = int(l)
		}
		results, err := s.searchProblems(query, limit)
		if err != nil {
			resp = MCPResponse{Success: false, Error: err.Error()}
		} else {
			resp = MCPResponse{Success: true, Data: results}
		}

	case "get_problem":
		pid := int(req.Params["problem_id"].(float64))
		problem, err := s.getProblem(pid)
		if err != nil {
			resp = MCPResponse{Success: false, Error: err.Error()}
		} else {
			resp = MCPResponse{Success: true, Data: problem}
		}

	case "get_solution_code":
		pid := int(req.Params["problem_id"].(float64))
		var code string
		err := s.db.QueryRow("SELECT solution_code FROM problem_solutions WHERE problem_id = ?", pid).Scan(&code)
		if err != nil {
			resp = MCPResponse{Success: false, Error: err.Error()}
		} else {
			resp = MCPResponse{Success: true, Data: map[string]string{"code": code}}
		}

	default:
		resp = MCPResponse{Success: false, Error: "unknown request type"}
	}

	// 记录日志
	latency := time.Since(start).Milliseconds()
	logMCPCall(s.db, req.SessionID, req.Type, req.Params, resp, int(latency))

	respond(w, resp)
}

func respond(w http.ResponseWriter, resp MCPResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func logMCPCall(db *sql.DB, sessionID, reqType string, params, response interface{}, latency int) {
	paramsJSON, _ := json.Marshal(params)
	responseJSON, _ := json.Marshal(response)
	
	db.Exec(`
		INSERT INTO mcp_call_logs (session_id, request_type, request_params, response_data, latency_ms)
		VALUES (?, ?, ?, ?, ?)
	`, sessionID, reqType, string(paramsJSON), string(responseJSON), latency)
}

func main() {
	db, err := sql.Open("sqlite3", "oj_platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := &MCPServer{db: db}

	http.HandleFunc("/mcp", server.handleMCP)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	fmt.Println("MCP Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
