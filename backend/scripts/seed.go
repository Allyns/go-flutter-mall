package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"go-flutter-mall/backend/config"
	"go-flutter-mall/backend/models"
	"go-flutter-mall/backend/utils"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// main 执行数据填充脚本
// 运行方法: cd backend && go run scripts/seed.go
func main() {
	// 初始化数据库连接
	config.ConnectDatabase()
	db := config.DB

	log.Println("🌱 开始填充数据...")

	// 1. 清理现有数据
	log.Println("正在清理旧数据...")
	db.Exec("TRUNCATE TABLE reviews, order_items, orders, addresses, cart_items, product_skus, products, categories, admin_users, users RESTART IDENTITY CASCADE")

	// 2. 创建管理员
	adminPassword, _ := utils.HashPassword("admin123")
	admin := models.AdminUser{
		Username: "admin",
		Password: adminPassword,
		Role:     "admin",
		Avatar:   "https://ui-avatars.com/api/?name=Admin&background=random",
	}
	db.Create(&admin)
	log.Println("已创建管理员: admin / admin123")

	// 3. 创建商品分类
	digital := models.Category{Name: "数码", SortOrder: 1, Icon: "phone_iphone"}
	clothing := models.Category{Name: "服饰", SortOrder: 2, Icon: "checkroom"}
	food := models.Category{Name: "食品", SortOrder: 3, Icon: "restaurant"}
	fresh := models.Category{Name: "生鲜", SortOrder: 4, Icon: "local_florist"}
	appliances := models.Category{Name: "家电", SortOrder: 5, Icon: "kitchen"}

	db.Create(&digital)
	db.Create(&clothing)
	db.Create(&food)
	db.Create(&fresh)
	db.Create(&appliances)

	// 图片链接 (已修复失效链接)
	const imgIphone = "https://images.unsplash.com/photo-1695048133142-1a20484d2569?q=80&w=800&auto=format&fit=crop"
	const imgHeadphone = "https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?q=80&w=800&auto=format&fit=crop"
	const imgTshirt = "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?q=80&w=800&auto=format&fit=crop"
	const imgJacket = "https://images.unsplash.com/photo-1523275335684-37898b6baf30?q=80&w=800&auto=format&fit=crop"
	const imgSalad = "https://images.unsplash.com/photo-1546069901-ba9599a7e63c?q=80&w=800&auto=format&fit=crop"
	const imgFruit = "https://images.unsplash.com/photo-1610832958506-aa56368176cf?q=80&w=800&auto=format&fit=crop"
	const imgSeafood = "https://images.unsplash.com/photo-1534483509719-3feaee7c30da?q=80&w=800&auto=format&fit=crop"
	const imgFridge = "https://images.unsplash.com/photo-1588854337221-4cf9fa96059c?q=80&w=800&auto=format&fit=crop"
	const imgWasher = "https://images.unsplash.com/photo-1626806819282-2c1dc01a5e0c?q=80&w=800&auto=format&fit=crop"
	const imgLamp = "https://images.unsplash.com/photo-1565814329452-e1efa11c5b89?q=80&w=800&auto=format&fit=crop"
	const imgKeyboard = "https://images.unsplash.com/photo-1595225476474-87563907a212?q=80&w=800&auto=format&fit=crop"
	const imgShoes = "https://images.unsplash.com/photo-1542291026-7eec264c27ff?q=80&w=800&auto=format&fit=crop"

	// 4. 创建商品列表
	products := []models.Product{
		// 数码
		{CategoryID: digital.ID, Name: "iPhone 15 Pro Max", Description: "钛金属设计，A17 Pro 芯片，史上最强大的 iPhone。", Price: 9999.00, Stock: 100, CoverImage: imgIphone, Images: pq.StringArray{imgIphone}, Status: 1},
		{CategoryID: digital.ID, Name: "Sony WH-1000XM5", Description: "行业领先的降噪耳机，配备自动降噪优化器。", Price: 2499.00, Stock: 50, CoverImage: imgHeadphone, Images: pq.StringArray{imgHeadphone}, Status: 1},
		{CategoryID: digital.ID, Name: "机械键盘 RGB", Description: "RGB 背光，红轴，紧凑设计，打字手感极佳。", Price: 499.00, Stock: 150, CoverImage: imgKeyboard, Images: pq.StringArray{imgKeyboard}, Status: 1},
		// 服饰
		{CategoryID: clothing.ID, Name: "经典纯棉T恤", Description: "优质纯棉，透气舒适，百搭款式。", Price: 99.00, Stock: 200, CoverImage: imgTshirt, Images: pq.StringArray{imgTshirt}, Status: 1},
		{CategoryID: clothing.ID, Name: "复古牛仔夹克", Description: "经典款式牛仔夹克，适合任何季节穿着。", Price: 399.00, Stock: 80, CoverImage: imgJacket, Images: pq.StringArray{imgJacket}, Status: 1},
		{CategoryID: clothing.ID, Name: "专业跑步鞋", Description: "轻量化设计，减震鞋底，完美适合慢跑和训练。", Price: 599.00, Stock: 120, CoverImage: imgShoes, Images: pq.StringArray{imgShoes}, Status: 1},
		// 食品
		{CategoryID: food.ID, Name: "健康沙拉碗", Description: "新鲜蔬菜搭配特制酱料，健康美味。", Price: 35.00, Stock: 999, CoverImage: imgSalad, Images: pq.StringArray{imgSalad}, Status: 1},
		// 生鲜
		{CategoryID: fresh.ID, Name: "进口甜橙 (5kg)", Description: "阳光充足，果肉饱满，汁多味甜。", Price: 88.00, Stock: 300, CoverImage: imgFruit, Images: pq.StringArray{imgFruit}, Status: 1},
		{CategoryID: fresh.ID, Name: "新鲜三文鱼切片", Description: "深海捕捞，极速冷链，口感鲜美。", Price: 128.00, Stock: 50, CoverImage: imgSeafood, Images: pq.StringArray{imgSeafood}, Status: 1},
		// 家电
		{CategoryID: appliances.ID, Name: "智能双开门冰箱", Description: "大容量，风冷无霜，智能温控。", Price: 3999.00, Stock: 20, CoverImage: imgFridge, Images: pq.StringArray{imgFridge}, Status: 1},
		{CategoryID: appliances.ID, Name: "全自动滚筒洗衣机", Description: "洗烘一体，静音变频，除菌洗。", Price: 2599.00, Stock: 30, CoverImage: imgWasher, Images: pq.StringArray{imgWasher}, Status: 1},
		{CategoryID: appliances.ID, Name: "现代护眼台灯", Description: "LED 护眼台灯，可调节亮度和色温。", Price: 159.00, Stock: 300, CoverImage: imgLamp, Images: pq.StringArray{imgLamp}, Status: 1},
	}

	var savedProducts []models.Product
	for _, p := range products {
		if err := db.Create(&p).Error; err != nil {
			log.Printf("创建商品失败 %s: %v", p.Name, err)
			continue
		}
		// SKU
		db.Create(&models.ProductSKU{
			ProductID: p.ID,
			Name:      p.Name + " - 标准版",
			Specs:     `{"type": "标准版"}`,
			Price:     p.Price,
			Stock:     p.Stock,
			Image:     p.CoverImage,
		})
		savedProducts = append(savedProducts, p)
		log.Printf("已创建商品: %s", p.Name)
	}

	// 5. 创建用户 (1个主测试用户 + 10个随机用户)
	userPassword, _ := utils.HashPassword("123456")
	mainUser := models.User{
		Username: "user",
		Email:    "user@example.com",
		Password: userPassword,
		Avatar:   "https://ui-avatars.com/api/?name=User&background=random",
	}
	db.Create(&mainUser)

	// 创建地址
	db.Create(&models.Address{
		UserID: mainUser.ID, ReceiverName: "张三", Phone: "13800138000", Province: "北京市", City: "北京市", District: "朝阳区", DetailAddress: "三里屯 SOHO", IsDefault: true,
	})

	var users []models.User
	users = append(users, mainUser)

	for i := 0; i < 10; i++ {
		u := models.User{
			Username: fmt.Sprintf("user%d", i+1),
			Email:    fmt.Sprintf("user%d@example.com", i+1),
			Password: userPassword,
			Avatar:   fmt.Sprintf("https://ui-avatars.com/api/?name=User%d&background=random", i+1),
		}
		db.Create(&u)
		users = append(users, u)
		// 地址
		db.Create(&models.Address{
			UserID: u.ID, ReceiverName: fmt.Sprintf("用户%d", i+1), Phone: fmt.Sprintf("1390000%04d", i), Province: "上海市", City: "上海市", District: "浦东新区", DetailAddress: "陆家嘴环路 100 号", IsDefault: true,
		})
	}
	log.Println("已创建 11 个用户")

	// 6. 生成大量评论
	comments := []string{
		"非常喜欢，质量很好！", "物流很快，第二天就到了。", "包装有点简陋，但东西不错。", "性价比很高，推荐购买。",
		"不太满意，颜色有色差。", "客服态度很好，解决了我的问题。", "第二次购买了，一如既往的好。", "这是送给朋友的礼物，他很喜欢。",
		"功能很强大，完全符合预期。", "一般般吧，习惯好评。", "真的是物超所值！", "有点小贵，但品质对得起价格。",
	}

	for _, p := range savedProducts {
		// 每个商品生成 3-8 条评论
		count := rand.Intn(6) + 3
		for i := 0; i < count; i++ {
			randomUser := users[rand.Intn(len(users))]
			randomTime := time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour) // 过去30天内

			db.Create(&models.Review{
				UserID:    randomUser.ID,
				ProductID: p.ID,
				Content:   comments[rand.Intn(len(comments))],
				Rating:    rand.Intn(3) + 3, // 3-5 星
				Status:    1,
				Model:     gorm.Model{CreatedAt: randomTime},
			})
		}
	}
	log.Println("已生成商品评论")

	// 7. 生成订单 (覆盖所有状态)
	// 状态: 0:待支付, 1:待发货, 2:待收货, 3:待评价, 4:已完成, 5:售后中, -1:已取消
	statuses := []int{0, 1, 2, 3, 4, 5, -1}

	for _, status := range statuses {
		// 每个状态生成 2-3 个订单
		count := rand.Intn(2) + 2
		for i := 0; i < count; i++ {
			// 随机选一个用户 (主要是主用户，方便查看)
			targetUser := mainUser
			if rand.Float32() > 0.7 {
				targetUser = users[rand.Intn(len(users))]
			}

			// 随机选 1-3 个商品
			itemCount := rand.Intn(3) + 1
			var orderItems []models.OrderItem
			var totalAmount float64

			for j := 0; j < itemCount; j++ {
				p := savedProducts[rand.Intn(len(savedProducts))]
				qty := rand.Intn(2) + 1
				price := p.Price
				totalAmount += price * float64(qty)

				orderItems = append(orderItems, models.OrderItem{
					ProductID:    p.ID,
					ProductName:  p.Name,
					ProductImage: p.CoverImage,
					Price:        price,
					Quantity:     qty,
				})
			}

			// 随机时间
			createdAt := time.Now().Add(-time.Duration(rand.Intn(7*24)) * time.Hour)

			order := models.Order{
				CreatedAt:   createdAt, // 修正
				OrderNo:     fmt.Sprintf("%d%d", createdAt.UnixNano(), targetUser.ID),
				UserID:      targetUser.ID,
				TotalAmount: totalAmount,
				Status:      status,
				AddressID:   1, // 简化，假设都有 AddressID 1 (或者查询该用户的地址)
				Items:       orderItems,
			}

			// 查找用户真实地址 ID
			var addr models.Address
			if err := db.Where("user_id = ?", targetUser.ID).First(&addr).Error; err == nil {
				order.AddressID = addr.ID
			}

			db.Create(&order)
		}
	}
	log.Println("已生成各状态订单数据")

	log.Println("✅ 所有数据填充完成")
}
