import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_flutter_mall/features/order/providers/order_provider.dart';
import 'package:go_flutter_mall/features/order/models/order.dart';

/// 订单列表屏幕
/// 展示用户的历史订单
class OrderListScreen extends HookConsumerWidget {
  const OrderListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final filterStatus = ref.watch(orderStatusFilterProvider);

    // 映射状态到 Tab 索引
    // null -> 0 (全部)
    // 0 -> 1 (待付款)
    // 1 -> 2 (待发货)
    // 2 -> 3 (待收货)
    // 3 -> 4 (待评价)
    // 5 -> 5 (售后)
    int getIndexFromStatus(int? status) {
      if (status == null) return 0;
      if (status <= 3) return status + 1;
      if (status == 5) return 5;
      return 0; // Default
    }

    // 映射 Tab 索引到状态
    int? getStatusFromIndex(int index) {
      if (index == 0) return null;
      if (index <= 4) return index - 1;
      if (index == 5) return 5;
      return null;
    }

    final tabController = useTabController(
      initialLength: 6,
      initialIndex: getIndexFromStatus(filterStatus),
    );

    // 监听 Tab 切换，更新 Provider
    useEffect(() {
      void listener() {
        if (!tabController.indexIsChanging) {
          // 只在动画结束或直接点击时更新
          final newStatus = getStatusFromIndex(tabController.index);
          if (ref.read(orderStatusFilterProvider) != newStatus) {
            Future.microtask(
              () => ref.read(orderStatusFilterProvider.notifier).state =
                  newStatus,
            );
          }
        }
      }

      tabController.addListener(listener);
      return () => tabController.removeListener(listener);
    }, [tabController]);

    final ordersAsync = ref.watch(orderListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('我的订单'),
        bottom: TabBar(
          controller: tabController,
          isScrollable: true,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.grey,
          indicatorColor: Colors.green,
          tabs: const [
            Tab(text: '全部'),
            Tab(text: '待付款'),
            Tab(text: '待发货'),
            Tab(text: '待收货'),
            Tab(text: '待评价'),
            Tab(text: '售后'),
          ],
        ),
      ),
      body: ordersAsync.when(
        data: (orders) {
          if (orders.isEmpty) {
            return const Center(child: Text('暂无订单'));
          }
          return ListView.builder(
            itemCount: orders.length,
            padding: const EdgeInsets.only(bottom: 20),
            itemBuilder: (context, index) {
              final order = orders[index];
              return _OrderCard(order: order);
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) {
          debugPrint("加载失败:$err");
          return Center(child: Text('加载失败: $err'));
        },
      ),
    );
  }
}

class _OrderCard extends HookConsumerWidget {
  final Order order;

  const _OrderCard({required this.order});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 倒计时逻辑
    final remainingTime = useState<Duration>(Duration.zero);

    useEffect(() {
      if (order.status != 0) return null;

      Timer? timer;
      void calculateTime() {
        try {
          final created = DateTime.parse(order.createdAt);
          final expireTime = created.add(const Duration(minutes: 30));
          final now = DateTime.now();
          final diff = expireTime.difference(now);

          if (diff.isNegative) {
            remainingTime.value = Duration.zero;
            // 如果倒计时结束但状态仍为 0，可能需要刷新列表
            // ref.invalidate(orderListProvider); // Optional: Auto refresh
          } else {
            remainingTime.value = diff;
          }
        } catch (e) {
          // Parse error
        }
      }

      calculateTime(); // Initial
      timer = Timer.periodic(
        const Duration(seconds: 1),
        (_) => calculateTime(),
      );

      return () => timer?.cancel();
    }, [order]);

