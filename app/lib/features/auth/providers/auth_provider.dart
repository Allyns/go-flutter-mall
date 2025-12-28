import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:dio/dio.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';
import 'package:go_flutter_mall/features/auth/models/user.dart';

/// 认证状态类
class AuthState {
  final User? user; // 当前登录用户，为空表示未登录
  final bool isLoading; // 是否正在加载
  final String? error; // 错误信息

  AuthState({this.user, this.isLoading = false, this.error});

  // 是否已认证
  bool get isAuthenticated => user != null;

  // 复制并修改状态
  AuthState copyWith({User? user, bool? isLoading, String? error}) {
    return AuthState(
      user: user ?? this.user,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// 认证状态管理器 (Notifier)
/// 负责处理登录、注册、注销逻辑并更新状态
class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier() : super(AuthState()) {
    _checkLoginStatus();
  }

  // 检查本地存储中是否有 Token，尝试恢复登录状态
  Future<void> _checkLoginStatus() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('token');
    
    if (token != null) {
      try {
        // 调用 /auth/me 获取最新用户信息
        final response = await HttpClient().dio.get('/auth/me');
        final userData = response.data;
        
        // 更新本地存储和状态
        await prefs.setString('user', jsonEncode(userData));
        state = AuthState(user: User.fromJson(userData));
      } catch (e) {
        // 如果获取失败（如 token 过期），尝试使用本地缓存
        print('Failed to refresh user profile: $e');
        
        final userJson = prefs.getString('user');
        if (userJson != null) {
          try {
            final user = User.fromJson(jsonDecode(userJson));
            state = AuthState(user: user);
          } catch (_) {}
        }
        
        // 如果连本地缓存也没有，或者 Token 失效严重，可能需要强制登出
        // 这里暂时保持现状，等待后续请求 401 触发登出
      }
    }
  }

  /// 登录方法
  Future<void> login(String account, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      // Determine if account is email or username
      final isEmail = RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(account);
      
      final data = {
        'password': password,
        if (isEmail) 'email': account else 'username': account,
      };

      final response = await HttpClient().dio.post(
        '/auth/login',
        data: data,
      );

      final token = response.data['token'];
      final userData = response.data['user'];

      // 保存 Token 和用户信息到本地
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('token', token);
      await prefs.setString('user', jsonEncode(userData));

      // 更新状态
      state = AuthState(user: User.fromJson(userData));
    } on DioException catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.response?.data['error'] ?? 'Login failed',
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 注册方法
  Future<void> register(String username, String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await HttpClient().dio.post(
        '/auth/register',
        data: {'username': username, 'email': email, 'password': password},
      );
      // 注册成功后自动登录
      await login(email, password);
    } on DioException catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.response?.data['error'] ?? 'Registration failed',
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 注销方法
  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('token');
    await prefs.remove('user');
    state = AuthState();
  }
}

/// 全局 AuthProvider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier();
});
