package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// 分词器：提取中文和英文单词
func tokenize(text string) []string {
	var tokens []string
	// 匹配中文
	chineseRe := regexp.MustCompile(`[\u4e00-\u9fa5]+`)
	chineseMatches := chineseRe.FindAllString(text, -1)
	tokens = append(tokens, chineseMatches...)
	
	// 匹配英文单词
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for _, w := range words {
		w = strings.ToLower(strings.TrimSpace(w))
		if len(w) > 1 {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

// 构建倒排索引
func buildInvertedIndex(db *sql.DB) error {
	// 清空现有索引
	_, err := db.Exec("DELETE FROM inverted_index")
	if err != nil {
		return err
	}

	// 获取所有题目解决方案
	rows, err := db.Query(`
		SELECT problem_id, title, description, solution_approach, solution_code, tags 
		FROM problem_solutions
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	stmt, err := db.Prepare(`
		INSERT INTO inverted_index (word, problem_id, field_type, position, weight)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for rows.Next() {
		var pid int
		var title, desc, approach, code, tags string
		err := rows.Scan(&pid, &title, &desc, &approach, &code, &tags)
		if err != nil {
			continue
		}

		// 为不同字段建立索引，设置不同权重
		fields := map[string]struct {
			content string
			weight  float64
		}{
			"title":       {title, 3.0},
			"tags":        {tags, 2.5},
			"description": {desc, 2.0},
			"solution":    {approach, 1.5},
			"code":        {code, 1.0},
		}

		for fieldType, field := range fields {
			tokens := tokenize(field.content)
			for pos, word := range tokens {
				_, err := stmt.Exec(word, pid, fieldType, pos, field.weight)
				if err != nil {
					log.Printf("Error inserting token: %v", err)
				}
			}
		}
	}

	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "oj_platform.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Building inverted index...")
	if err := buildInvertedIndex(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inverted index built successfully!")
}