    String formatDuration(Duration d) {
      if (d.inSeconds <= 0) return "已超时";
      final minutes = d.inMinutes.remainder(60).toString().padLeft(2, '0');
      final seconds = d.inSeconds.remainder(60).toString().padLeft(2, '0');
      return "$minutes:$seconds";
    }

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            // Header
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '订单号: ${order.orderNo}',
                      style: const TextStyle(fontSize: 12, color: Colors.grey),
                    ),
                    if (order.status == 0 && remainingTime.value.inSeconds > 0)
                      Text(
                        '支付剩余: ${formatDuration(remainingTime.value)}',
                        style: const TextStyle(fontSize: 12, color: Colors.red),
                      ),
                  ],
                ),
                Text(
                  order.statusText,
                  style: TextStyle(
                    color: order.status == 0
                        ? Colors.green
                        : Colors.white,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const Divider(),
            // Items
            ...order.items.map(
              (item) => Padding(
                padding: const EdgeInsets.symmetric(vertical: 8),
                child: Row(
                  children: [
                    ClipRRect(
                      borderRadius: BorderRadius.circular(4),
                      child: Image.network(
                        item.productImage,
                        width: 60,
                        height: 60,
                        fit: BoxFit.cover,
                        errorBuilder: (_, __, ___) => Container(
                          width: 60,
                          height: 60,
                          color: Colors.grey[200],
                        ),
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            item.productName,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 4),
                          Text(
                            '¥${item.price}',
                            style: const TextStyle(color: Colors.grey),
                          ),
                        ],
                      ),
                    ),
                    Text('x${item.quantity}'),
                  ],
                ),
              ),
            ),
            const Divider(),
            // Footer & Actions
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '总价: ¥${order.totalAmount}',
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                Row(
                  children: [
                    if (order.status == 0)
                      ElevatedButton(
                        onPressed: () => ref
                            .read(orderControllerProvider)
                            .payOrder(order.id),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFFFF5000),
                          foregroundColor: Colors.white,
                        ),
                        child: const Text('立即支付'),
                      ),
                    if (order.status == 2)
                      ElevatedButton(
                        onPressed: () =>
                            _showConfirmDialog(context, ref, order.id),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFFFF5000),
                          foregroundColor: Colors.white,
                        ),
                        child: const Text('确认收货'),
                      ),
                    if (order.status == 3)
                      ElevatedButton(
                        onPressed: () =>
                            _showReviewDialog(context, ref, order.id),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF1E1E1E),
                          foregroundColor: Colors.white,
                          side: const BorderSide(color: Colors.white),
                        ),
                        child: const Text('评价'),
                      ),
                    if (order.status == 4)
                      TextButton(
                        onPressed: () =>
                            _applyAfterSales(context, ref, order.id),
                        child: const Text(
                          '申请售后',
                          style: TextStyle(color: Colors.red),
                        ),
                      ),
                  ],
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _showConfirmDialog(BuildContext context, WidgetRef ref, int orderId) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认收货'),
        content: const Text('确认已收到商品？确认后将不能退款。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              ref.read(orderControllerProvider).confirmReceipt(orderId);
            },
            child: const Text('确认'),
          ),
        ],
      ),
    );
  }

  void _showReviewDialog(BuildContext context, WidgetRef ref, int orderId) {
    final contentController = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('评价商品'),
        content: TextField(
          controller: contentController,
          decoration: const InputDecoration(hintText: '请输入评价内容...'),
          maxLines: 3,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              ref
                  .read(orderControllerProvider)
                  .reviewOrder(orderId, contentController.text, 5);
            },
            child: const Text('提交评价'),
          ),
        ],
      ),
    );
  }

  // 申请售后 (Status 4 -> 5)
  Future<void> _applyAfterSales(
    BuildContext context,
    WidgetRef ref,
    int orderId,
  ) async {
    try {
      await ref
          .read(orderControllerProvider)
          .applyAfterSales(orderId); // Corrected Provider
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('售后申请已提交')));
    } catch (e) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('申请失败: $e')));
    }
  }

  // 辅助方法：根据状态码返回状态文本
  String _getStatusText(int status) {
    switch (status) {
      case 0:
        return '待付款';
      case 1:
        return '待发货';
      case 2:
        return '待收货';
      case 3:
        return '待评价';
      case 4:
        return '已完成';
      case 5:
        return '售后中';
      case -1:
        return '已取消';
      default:
        return '未知状态';
    }
  }
}
