# MongoDB 数据库使用指南

本文档详细说明了 Go-Flutter Mall 项目中 MongoDB 的配置、模型设计、常用操作及最佳实践。

## 1. 简介与使用场景

本项目使用 **MongoDB** 作为非关系型数据库，主要用于存储结构灵活、写入量大或需要快速迭代的数据。

### 1.1 核心使用场景
*   **搜索历史 (Search History)**: 用户的搜索关键词记录。由于无需复杂的事务支持，且写入频繁，非常适合用 MongoDB 存储。
*   **日志与审计 (Logs)**: 系统操作日志、API 请求日志（如有）。
*   **聊天记录 (Chat History)**: （可选）海量聊天消息的持久化存储，支持灵活的消息格式（文本、图片、富媒体）。

### 1.2 环境配置 (`backend/config/mongo.go`)

项目启动时会自动连接 MongoDB。

```go
var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongoDB() {
    // 默认连接本地 MongoDB，端口 27017
    uri := "mongodb://127.0.0.1:27017"
    // ... 连接逻辑 ...
    MongoDB = client.Database("go_flutter_mall")
}
```

---

## 2. 数据模型 (BSON)

在 Go 中使用 MongoDB，通常使用 `go.mongodb.org/mongo-driver/bson` 进行数据映射。

### 2.1 示例：搜索历史 (`backend/models/search_history.go`)

```go
package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchHistory struct {
    // _id 是 MongoDB 的主键，使用 primitive.ObjectID 类型
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID    uint               `bson:"user_id" json:"user_id"`
    Keyword   string             `bson:"keyword" json:"keyword"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
```

**关键点**:
*   使用 `bson:"field_name"` 标签定义数据库中的字段名。
*   `_id` 字段通常对应 `primitive.ObjectID`。
*   `omitempty` 表示如果字段为空，则在序列化时不包含该字段（主要用于插入时自动生成 ID）。

---

## 3. 常用操作 (CRUD)

以下代码示例基于 `backend/controllers/search/search_controller.go`。

### 3.1 获取集合 (Collection)

```go
import "go-flutter-mall/backend/config"

// 获取 search_history 集合的句柄
collection := config.MongoDB.Collection("search_history")
```

### 3.2 插入或更新 (Upsert)

**场景**: 用户搜索关键词，如果该词已存在，则更新搜索时间；如果不存在，则插入新记录。

```go
import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func SaveHistory(userID uint, keyword string) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 1. 定义查询条件 (Filter)
    filter := bson.M{"user_id": userID, "keyword": keyword}

    // 2. 定义更新内容 (Update)
    // $set: 设置字段值
    // $setOnInsert: 仅在插入时设置的字段（如创建时间，可选）
    update := bson.M{
        "$set": bson.M{
            "created_at": time.Now(),
        },
    }

    // 3. 启用 Upsert 选项
    opts := options.Update().SetUpsert(true)

    // 4. 执行更新
    _, err := collection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        // handle error
    }
}
```

### 3.3 查询列表 (Find)

**场景**: 获取用户的最近 10 条搜索历史，按时间倒序排列。

```go
func GetHistory(userID uint) ([]models.SearchHistory, error) {
    var history []models.SearchHistory
    ctx := context.Background()

    // 1. 查询条件
    filter := bson.M{"user_id": userID}

    // 2. 排序与限制
    opts := options.Find().
        SetSort(bson.D{{Key: "created_at", Value: -1}}). // -1 表示倒序
        SetLimit(10)

    // 3. 执行查询
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    // 4. 解码结果
    if err = cursor.All(ctx, &history); err != nil {
        return nil, err
    }

    return history, nil
}
```

### 3.4 删除 (Delete)

**场景**: 清空用户的搜索历史。

```go
func ClearHistory(userID uint) error {
    filter := bson.M{"user_id": userID}
    _, err := collection.DeleteMany(context.Background(), filter)
    return err
}
```

---

## 4. 命令行工具与调试

本项目内置了 MongoDB 二进制文件（MacOS ARM64），方便快速启动。

### 4.1 启动 MongoDB 服务
在 `backend` 目录下：

```bash
# 指定数据存储目录 ./data/db
./mongo_bin/bin/mongod --dbpath ./data/db --bind_ip 127.0.0.1 --port 27017
```

### 4.2 使用 Mongo Shell
连接到本地数据库进行调试：

```bash
./mongo_bin/bin/mongos --port 27017
# 注意：如果是旧版 mongo shell，命令可能是 ./mongo
```

**常用 Shell 命令**:
```javascript
// 切换数据库
use go_flutter_mall

// 查看集合
show collections

// 查询所有记录
db.search_history.find()

// 格式化查询
db.search_history.find().pretty()

// 统计数量
db.search_history.countDocuments()

// 删除集合
db.search_history.drop()
```

---

## 5. 最佳实践

1.  **索引 (Indexes)**:
    *   MongoDB 对未索引字段的查询性能较差。
    *   对于 `search_history`，建议为 `user_id` 和 `created_at` 建立复合索引，以加速“查询某用户的最近记录”。
    *   **创建索引示例**:
        ```go
        model := mongo.IndexModel{
            Keys: bson.D{
                {Key: "user_id", Value: 1},
                {Key: "created_at", Value: -1},
            },
        }
        collection.Indexes().CreateOne(context.Background(), model)
        ```

2.  **连接管理**:
    *   `mongo.Client` 是线程安全的，应作为单例使用（本项目已在 `config.MongoClient` 中实现）。
    *   不要为每个请求创建新的 Client。

3.  **BSON 构建**:
    *   `bson.M`: Map，用于构建无序的文档（如 Filter）。
    *   `bson.D`: Slice，用于构建有序的文档（如 Sort, Command），在某些命令中顺序很重要。
    *   `bson.A`: Array，用于构建数组。

4.  **错误处理**:
    *   查询为空时，`FindOne` 会返回 `mongo.ErrNoDocuments`，需要单独处理。
    *   `Find` 返回空列表时不会报错，而是 `cursor` 为空。
