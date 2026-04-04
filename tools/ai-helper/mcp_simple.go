package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type MCPRequest struct {
	SessionID string                 `json:"session_id"`
	Type      string                 `json:"type"`
	Params    map[string]interface{} `json:"params"`
}

type MCPResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	db, err := sql.Open("sqlite3", "oj_platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		var req MCPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(MCPResponse{Success: false, Error: err.Error()})
			return
		}

		var resp MCPResponse
		switch req.Type {
		case "search":
			query, _ := req.Params["query"].(string)
			rows, err := db.Query(`
				SELECT DISTINCT ps.problem_id, ps.title, ps.difficulty, ps.tags, ps.solution_approach
				FROM problem_solutions ps
				JOIN inverted_index ii ON ps.problem_id = ii.problem_id
				WHERE ii.word LIKE ?
				ORDER BY ii.weight DESC
				LIMIT 3
			`, "%"+query+"%")
			if err != nil {
				resp = MCPResponse{Success: false, Error: err.Error()}
			} else {
				defer rows.Close()
				var results []map[string]interface{}
				for rows.Next() {
					var pid int
					var title, difficulty, tags, approach string
					rows.Scan(&pid, &title, &difficulty, &tags, &approach)
					results = append(results, map[string]interface{}{
						"problem_id": pid,
						"title": title,
						"difficulty": difficulty,
						"tags": tags,
						"approach": approach,
					})
				}
				resp = MCPResponse{Success: true, Data: results}
			}

		case "get_problem":
			pid := int(req.Params["problem_id"].(float64))
			var title, difficulty, tags, approach, code string
			err := db.QueryRow(`
				SELECT title, difficulty, tags, solution_approach, solution_code
				FROM problem_solutions WHERE problem_id = ?
			`, pid).Scan(&title, &difficulty, &tags, &approach, &code)
			if err != nil {
				resp = MCPResponse{Success: false, Error: err.Error()}
			} else {
				resp = MCPResponse{Success: true, Data: map[string]interface{}{
					"title": title,
					"difficulty": difficulty,
					"tags": tags,
					"approach": approach,
					"code": code,
				}}
			}

		case "get_solution_code":
			pid := int(req.Params["problem_id"].(float64))
			var code string
			err := db.QueryRow("SELECT solution_code FROM problem_solutions WHERE problem_id = ?", pid).Scan(&code)
			if err != nil {
				resp = MCPResponse{Success: false, Error: err.Error()}
			} else {
				resp = MCPResponse{Success: true, Data: map[string]string{"code": code}}
			}

		default:
			resp = MCPResponse{Success: false, Error: "unknown type"}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("MCP Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
