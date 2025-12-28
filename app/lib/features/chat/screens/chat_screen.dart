import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_flutter_mall/features/chat/services/chat_service.dart';
import 'package:go_flutter_mall/features/auth/providers/auth_provider.dart';
import 'package:go_flutter_mall/core/providers/unread_provider.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';

class ChatScreen extends HookConsumerWidget {
  const ChatScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final messages = ref.watch(chatProvider);
    final user = ref.watch(authProvider).user;
    final textController = useTextEditingController();
    final scrollController = useScrollController();

    // 自动连接 WS 并标记已读
    useEffect(() {
      final notifier = ref.read(chatProvider.notifier);
      notifier.connect();
      
      // 进入聊天页面时，标记所有消息为已读
      Future.microtask(() async {
        try {
          await HttpClient().dio.put('/chat/read');
          ref.read(unreadCountProvider.notifier).fetchUnreadCount(); // 更新未读数
        } catch (e) {
          debugPrint('Failed to mark messages as read: $e');
        }
      });

      return () => notifier.disconnect();
    }, []);

    // 自动滚动到底部
    useEffect(() {
      if (scrollController.hasClients) {
        Future.delayed(const Duration(milliseconds: 100), () {
          scrollController.animateTo(
            scrollController.position.maxScrollExtent,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
          );
        });
      }
      return null;
    }, [messages.length]);

    void handleSend() {
      final content = textController.text.trim();
      if (content.isNotEmpty) {
        ref.read(chatProvider.notifier).sendMessage(content);
        textController.clear();
      }
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('联系客服'),
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              controller: scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: messages.length,
              itemBuilder: (context, index) {
                final msg = messages[index];
                final isMe = msg.senderType == 'user'; // 当前用户发送的消息

                return Align(
                  alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.symmetric(vertical: 4),
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                    decoration: BoxDecoration(
                      color: isMe ? Colors.green : const Color(0xFF1E1E1E),
                      borderRadius: BorderRadius.circular(20).copyWith(
                        topLeft: isMe ? const Radius.circular(20) : const Radius.circular(0),
                        topRight: isMe ? const Radius.circular(0) : const Radius.circular(20),
                      ),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.05),
                          blurRadius: 5,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: Text(
                      msg.content,
                      style: const TextStyle(
                        color: Colors.white,
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: const BoxDecoration(
              color: Color(0xFF1E1E1E),
              border: Border(top: BorderSide(color: Color(0xFF333333))),
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: textController,
                    style: const TextStyle(color: Colors.white),
                    decoration: InputDecoration(
                      hintText: '请输入消息...',
                      hintStyle: TextStyle(color: Colors.grey[400]),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(30),
                        borderSide: BorderSide.none,
                      ),
                      filled: true,
                      fillColor: const Color(0xFF2C2C2C),
                    ),
                    onSubmitted: (_) => handleSend(),
                  ),
                ),
                const SizedBox(width: 12),
                FloatingActionButton(
                  onPressed: handleSend,
                  mini: true,
                  backgroundColor: Colors.green,
                  child: const Icon(Icons.send, color: Colors.white),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
