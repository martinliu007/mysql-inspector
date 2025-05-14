package handler

import (
	"database/sql"
	"encoding/json"
	"mysql-inspector/internal/database"
	"net/http"
)

// Inspector 处理检查相关的HTTP请求
type Inspector struct {
	db *sql.DB
}

// NewInspector 创建新的Inspector实例
func NewInspector(db *sql.DB) *Inspector {
	return &Inspector{db: db}
}

// setCORSHeaders 设置CORS响应头
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// GetTopTables 处理获取表行数top N的请求
func (i *Inspector) GetTopTables(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 默认返回前10张表
	limit := 10

	// 获取表信息
	tables, err := database.GetTopTablesByRowCount(i.db, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理OPTIONS请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理OPTIONS请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理OPTIONS请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(tables); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetNonInnoDBTables 处理获取表引擎类型的请求
func (i *Inspector) GetNonInnoDBTables(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取非InnoDB表信息
	tables, err := database.GetNonInnoDBTables(i.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(tables); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTopFragmentedTables 处理获取碎片率最高表的请求
func (i *Inspector) GetTopFragmentedTables(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 默认返回前10张表
	limit := 10

	// 获取碎片率信息
	tables, err := database.GetTopFragmentedTables(i.db, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(tables); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
