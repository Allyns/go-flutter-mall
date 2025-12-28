import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';
import 'package:go_flutter_mall/features/cart/models/cart_item.dart';

/// 购物车 Notifier
/// 管理购物车的加载、添加、更新和删除
class CartNotifier extends StateNotifier<AsyncValue<List<CartItem>>> {
  CartNotifier() : super(const AsyncValue.loading()) {
    fetchCart();
  }

  /// 获取购物车列表
  Future<void> fetchCart() async {
    try {
      state = const AsyncValue.loading();
      final response = await HttpClient().dio.get('/cart');
      final List<dynamic> data = response.data;
      final items = data.map((json) => CartItem.fromJson(json)).toList();
      state = AsyncValue.data(items);
    } catch (e, stack) {
      state = AsyncValue.error(e, stack);
    }
  }

  /// 添加商品到购物车
  Future<void> addToCart(int productId, int quantity, {int? skuId}) async {
    try {
      await HttpClient().dio.post('/cart', data: {
        'product_id': productId,
        'quantity': quantity,
        if (skuId != null) 'sku_id': skuId,
      });
      // 重新获取购物车数据以更新 UI
      fetchCart();
    } catch (e) {
      // 处理错误，如提示用户
      rethrow;
    }
  }

  /// 更新购物车项 (如修改数量、选中状态)
  Future<void> updateCartItem(int id, int quantity, bool selected) async {
    try {
      // 乐观更新 UI (先更新本地状态，再请求 API)
      state.whenData((items) {
        state = AsyncValue.data(items.map((item) {
          if (item.id == id) {
            return CartItem(
              id: item.id,
              productId: item.productId,
              product: item.product,
              quantity: quantity,
              selected: selected,
            );
          }
          return item;
        }).toList());
      });

      await HttpClient().dio.put('/cart/$id', data: {
        'quantity': quantity,
        'selected': selected,
      });
      // 也可以选择在这里重新 fetchCart 确保一致性
    } catch (e) {
      // 如果失败，回滚状态 (需重新 fetch)
      fetchCart();
    }
  }

  /// 删除购物车项
  Future<void> deleteCartItem(int id) async {
    try {
       state.whenData((items) {
        state = AsyncValue.data(items.where((item) => item.id != id).toList());
      });
      
      await HttpClient().dio.delete('/cart/$id');
    } catch (e) {
      fetchCart();
    }
  }
}

/// 购物车 Provider
final cartProvider = StateNotifierProvider<CartNotifier, AsyncValue<List<CartItem>>>((ref) {
  return CartNotifier();
});

/// 计算购物车总价的 Provider
final cartTotalProvider = Provider<double>((ref) {
  final cartState = ref.watch(cartProvider);
  return cartState.maybeWhen(
    data: (items) => items
        .where((item) => item.selected)
        .fold(0, (sum, item) => sum + (item.product.price * item.quantity)),
    orElse: () => 0.0,
  );
});
