import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';

// Notification Model
class AppNotification {
  final int id;
  final String title;
  final String content;
  final bool isRead;
  final String createdAt;

  AppNotification({
    required this.id,
    required this.title,
    required this.content,
    required this.isRead,
    required this.createdAt,
  });

  factory AppNotification.fromJson(Map<String, dynamic> json) {
    return AppNotification(
      id: json['id'],
      title: json['title'],
      content: json['content'],
      isRead: json['is_read'],
      createdAt: json['created_at'],
    );
  }
}

// Notification Notifier
class NotificationNotifier extends StateNotifier<AsyncValue<List<AppNotification>>> {
  NotificationNotifier() : super(const AsyncValue.loading());

  Future<void> fetchNotifications() async {
    state = const AsyncValue.loading();
    try {
      final response = await HttpClient().dio.get('/notifications');
      final List<dynamic> data = response.data;
      final notifications = data.map((json) => AppNotification.fromJson(json)).toList();
      state = AsyncValue.data(notifications);
    } catch (e, stack) {
      state = AsyncValue.error(e, stack);
    }
  }

  Future<void> markAsRead(int id) async {
    try {
      await HttpClient().dio.put('/notifications/$id/read');
      // Optimistic update
      state.whenData((notifications) {
        state = AsyncValue.data(
          notifications.map((n) => n.id == id ? AppNotification(
            id: n.id,
            title: n.title,
            content: n.content,
            isRead: true,
            createdAt: n.createdAt,
          ) : n).toList(),
        );
      });
    } catch (e) {
      // Ignore
    }
  }
}

final notificationProvider = StateNotifierProvider<NotificationNotifier, AsyncValue<List<AppNotification>>>((ref) {
  return NotificationNotifier();
});
