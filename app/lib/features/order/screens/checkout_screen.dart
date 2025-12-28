import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';
import 'package:go_flutter_mall/features/cart/providers/cart_provider.dart';
import 'package:go_flutter_mall/features/order/providers/order_provider.dart';

// 简单的地址模型
class SimpleAddress {
  final int id;
  final String receiverName;
  final String receiverPhone;
  final String detailAddress;

  SimpleAddress({
    required this.id,
    required this.receiverName,
    required this.receiverPhone,
    required this.detailAddress,
  });

  factory SimpleAddress.fromJson(Map<String, dynamic> json) {
    return SimpleAddress(
      id: json['ID'] ?? 0,
      receiverName: json['receiver_name'] ?? '',
      receiverPhone: json['receiver_phone'] ?? '',
      detailAddress: "${json['province']} ${json['city']} ${json['district']} ${json['detail_address']}",
    );
  }
}

// 默认地址 Provider
final defaultAddressProvider = FutureProvider.autoDispose<SimpleAddress?>((ref) async {
  try {
    // 这里简单实现：获取地址列表并取第一个作为默认地址
    // 实际项目中应该有一个 /api/addresses/default 接口
    final response = await HttpClient().dio.get('/addresses'); // 需要后端实现此接口，或者复用现有逻辑
    // 暂时 mock 一下，或者让用户必须先去设置地址
    // 由于后端还没做完整的地址管理，这里先尝试获取列表
    final List<dynamic> data = response.data;
    if (data.isNotEmpty) {
      return SimpleAddress.fromJson(data.first);
    }
    return null;
  } catch (e) {
    return null;
  }
});

/// 结算/确认订单屏幕
class CheckoutScreen extends ConsumerWidget {
  const CheckoutScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 获取购物车总价
    final totalAmount = ref.watch(cartTotalProvider);
    // 获取默认地址
    final addressAsyncValue = ref.watch(defaultAddressProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('确认订单')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 地址选择部分
            addressAsyncValue.when(
              data: (address) {
                if (address == null) {
                  return Card(
                    child: ListTile(
                      leading: const Icon(Icons.location_off, color: Colors.red),
                      title: const Text('暂无收货地址'),
                      subtitle: const Text('请先添加收货地址 (点击添加)'),
                      trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                      onTap: () async {
                        await context.push('/addresses');
                        // 从地址页返回后刷新默认地址
                        ref.invalidate(defaultAddressProvider);
                      },
                    ),
                  );
                }
                return Card(
                  child: ListTile(
                    leading: const Icon(Icons.location_on, color: Colors.green),
                    title: Text('${address.receiverName}, ${address.receiverPhone}'),
                    subtitle: Text(address.detailAddress),
                    trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                    onTap: () async {
                      await context.push('/addresses');
                      // 从地址页返回后刷新默认地址
                      ref.invalidate(defaultAddressProvider);
                    },
                  ),
                );
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (err, stack) => const Card(child: ListTile(title: Text('加载地址失败'))),
            ),
            
            const SizedBox(height: 24),
            // 金额概览
            Text(
              '订单详情',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('订单总额'),
                Text(
                  '¥${totalAmount.toStringAsFixed(2)}',
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 18, color: Color(0xFFFF5000)),
                ),
              ],
            ),
            const Spacer(),
            // 提交订单按钮
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () async {
                  final address = addressAsyncValue.value;
                  if (address == null) {
                     ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('请选择收货地址')));
                     return;
                  }

                  try {
                    // 使用真实的地址 ID 创建订单
                    await ref.read(orderControllerProvider).createOrder(address.id);
                    // 刷新购物车 (因为已购买项被清空)
                    ref.invalidate(cartProvider);

                    if (context.mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('订单创建成功'),
                        ),
                      );
                      // 跳转回首页或订单列表
                      context.go('/');
                    }
                  } catch (e) {
                    if (context.mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('创建订单失败: $e')),
                      );
                    }
                  }
                },
                child: const Text('提交订单'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
