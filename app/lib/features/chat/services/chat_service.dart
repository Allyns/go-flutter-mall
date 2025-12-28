import 'dart:convert';
import 'dart:io';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:go_flutter_mall/features/auth/providers/auth_provider.dart';

// 聊天消息模型
class ChatMessage {
  final int senderId;
  final String senderType;
  final String content;
  final String type; // text, image

  ChatMessage({
    required this.senderId,
    required this.senderType,
    required this.content,
    this.type = 'text',
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      senderId: json['sender_id'],
      senderType: json['sender_type'],
      content: json['content'],
      type: json['type'] ?? 'text',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'sender_id': senderId,
      'sender_type': senderType,
      'content': content,
      'type': type,
    };
  }
}

// 聊天状态管理
class ChatNotifier extends StateNotifier<List<ChatMessage>> {
  final Ref ref;
  WebSocketChannel? _channel;

  ChatNotifier(this.ref) : super([]);

  void connect() {
    final user = ref.read(authProvider).user;
    if (user == null) return;

    // 防止重复连接
    if (_channel != null) return;

    // TODO: 生产环境应从配置读取 host, 暂时硬编码 localhost
    // Android 模拟器需用 10.0.2.2，iOS 模拟器用 localhost
    // 真机需用局域网 IP (192.168.5.165)
    
    String host;
    if (Platform.isAndroid) {
      // 如果是 Android 模拟器，使用 10.0.2.2
      // 如果是真机，使用局域网 IP
      // 这里无法自动判断是否为模拟器，通常建议开发时配置化
      // 简单起见，这里假设您正在使用真机调试，使用局域网 IP
      // 如果切回模拟器，请改回 '10.0.2.2:8080'
      host = '192.168.5.165:8080'; 
    } else {
      host = 'localhost:8080';
    }
    
    final wsUrl = Uri.parse('ws://$host/api/ws?user_id=${user.id}&type=user');
    
    try {
      _channel = WebSocketChannel.connect(wsUrl);
      print('WS Connected to $wsUrl');

      _channel!.stream.listen(
        (message) {
          try {
            final data = jsonDecode(message);
            if (data['type'] == 'message') {
               final msg = ChatMessage.fromJson(data['payload']);
               
               // 如果是当前用户发送的消息，且在发送时已经乐观添加到列表中，则忽略
               // 实际上后端发回来的消息 ID 可能还没有生成，或者生成了但不一致
               // 简单做法：如果是 'user' 类型且 senderId == 当前用户 ID，则认为是自己发的
               final currentUser = ref.read(authProvider).user;
               if (currentUser != null && msg.senderType == 'user' && msg.senderId == currentUser.id) {
                 return;
               }
               
               // 添加到消息列表
               state = [...state, msg];
            }
          } catch (e) {
            print('Parse Error: $e');
          }
        },
        onError: (error) {
          print('WS Error: $error');
          _channel = null;
        },
        onDone: () {
          print('WS Closed');
          _channel = null;
        },
      );
    } catch (e) {
      print('Connection Error: $e');
    }
  }

  void sendMessage(String content) {
    if (_channel == null) {
      print('WS not connected, trying to connect...');
      connect();
      // 简单重试逻辑：延迟一下再发（实际应有更好的队列机制）
      // 这里如果连接未建立，本次发送可能会失败
      if (_channel == null) return;
    }
    
    final user = ref.read(authProvider).user;
    if (user == null) return;

    final msg = ChatMessage(
      senderId: user.id,
      senderType: 'user',
      content: content,
    );

    // 乐观更新 UI
    state = [...state, msg];

    try {
      _channel!.sink.add(jsonEncode({
        'type': 'message',
        'payload': msg.toJson(),
      }));
    } catch (e) {
      print('Send Error: $e');
      // 发送失败可能需要从 state 移除或标记失败
    }
  }

  void disconnect() {
    _channel?.sink.close();
    _channel = null;
  }
}

final chatProvider = StateNotifierProvider<ChatNotifier, List<ChatMessage>>((ref) {
  return ChatNotifier(ref);
});
