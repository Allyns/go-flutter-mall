import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';
import 'package:go_flutter_mall/features/product/models/product.dart';

/// 搜索关键字 Provider
final productSearchQueryProvider = StateProvider<String>((ref) => '');

/// 分类 ID Provider (0 表示全部)
final productCategoryProvider = StateProvider<int>((ref) => 0);

/// 商品列表 Provider (FutureProvider)
/// 自动处理异步加载、缓存和错误状态
final productListProvider = FutureProvider<List<Product>>((ref) async {
  final searchQuery = ref.watch(productSearchQueryProvider);
  final categoryId = ref.watch(productCategoryProvider);

  final queryParams = <String, dynamic>{};
  if (searchQuery.isNotEmpty) {
    queryParams['search'] = searchQuery;
  }
  // 假设后端支持 category_id 筛选，如果后端还没实现，这里暂时只是客户端筛选或忽略
  // 为了演示，我们先加上参数，如果后端忽略也没关系
  if (categoryId > 0) {
    queryParams['category_id'] = categoryId;
  }

  final response = await HttpClient().dio.get(
    '/products',
    queryParameters: queryParams.isNotEmpty ? queryParams : null,
  );

  final List<dynamic> data = response.data;
  var products = data.map((json) => Product.fromJson(json)).toList();

  // 如果后端未实现分类筛选，我们在前端做一层过滤 (临时方案)
  if (categoryId > 0) {
    products = products.where((p) => p.categoryId == categoryId).toList();
  }

  return products;
});

/// 单个商品详情 Provider Family
/// 接收商品 ID 作为参数，返回该商品的详情
final productDetailProvider =
    FutureProvider.family<Product, int>((ref, id) async {
  final response = await HttpClient().dio.get('/products/$id');
  return Product.fromJson(response.data);
});

/// 商品评价列表 Provider
final productReviewsProvider =
    FutureProvider.family<List<Review>, int>((ref, id) async {
  final response = await HttpClient().dio.get('/products/$id/reviews');
  final List<dynamic> data = response.data;
  return data.map((json) => Review.fromJson(json)).toList();
});
