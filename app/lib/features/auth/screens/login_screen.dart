import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:go_flutter_mall/features/auth/providers/auth_provider.dart';

/// 登录屏幕
/// 使用 HooksConsumerWidget 以同时使用 Hooks 和 Riverpod
class LoginScreen extends HookConsumerWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 使用 Flutter Hooks 管理 TextController
    final emailController = useTextEditingController();
    final passwordController = useTextEditingController();
    
    // 监听 Auth 状态
    final authState = ref.watch(authProvider);

    // 显示错误提示
    useEffect(() {
      if (authState.error != null) {
        // 使用 Future.microtask 避免在构建过程中显示 SnackBar
        Future.microtask(() {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(authState.error!)),
            );
          }
        });
      }
      return null;
    }, [authState.error]);

    return Scaffold(
      appBar: AppBar(title: const Text('登录')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // 邮箱输入框
            TextField(
              controller: emailController,
              decoration: const InputDecoration(labelText: 'Email or Username'),
              keyboardType: TextInputType.emailAddress,
            ),
            const SizedBox(height: 16),
            // 密码输入框
            TextField(
              controller: passwordController,
              decoration: const InputDecoration(labelText: '密码'),
              obscureText: true,
            ),
            const SizedBox(height: 24),
            // 登录按钮
            if (authState.isLoading)
              const CircularProgressIndicator()
            else
              ElevatedButton(
                onPressed: () {
                  // 调用 AuthProvider 进行登录
                  ref.read(authProvider.notifier).login(
                        emailController.text,
                        passwordController.text,
                      );
                },
                child: const Text('登录'),
              ),
            TextButton(
              onPressed: () => context.push('/register'),
              child: const Text('创建新账号'),
            ),
          ],
        ),
      ),
    );
  }
}
