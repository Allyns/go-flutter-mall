import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/features/auth/providers/auth_provider.dart';
import 'package:go_flutter_mall/features/auth/screens/login_screen.dart';
import 'package:go_flutter_mall/features/auth/screens/register_screen.dart';
import 'package:go_flutter_mall/features/home/screens/home_screen.dart';
import 'package:go_flutter_mall/features/product/screens/product_detail_screen.dart';
import 'package:go_flutter_mall/features/order/screens/checkout_screen.dart';
import 'package:go_flutter_mall/features/address/screens/address_list_screen.dart';
import 'package:go_flutter_mall/features/address/screens/address_edit_screen.dart';
import 'package:go_flutter_mall/features/chat/screens/chat_screen.dart';

/// 路由配置 Provider
/// 使用 Riverpod 监听 authProvider 的状态变化，以实现重定向
final routerProvider = Provider<GoRouter>((ref) {
  // 监听认证状态
  final authState = ref.watch(authProvider);

  return GoRouter(
    // 初始路由
    initialLocation: '/',
    // 调试日志
    debugLogDiagnostics: true,
    // 定义路由表
    routes: [
      // 首页
      GoRoute(path: '/', builder: (context, state) => const HomeScreen()),
      // 登录页
      GoRoute(path: '/login', builder: (context, state) => const LoginScreen()),
      // 注册页
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterScreen(),
      ),
      // 商品详情页，带参数 id
      GoRoute(
        path: '/product/:id',
        builder: (context, state) {
          final id = state.pathParameters['id']!;
          return ProductDetailScreen(productId: int.parse(id));
        },
      ),
      // 结算页
      GoRoute(
        path: '/checkout',
        builder: (context, state) => const CheckoutScreen(),
      ),
      // 客服聊天页
      GoRoute(
        path: '/chat',
        builder: (context, state) => const ChatScreen(),
      ),
      // 地址列表页
      GoRoute(
        path: '/addresses',
        builder: (context, state) => const AddressListScreen(),
        routes: [
          // 添加地址
          GoRoute(
            path: 'add',
            builder: (context, state) => const AddressEditScreen(),
          ),
          // 编辑地址
          GoRoute(
            path: 'edit',
            builder: (context, state) {
              final address = state.extra as Address;
              return AddressEditScreen(address: address);
            },
          ),
        ],
      ),
    ],
    // 重定向逻辑
    redirect: (context, state) {
      final isLoggedIn = authState.isAuthenticated;
      final isLoggingIn = state.uri.path == '/login';
      final isRegistering = state.uri.path == '/register';

      // 如果未登录，且不是去登录或注册页，则重定向到登录页
      if (!isLoggedIn && !isLoggingIn && !isRegistering) {
        return '/login';
      }

      // 如果已登录，且正在访问登录或注册页，则重定向到首页
      if (isLoggedIn && (isLoggingIn || isRegistering)) {
        return '/';
      }

      // 否则不进行重定向
      return null;
    },
  );
});
