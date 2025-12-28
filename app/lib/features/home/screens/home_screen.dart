import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_flutter_mall/features/product/screens/product_list_screen.dart';
import 'package:go_flutter_mall/features/cart/screens/cart_screen.dart';
import 'package:go_flutter_mall/features/order/screens/order_list_screen.dart';
import 'package:go_flutter_mall/features/profile/screens/profile_screen.dart';
import 'package:go_flutter_mall/features/notification/screens/notification_screen.dart';
import 'package:go_flutter_mall/features/notification/providers/notification_provider.dart';
import 'package:go_flutter_mall/features/home/providers/home_provider.dart';
import 'package:go_flutter_mall/core/providers/unread_provider.dart';

class HomeScreen extends HookConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final currentIndex = ref.watch(homeTabIndexProvider);
    final unreadCount = ref.watch(unreadCountProvider);

    // 初始加载未读数
    useEffect(() {
      Future.microtask(() => ref.read(unreadCountProvider.notifier).fetchUnreadCount());
      return null;
    }, []);

    final pages = [
      const ProductListScreen(),
      const CartScreen(),
      const OrderListScreen(),
      const NotificationScreen(), // 消息页面
      const ProfileScreen(),
    ];

    return Scaffold(
      body: pages[currentIndex],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: currentIndex,
        onTap: (index) {
          ref.read(homeTabIndexProvider.notifier).state = index;
          // 如果点击的是消息 Tab (索引为 3)，则刷新消息列表并清除红点
          if (index == 3) {
            ref.read(notificationProvider.notifier).fetchNotifications();
            ref.read(unreadCountProvider.notifier).fetchUnreadCount(); // Re-fetch to be sure, or just clear locally if we assume all read
          }
        },
        type: BottomNavigationBarType.fixed,
        // selectedItemColor: Colors.black, // Removed to use Theme
        // unselectedItemColor: Colors.grey, // Removed to use Theme
        selectedFontSize: 12,
        unselectedFontSize: 12,
        items: [
          const BottomNavigationBarItem(icon: Icon(Icons.home_outlined), activeIcon: Icon(Icons.home), label: '首页'),
          const BottomNavigationBarItem(icon: Icon(Icons.shopping_cart_outlined), activeIcon: Icon(Icons.shopping_cart), label: '购物车'),
          const BottomNavigationBarItem(icon: Icon(Icons.assignment_outlined), activeIcon: Icon(Icons.assignment), label: '订单'),
          BottomNavigationBarItem(
            icon: Badge(
              isLabelVisible: unreadCount > 0,
              label: Text('$unreadCount'),
              child: const Icon(Icons.notifications_outlined),
            ),
            activeIcon: Badge(
              isLabelVisible: unreadCount > 0,
              label: Text('$unreadCount'),
              child: const Icon(Icons.notifications),
            ),
            label: '消息',
          ),
          const BottomNavigationBarItem(icon: Icon(Icons.person_outline), activeIcon: Icon(Icons.person), label: '我的'),
        ],
      ),
    );
  }
}
