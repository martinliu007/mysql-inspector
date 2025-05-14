package database

import (
	"database/sql"
	"fmt"
	"mysql-inspector/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// TableInfo 存储表信息
type TableInfo struct {
	TableName   string `json:"table_name"`
	RowCount    int64  `json:"row_count"`
	Engine      string `json:"engine,omitempty"`
	TableSchema string `json:"table_schema,omitempty"`
	DataLength  int64  `json:"data_length,omitempty"`
	IndexLength int64  `json:"index_length,omitempty"`
}

// GetNonInnoDBTables 获取所有表的引擎类型（包括InnoDB和非InnoDB）
func GetNonInnoDBTables(db *sql.DB) ([]TableInfo, error) {
	query := `
		SELECT 
			TABLE_SCHEMA,
			TABLE_NAME,
			ENGINE,
			TABLE_ROWS,
			DATA_LENGTH,
			INDEX_LENGTH
		FROM 
			information_schema.TABLES 
		WHERE 
			TABLE_SCHEMA = DATABASE()
		ORDER BY 
			ENGINE != 'InnoDB', TABLE_NAME
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		if err := rows.Scan(
			&table.TableSchema,
			&table.TableName,
			&table.Engine,
			&table.RowCount,
			&table.DataLength,
			&table.IndexLength,
		); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// FragmentationInfo 存储表碎片信息
type FragmentationInfo struct {
	TableName     string  `json:"table_name"`
	TableSchema   string  `json:"table_schema"`
	Engine        string  `json:"engine"`
	DataFree      int64   `json:"data_free"`
	DataLength    int64   `json:"data_length"`
	FragmentRatio float64 `json:"fragment_ratio"`
	FragmentSize  int64   `json:"fragment_size"`
}

// GetTopFragmentedTables 获取所有表的碎片信息
func GetTopFragmentedTables(db *sql.DB, limit int) ([]FragmentationInfo, error) {
	query := `
		SELECT 
			TABLE_SCHEMA,
			TABLE_NAME,
			ENGINE,
			DATA_FREE,
			DATA_LENGTH,
			ROUND(IFNULL(DATA_FREE * 100.0 / NULLIF(DATA_LENGTH, 0), 0), 2) as FRAGMENT_RATIO
		FROM 
			information_schema.TABLES 
		WHERE 
			TABLE_SCHEMA = DATABASE()
			AND DATA_LENGTH > 0
		ORDER BY 
			FRAGMENT_RATIO DESC, DATA_LENGTH DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []FragmentationInfo
	for rows.Next() {
		var table FragmentationInfo
		if err := rows.Scan(
			&table.TableSchema,
			&table.TableName,
			&table.Engine,
			&table.DataFree,
			&table.DataLength,
			&table.FragmentRatio,
		); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// NewMySQLConnection 创建新的MySQL连接
func NewMySQLConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// GetTopTablesByRowCount 获取行数最多的前N张表
func GetTopTablesByRowCount(db *sql.DB, limit int) ([]TableInfo, error) {
	query := `
		SELECT 
			TABLE_NAME,
			TABLE_ROWS
		FROM 
			information_schema.TABLES
		WHERE 
			TABLE_SCHEMA = DATABASE()
		ORDER BY 
			TABLE_ROWS DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		if err := rows.Scan(&table.TableName, &table.RowCount); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
