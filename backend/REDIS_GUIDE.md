# Redis 使用指南

本文档详细说明了 Go-Flutter Mall 项目中 Redis 的配置、使用场景、核心模式及最佳实践。

## 1. 简介与环境配置

本项目使用 **Redis** 作为高性能键值存储，主要用于解决高并发场景下的数据一致性问题（分布式锁）和异步延时任务（订单超时）。

### 1.1 依赖库
使用官方推荐的 Go 客户端：
- `github.com/redis/go-redis/v9`

### 1.2 连接配置 (`backend/config/redis.go`)
```go
var RedisClient *redis.Client

func ConnectRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     "127.0.0.1:6379",
        Password: "", // 生产环境请设置密码
        DB:       0,
    })
    // ... Ping 测试连接
}
```

---

## 2. 核心使用场景

### 2.1 分布式锁 (Distributed Lock)

**场景**: 防止商品超卖。当多个用户同时购买同一件商品时，我们需要保证库存扣减的原子性。

**实现原理**: 使用 Redis 的 `SETNX` (Set if Not Exists) 命令。
- **Key**: `lock:product:{id}`
- **Value**: 占位符 (如 "1")
- **Expiration**: 防止死锁 (如 5秒)

**代码示例 (`controllers/order/order_controller.go`)**:

```go
// 尝试获取锁
lockKey := fmt.Sprintf("lock:product:%d", item.ProductID)
ctx := context.Background()

// SET lock:product:1 1 NX EX 5
// 成功返回 true，失败返回 false
locked, err := config.RedisClient.SetNX(ctx, lockKey, 1, 5*time.Second).Result()

if err != nil {
    // Redis 连接错误处理
    log.Println("Redis error:", err)
} else if !locked {
    // 获取锁失败，说明有其他请求正在处理该商品
    return errors.New("系统繁忙，请重试")
} else {
    // 业务处理完成后释放锁
    defer config.RedisClient.Del(ctx, lockKey)
}

// ... 执行扣库存逻辑 ...
```

### 2.2 延时队列 (Delay Queue)

**场景**: 订单创建后，如果 30 分钟内未支付，系统自动取消订单并回滚库存。

**实现原理**: 使用 Redis 的 **Sorted Set (ZSet)**。
- **Key**: `order:delay_queue`
- **Member**: 任务标识 (如 `orderID:userID`)
- **Score**: 执行时间的时间戳 (Unix Timestamp)

**生产者：添加任务 (`pkg/scheduler/scheduler.go`)**:

```go
const DelayQueueKey = "order:delay_queue"

func AddToDelayQueue(orderID uint, userID uint, delay time.Duration) error {
    ctx := context.Background()
    // Score = 当前时间 + 延时时长 (e.g. 30分钟)
    score := float64(time.Now().Add(delay).Unix())
    member := fmt.Sprintf("%d:%d", orderID, userID)

    // ZADD order:delay_queue <timestamp> <orderID:userID>
    return config.RedisClient.ZAdd(ctx, DelayQueueKey, redis.Z{
        Score:  score,
        Member: member,
    }).Err()
}
```

**消费者：轮询任务 (`pkg/scheduler/scheduler.go`)**:

```go
func StartScheduler() {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        for range ticker.C {
            ctx := context.Background()
            now := float64(time.Now().Unix())

            // 1. 获取已到期的任务 (Score <= Now)
            // ZRANGEBYSCORE order:delay_queue -inf <now> LIMIT 0 10
            vals, _ := config.RedisClient.ZRangeByScore(ctx, DelayQueueKey, &redis.ZRangeBy{
                Min:    "-inf",
                Max:    fmt.Sprintf("%f", now),
                Offset: 0,
                Count:  10,
            }).Result()

            if len(vals) > 0 {
                for _, val := range vals {
                    // 2. 移除任务 (防止重复处理)
                    // ZREM order:delay_queue <val>
                    removed, _ := config.RedisClient.ZRem(ctx, DelayQueueKey, val).Result()
                    
                    if removed > 0 {
                        // 3. 触发业务逻辑 (如发送 Kafka 事件取消订单)
                        handleOrderTimeout(val)
                    }
                }
            }
        }
    }()
}
```

### 2.3 数据缓存 (Caching) - *推荐实践*

**场景**: 首页轮播图、商品详情等高频读取、低频修改的数据。

**代码示例**:

```go
func GetProductWithCache(id uint) (*models.Product, error) {
    ctx := context.Background()
    cacheKey := fmt.Sprintf("cache:product:%d", id)

    // 1. 尝试从 Redis 获取
    val, err := config.RedisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        var product models.Product
        json.Unmarshal([]byte(val), &product)
        return &product, nil
    }

    // 2. Redis 未命中，查询数据库
    var product models.Product
    if err := config.DB.First(&product, id).Error; err != nil {
        return nil, err
    }

    // 3. 写入 Redis (设置 1 小时过期)
    jsonBytes, _ := json.Marshal(product)
    config.RedisClient.Set(ctx, cacheKey, jsonBytes, 1*time.Hour)

    return &product, nil
}
```

---

## 3. 常用 CLI 命令调试

在开发过程中，可以使用 `redis-cli` 工具查看数据。

```bash
# 连接 Redis
redis-cli

# 查看所有 Key (慎用，生产环境禁用)
KEYS *

# 1. 检查分布式锁
GET lock:product:123
TTL lock:product:123  # 查看剩余过期时间

# 2. 检查延时队列
ZRANGE order:delay_queue 0 -1 WITHSCORES

# 3. 手动清空数据库
FLUSHDB
```

---

## 4. 最佳实践

1.  **Key 命名规范**:
    *   使用冒号分隔层级: `service:module:id`
    *   示例: `mall:auth:token:1001`, `mall:product:view_count:50`
2.  **设置过期时间 (TTL)**:
    *   所有的缓存 Key 都应该有过期时间，防止内存无限膨胀。
    *   分布式锁必须有过期时间，防止程序崩溃导致死锁。
3.  **错误处理**:
    *   Redis 操作可能会失败（网络波动）。
    *   对于非关键业务（如缓存），Redis 失败应降级查库，不应导致接口报错。
    *   对于关键业务（如分布式锁），Redis 失败通常意味着服务降级或拒绝请求。
4.  **连接池**:
    *   `go-redis` 客户端内部已实现了连接池，通常无需手动配置，但可通过 `Options` 调整 `PoolSize`。

## 5. 常见问题排查

*   **Error: `connect: connection refused`**:
    *   检查 Redis 服务是否启动 (`docker ps` 或 `brew services list`)。
    *   检查配置中的 IP 端口是否正确。
*   **延时任务未触发**:
    *   检查系统时间是否准确。
    *   检查 `StartScheduler()` 是否在 `main.go` 中被调用。
