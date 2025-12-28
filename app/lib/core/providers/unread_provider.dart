import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';

/// 未读消息计数 Provider
class UnreadCountNotifier extends StateNotifier<int> {
  UnreadCountNotifier() : super(0);

  /// 从服务器获取未读数量
  Future<void> fetchUnreadCount() async {
    try {
      final response = await HttpClient().dio.get('/notifications/unread-count');
      final data = response.data;
      if (data != null && data['unread_count'] != null) {
        state = (data['unread_count'] as num).toInt();
      }
    } catch (e) {
      // 忽略错误
    }
  }

  /// 手动更新数量 (例如 WS 推送时 +1)
  void increment() {
    state++;
  }

  /// 手动减少 (例如已读)
  void decrement() {
    if (state > 0) state--;
  }
  
  /// 重置为0
  void clear() {
    state = 0;
  }
}

final unreadCountProvider = StateNotifierProvider<UnreadCountNotifier, int>((ref) {
  return UnreadCountNotifier();
});
