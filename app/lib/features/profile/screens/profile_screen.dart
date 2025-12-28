import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_flutter_mall/features/auth/providers/auth_provider.dart';
import 'package:go_flutter_mall/features/order/providers/order_provider.dart';
import 'package:go_flutter_mall/features/home/providers/home_provider.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter/services.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authProvider).user;
    
    // 强制刷新订单数据以获取最新状态
    ref.watch(orderListProvider);
    // 获取订单数量统计
    final orderCountsAsync = ref.watch(orderCountsProvider);
    final orderCounts = orderCountsAsync.value ?? {};

    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            // 头部区域 (黑色 -> 绿色渐变)
            Container(
              padding: const EdgeInsets.only(top: 60, bottom: 20, left: 20, right: 20),
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [Colors.black, Color(0xFF1B5E20)], // Black -> Dark Green
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
              ),
              child: Row(
                children: [
                  CircleAvatar(
                    radius: 30,
                    backgroundColor: Colors.white,
                    backgroundImage: (user?.avatar != null && user!.avatar!.isNotEmpty)
                        ? NetworkImage(user!.avatar!)
                        : null,
                    child: (user?.avatar == null || user!.avatar!.isEmpty)
                        ? const Icon(Icons.person, size: 40, color: Colors.green)
                        : null,
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          user?.username ?? '登录 / 注册',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 4),
                        if (user != null)
                          GestureDetector(
                            onTap: () {
                              Clipboard.setData(ClipboardData(text: user.id.toString()));
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(content: Text('ID 已复制到剪贴板')),
                              );
                            },
                            child: Row(
                              children: [
                                Text(
                                  'ID: ${user.id}',
                                  style: const TextStyle(color: Colors.white70),
                                ),
                                const SizedBox(width: 4),
                                const Icon(Icons.copy, size: 14, color: Colors.white70),
                              ],
                            ),
                          )
                        else
                          const Text(
                            '欢迎来到 Go 商城',
                            style: TextStyle(color: Colors.white70),
                          ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            
            // 订单状态栏
            Padding(
              padding: const EdgeInsets.all(12.0),
              child: Card(
                elevation: 0,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                color: const Color(0xFF1E1E1E),
                child: Column(
                  children: [
                    // Header
                    ListTile(
                      title: const Text('我的订单', style: TextStyle(fontWeight: FontWeight.bold, color: Colors.white)),
                      trailing: const Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text('全部订单', style: TextStyle(fontSize: 12, color: Colors.grey)),
                          Icon(Icons.arrow_forward_ios, size: 12, color: Colors.grey),
                        ],
                      ),
                      onTap: () {
                        ref.read(orderStatusFilterProvider.notifier).state = null; // null for All
                        ref.read(homeTabIndexProvider.notifier).state = 2; // Order Tab
                      },
                    ),
                    const Divider(height: 1),
                    // Status Grid
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: [
                          _buildOrderStatusItem(context, ref, Icons.payment, '待付款', 0, orderCounts[0] ?? 0),
                          _buildOrderStatusItem(context, ref, Icons.local_shipping, '待发货', 1, orderCounts[1] ?? 0),
                          _buildOrderStatusItem(context, ref, Icons.assignment_turned_in, '待收货', 2, orderCounts[2] ?? 0),
                          _buildOrderStatusItem(context, ref, Icons.rate_review, '待评价', 3, orderCounts[3] ?? 0),
                          _buildOrderStatusItem(context, ref, Icons.assignment_return, '退换/售后', 5, orderCounts[5] ?? 0), // Changed to 5
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // 功能菜单
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12.0),
              child: Card(
                elevation: 0,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                color: const Color(0xFF1E1E1E),
                child: Column(
                  children: [
                    _buildMenuItem(Icons.location_on_outlined, '收货地址', () {
                       context.push('/addresses');
                    }),
                    const Divider(height: 1, indent: 50),
                    _buildMenuItem(Icons.favorite_border, '我的收藏', () {
                       ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('收藏功能待开发')));
                    }),
                    const Divider(height: 1, indent: 50),
                    _buildMenuItem(Icons.headset_mic_outlined, '客户服务', () {
                       context.push('/chat');
                    }),
                    const Divider(height: 1, indent: 50),
                    _buildMenuItem(Icons.settings_outlined, '设置', () {
                       // 简单的设置菜单，暂时只包含退出登录
                       showModalBottomSheet(
                         context: context,
                         builder: (context) => Column(
                           mainAxisSize: MainAxisSize.min,
                           children: [
                             ListTile(
                               leading: const Icon(Icons.logout, color: Colors.red),
                               title: const Text('退出登录', style: TextStyle(color: Colors.red)),
                               onTap: () {
                                 Navigator.pop(context); // 关闭 BottomSheet
                                 ref.read(authProvider.notifier).logout();
                               },
                             ),
                             const SizedBox(height: 20),
                           ],
                         ),
                       );
                    }),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOrderStatusItem(BuildContext context, WidgetRef ref, IconData icon, String label, int? status, [int count = 0]) {
    return GestureDetector(
      onTap: () {
        ref.read(orderStatusFilterProvider.notifier).state = status;
        ref.read(homeTabIndexProvider.notifier).state = 2; // Order List Tab Index
      },
      child: Column(
        children: [
          Stack(
            clipBehavior: Clip.none,
            children: [
              Icon(icon, size: 28, color: Colors.grey[700]),
              if (count > 0)
                Positioned(
                  right: -6,
                  top: -6,
                  child: Container(
                    padding: const EdgeInsets.all(4),
                    decoration: const BoxDecoration(
                      color: Colors.red,
                      shape: BoxShape.circle,
                    ),
                    constraints: const BoxConstraints(
                      minWidth: 16,
                      minHeight: 16,
                    ),
                    child: Text(
                      count > 99 ? '99+' : count.toString(),
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 8),
          Text(label, style: const TextStyle(fontSize: 12, color: Colors.white70)),
        ],
      ),
    );
  }

  Widget _buildMenuItem(IconData icon, String title, VoidCallback? onTap, {Widget? trailing}) {
    return ListTile(
      leading: Icon(icon, color: Colors.green), // 绿色图标
      title: Text(title, style: const TextStyle(fontSize: 15, color: Colors.white)),
      trailing: trailing ?? const Icon(Icons.arrow_forward_ios, size: 16, color: Colors.grey),
      onTap: onTap,
    );
  }
}
