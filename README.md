# MySQL Inspector

MySQL数据库巡检工具，用于检查和监控MySQL数据库的各项指标。

## 功能特性

目前实现的功能：
- 查询表行数Top 10
- 非InnoDB表查询
- 表碎片率Top 10
- Web界面展示

特性说明：
- 自动每5分钟刷新数据
- 响应式布局，支持移动端访问
- 实时展示数据库状态

## 环境要求

- Go 1.16+
- MySQL 5.7+

## 快速开始

### 后端设置

1. 克隆项目
```bash
git clone https://github.com/martinliu007/mysql-inspector.git
cd mysql-inspector
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库连接
复制配置文件示例并修改：
```bash
cp config.json.example config.json
```
修改 `config.json` 中的数据库连接信息。

4. 运行后端服务
```bash
go run cmd/main.go
```

### 前端设置

1. 进入web目录
```bash
cd web
```

2. 启动Web服务器
可以使用任何静态文件服务器，例如：
- Python 3:
```bash
python -m http.server 3000
```
- Node.js (需要安装 http-server):
```bash
npx http-server -p 3000
```

3. 访问Web界面
打开浏览器访问：
```
http://localhost:3000
```

### 使用说明

Web界面分为三个主要部分：
1. 表行数统计：显示数据库中行数最多的10张表
2. 非InnoDB表列表：显示所有非InnoDB存储引擎的表
3. 表碎片统计：显示碎片率最高的10张表

特性：
- 数据自动每5分钟刷新一次
- 可以通过刷新按钮手动更新数据
- 支持响应式布局，适配不同屏幕大小

## API接口

### 获取表行数Top 10

```
GET /api/tables/top
```

响应示例：
```json
[
    {
        "table_name": "users",
        "row_count": 1000000
    },
    {
        "table_name": "orders",
        "row_count": 500000
    }
]
```

### 获取非InnoDB表

```
GET /api/tables/non-innodb
```

响应示例：
```json
[
    {
        "table_name": "old_logs",
        "table_schema": "mydb",
        "engine": "MyISAM",
        "row_count": 50000,
        "data_length": 1048576,
        "index_length": 524288
    },
    {
        "table_name": "archive_data",
        "table_schema": "mydb",
        "engine": "ARCHIVE",
        "row_count": 10000,
        "data_length": 2097152,
        "index_length": 0
    }
]
```

该接口返回所有非InnoDB存储引擎的表，包括：
- 表名和所属数据库
- 存储引擎类型
- 表行数
- 数据长度（字节）
- 索引长度（字节）

### 获取表碎片率Top 10

```
GET /api/tables/fragmentation
```

响应示例：
```json
[
    {
        "table_name": "large_table",
        "table_schema": "mydb",
        "engine": "InnoDB",
        "data_free": 1048576,
        "data_length": 5242880,
        "fragment_ratio": 20.0
    },
    {
        "table_name": "history_data",
        "table_schema": "mydb",
        "engine": "InnoDB",
        "data_free": 524288,
        "data_length": 3145728,
        "fragment_ratio": 16.67
    }
]
```

该接口返回InnoDB表的碎片信息，包括：
- 表名和所属数据库
- 存储引擎类型
- 空闲空间大小（字节）
- 数据长度（字节）
- 碎片率（百分比）

注意：
- 只返回InnoDB引擎的表
- 碎片率 = (空闲空间 / 数据长度) * 100
- 按碎片率降序排序
- 只返回有碎片的表（DATA_FREE > 0）

## 配置说明

配置文件 `config.json` 参数说明：

```json
{
    "server_port": 8080,        // HTTP服务器端口
    "mysql": {
        "host": "localhost",    // MySQL主机地址
        "port": 3306,          // MySQL端口
        "user": "root",        // MySQL用户名
        "password": "******",  // MySQL密码
        "database": "test"     // 要检查的数据库名
    }
}
```

## 注意事项

- 确保MySQL用户具有information_schema的查询权限
- 对于大型数据库，表行数统计可能需要一定时间
