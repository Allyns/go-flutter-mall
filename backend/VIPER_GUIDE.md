# Viper 配置管理使用指南

本文档详细说明了如何在 Go-Flutter Mall 项目中集成和使用 **Viper** 进行配置管理。Viper 是 Go 语言中最流行的配置解决方案，支持多种格式（JSON, YAML, TOML 等）和环境变量覆盖。

## 1. 为什么使用 Viper？

当前项目中，数据库连接串、密钥等配置可能硬编码在代码中（如 `config/database.go`）。使用 Viper 可以解决以下问题：
*   **配置分离**: 将配置从代码中分离到配置文件（如 `config.yaml`）。
*   **环境区分**: 轻松切换开发、测试、生产环境配置。
*   **安全**: 敏感信息（密码、Secret）可通过环境变量注入，不提交到代码仓库。
*   **默认值**: 为配置项提供合理的默认值。

## 2. 安装与集成

### 2.1 安装依赖
在 `backend` 目录下执行：
```bash
go get github.com/spf13/viper
```

### 2.2 配置文件设计 (`config.yaml`)
建议在 `backend` 根目录或 `config` 目录下创建 `config.yaml`。

```yaml
server:
  port: "8080"
  mode: "debug" # debug, release

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres" # 生产环境请使用环境变量覆盖
  dbname: "go_flutter_mall"
  sslmode: "disable"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

jwt:
  secret: "your_jwt_secret_key"
  expire_hours: 24
```

---

## 3. 代码实现

建议创建一个统一的配置加载模块。

### 3.1 定义配置结构体 (`config/config.go`)

使用结构体可以获得 IDE 的代码提示和类型检查。

```go
package config

import (
    "log"
    "github.com/spf13/viper"
)

// GlobalConfig 全局配置实例
var AppConfig Config

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
    Port string `mapstructure:"port"`
    Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
    Addr     string `mapstructure:"addr"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
    Secret      string `mapstructure:"secret"`
    ExpireHours int    `mapstructure:"expire_hours"`
}

// LoadConfig 加载配置
func LoadConfig() {
    viper.SetConfigName("config") // 配置文件名 (不带后缀)
    viper.SetConfigType("yaml")   // 配置文件类型
    viper.AddConfigPath(".")      // 搜索路径 (当前目录)
    viper.AddConfigPath("./config") 

    // 设置环境变量前缀 (例如: MALL_DATABASE_PASSWORD)
    viper.SetEnvPrefix("MALL") 
    viper.AutomaticEnv() // 自动读取环境变量

    // 读取配置文件
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // 配置文件未找到，可能仅依赖环境变量
            log.Println("Config file not found, using defaults/environment variables")
        } else {
            log.Fatalf("Error reading config file: %s", err)
        }
    }

    // 解析到结构体
    if err := viper.Unmarshal(&AppConfig); err != nil {
        log.Fatalf("Unable to decode into struct: %s", err)
    }
    
    log.Println("Configuration loaded successfully")
}
```

---

## 4. 使用场景与示例

### 4.1 场景一：初始化数据库连接

**修改前** (`config/database.go`):
```go
dsn := "host=localhost user=postgres ..."
```

**修改后**:
```go
import "fmt"

func ConnectDatabase() {
    cfg := AppConfig.Database
    
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
        cfg.Host,
        cfg.User,
        cfg.Password,
        cfg.DBName,
        cfg.Port,
        cfg.SSLMode,
    )

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    // ...
}
```

### 4.2 场景二：JWT 密钥配置

**修改前** (`utils/jwt.go`):
```go
var jwtSecret = []byte("hardcoded_secret")
```

**修改后**:
```go
import "go-flutter-mall/backend/config"

// 注意：需要在 main.go 中先调用 config.LoadConfig()
func GenerateToken(userID uint) (string, error) {
    secret := []byte(config.AppConfig.JWT.Secret)
    // ...
}
```

### 4.3 场景三：程序入口初始化

在 `backend/main.go` 的最开始加载配置。

```go
func main() {
    // 1. 加载配置 (必须是第一步)
    config.LoadConfig()

    // 2. 初始化各个组件
    config.ConnectDatabase()
    config.ConnectRedis()
    
    // ...
    
    // 3. 启动服务
    r.Run(":" + config.AppConfig.Server.Port)
}
```

---

## 5. 高级用法：环境变量覆盖 (Docker/K8s 部署)

Viper 强大的地方在于可以无缝支持环境变量覆盖，这在容器化部署中非常关键。

假设我们设置了 `viper.SetEnvPrefix("MALL")` 并开启了 `AutomaticEnv()`，Viper 会自动将环境变量映射到配置项（使用 `_` 替换 `.`）。

**示例**:
如果不修改代码，只想在生产环境中改变数据库密码：

```bash
# 本地运行
export MALL_DATABASE_PASSWORD="prod_password_123"
export MALL_SERVER_PORT="9090"

./backend
```

*   `MALL_DATABASE_PASSWORD` -> 覆盖 `database.password`
*   `MALL_SERVER_PORT` -> 覆盖 `server.port`

---

## 6. 最佳实践

1.  **不要提交 config.yaml**:
    *   通常将 `config.yaml` 加入 `.gitignore`。
    *   提供一个 `config.example.yaml` 作为模板提交到仓库。

2.  **默认值**:
    *   可以使用 `viper.SetDefault("server.port", "8080")` 设置合理的默认值，防止配置文件缺失导致程序崩溃。

3.  **热加载 (WatchConfig)**:
    *   Viper 支持在运行时监控配置文件变化并自动重新加载（适用于无需重启即可生效的配置，如日志级别）。
    *   ```go
        viper.WatchConfig()
        viper.OnConfigChange(func(e fsnotify.Event) {
            fmt.Println("Config file changed:", e.Name)
            viper.Unmarshal(&AppConfig) // 重新解析
        })
        ```
