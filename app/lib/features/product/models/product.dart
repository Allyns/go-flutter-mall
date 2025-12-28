/// 商品模型
class Product {
  final int id;
  final int categoryId;
  final String name;
  final String description;
  final double price;
  final int stock;
  final String coverImage;
  final List<String> images;
  final List<ProductSKU> skus;
  final List<Review> reviews;

  Product({
    required this.id,
    required this.categoryId,
    required this.name,
    required this.description,
    required this.price,
    required this.stock,
    required this.coverImage,
    required this.images,
    this.skus = const [],
    this.reviews = const [],
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['ID'], // GORM Model 使用 ID (大写)
      categoryId: json['category_id'] ?? 0,
      name: json['name'],
      description: json['description'],
      price: (json['price'] as num).toDouble(),
      stock: json['stock'],
      coverImage: json['cover_image'],
      // 处理字符串数组
      images: List<String>.from(json['images'] ?? []),
      skus:
          (json['skus'] as List<dynamic>?)
              ?.map((e) => ProductSKU.fromJson(e))
              .toList() ??
          [],
      reviews:
          (json['reviews'] as List<dynamic>?)
              ?.map((e) => Review.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class ProductSKU {
  final int id;
  final String name;
  final String specs;
  final double price;
  final int stock;

  ProductSKU({
    required this.id,
    required this.name,
    required this.specs,
    required this.price,
    required this.stock,
  });

  factory ProductSKU.fromJson(Map<String, dynamic> json) {
    return ProductSKU(
      id: json['ID'],
      name: json['name'],
      specs: json['specs'],
      price: (json['price'] as num).toDouble(),
      stock: json['stock'],
    );
  }
}

class Review {
  final int id;
  final int userId;
  final String userName;
  final String userAvatar;
  final String content;
  final int rating;
  final String createdAt;

  Review({
    required this.id,
    required this.userId,
    required this.userName,
    required this.userAvatar,
    required this.content,
    required this.rating,
    required this.createdAt,
  });

  factory Review.fromJson(Map<String, dynamic> json) {
    final user = json['user'] ?? {};
    return Review(
      id: json['ID'],
      userId: json['user_id'],
      userName: user['username'] ?? '匿名用户',
      userAvatar: user['avatar'] ?? '',
      content: json['content'],
      rating: json['rating'],
      createdAt: json['CreatedAt'] ?? '',
    );
  }
}
