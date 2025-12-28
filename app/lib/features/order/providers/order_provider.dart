import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';
import 'package:go_flutter_mall/features/order/models/order.dart';

/// 订单筛选状态 Provider (null 表示全部)
final orderStatusFilterProvider = StateProvider<int?>((ref) => null);

/// 订单列表 Provider
final orderListProvider = FutureProvider<List<Order>>((ref) async {
  final status = ref.watch(orderStatusFilterProvider);
  final response = await HttpClient().dio.get(
    '/orders',
    queryParameters: status != null ? {'status': status} : null,
  );
  final List<dynamic> data = response.data;
  return data.map((json) => Order.fromJson(json)).toList();
});

/// 订单数量统计 Provider
/// 返回 Map<int, int>，key 是状态码，value 是数量
final orderCountsProvider = FutureProvider<Map<int, int>>((ref) async {
  final response = await HttpClient().dio.get('/orders/counts');
  final Map<String, dynamic> data = response.data;
  return data.map((key, value) => MapEntry(int.parse(key), value as int));
});

/// 订单管理 Provider
class OrderController {
  final Ref ref;

  OrderController(this.ref);

  /// 创建订单
  Future<Order> createOrder(int addressId) async {
    final response = await HttpClient().dio.post(
      '/orders',
      data: {'address_id': addressId},
    );
    // 刷新订单列表和统计
    ref.invalidate(orderListProvider);
    ref.invalidate(orderCountsProvider);
    return Order.fromJson(response.data);
  }

  /// 支付订单
  Future<void> payOrder(int orderId) async {
    await HttpClient().dio.post('/orders/$orderId/pay');
    ref.invalidate(orderListProvider);
    ref.invalidate(orderCountsProvider);
  }

  /// 确认收货
  Future<void> confirmReceipt(int orderId) async {
    await HttpClient().dio.put('/orders/$orderId/receipt');
    ref.invalidate(orderListProvider);
    ref.invalidate(orderCountsProvider);
  }

  /// 评价订单
  Future<void> reviewOrder(int orderId, String content, int rating) async {
    await HttpClient().dio.post(
      '/orders/$orderId/review',
      data: {'content': content, 'rating': rating},
    );
    ref.invalidate(orderListProvider);
    ref.invalidate(orderCountsProvider);
  }

  /// 申请售后 (Status 4 -> 5)
  Future<void> applyAfterSales(int orderId) async {
    await HttpClient().dio.post('/orders/$orderId/after-sales');
    ref.invalidate(orderListProvider);
    ref.invalidate(orderCountsProvider);
  }
}

final orderControllerProvider = Provider((ref) => OrderController(ref));
