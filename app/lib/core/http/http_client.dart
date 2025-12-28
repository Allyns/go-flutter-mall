import 'dart:io';

import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 全局 HTTP 客户端封装
/// 负责统一配置请求 URL、超时时间以及拦截器
class HttpClient {
  static final HttpClient _instance = HttpClient._internal();
  late Dio _dio;

  // 工厂构造函数，实现单例模式
  factory HttpClient() {
    return _instance;
  }

  HttpClient._internal() {
    // 自动检测运行平台，设置正确的基础 URL
    // Android 模拟器: 10.0.2.2
    // iOS 模拟器 / macOS: localhost
    // 真机: 需要填写局域网 IP (例如 192.168.x.x)
    // 当前局域网 IP: 192.168.5.165
    // final String baseUrl = Platform.isAndroid ? 'http://10.0.2.2:8080/api' : 'http://localhost:8080/api';

    // 真机调试模式：直接使用局域网 IP
    const String baseUrl = 'http://192.168.5.165:8080/api';

    // 初始化 Dio 实例
    _dio = Dio(
      BaseOptions(
        // 后端 API 基础 URL
        baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 5),
        receiveTimeout: const Duration(seconds: 3),
      ),
    );

    // 1. 添加认证拦截器
    _dio.interceptors.add(
      InterceptorsWrapper(
        // 请求拦截器: 发送请求前执行
        onRequest: (options, handler) async {
          // 从本地存储获取 Token
          final prefs = await SharedPreferences.getInstance();
          final token = prefs.getString('token');

          // 如果 Token 存在，则添加到请求头 Authorization 中
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          return handler.next(options);
        },
      ),
    );

    // 2. 添加详细日志拦截器 (使用 Dio 内置的 LogInterceptor)
    _dio.interceptors.add(
      LogInterceptor(
        request: true, // 请求体
        requestHeader: true, // 请求头
        requestBody: true, // 请求数据
        responseHeader: true, // 响应头
        responseBody: true, // 响应数据
        error: true, // 错误信息
        logPrint: (object) {
          print('DIO LOG: $object');
        },
      ),
    );
  }

  /// 获取 Dio 实例
  Dio get dio => _dio;
}
