package main

import (
	"fmt"
	"log"
	"net/http"

	"mysql-inspector/internal/config"
	"mysql-inspector/internal/database"
	"mysql-inspector/internal/handler"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建处理器
	h := handler.NewInspector(db)

	// 注册API路由
	http.HandleFunc("/api/tables/top", h.GetTopTables)
	http.HandleFunc("/api/tables/non-innodb", h.GetNonInnoDBTables)
	http.HandleFunc("/api/tables/fragmentation", h.GetTopFragmentedTables)

	// 添加静态文件服务
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
		http.NotFound(w, r)
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
