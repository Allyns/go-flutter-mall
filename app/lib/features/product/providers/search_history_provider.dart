import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';

// Search History Model
class SearchHistory {
  final String keyword;

  SearchHistory({required this.keyword});

  factory SearchHistory.fromJson(Map<String, dynamic> json) {
    return SearchHistory(keyword: json['keyword']);
  }
}

// Search History Notifier
class SearchHistoryNotifier extends StateNotifier<List<SearchHistory>> {
  SearchHistoryNotifier() : super([]) {
    fetchHistory();
  }

  Future<void> fetchHistory() async {
    try {
      final response = await HttpClient().dio.get('/search/history');
      final List<dynamic> data = response.data;
      state = data.map((json) => SearchHistory.fromJson(json)).toList();
    } catch (e) {
      // Ignore errors for now
    }
  }

  Future<void> addHistory(String keyword) async {
    if (keyword.trim().isEmpty) return;
    try {
      // Optimistic update
      // Remove existing if present to move to top
      state = [
        SearchHistory(keyword: keyword),
        ...state.where((h) => h.keyword != keyword),
      ];
      
      await HttpClient().dio.post('/search/history', data: {'keyword': keyword});
    } catch (e) {
       // Revert or ignore
    }
  }

  Future<void> clearHistory() async {
    try {
      state = [];
      await HttpClient().dio.delete('/search/history');
    } catch (e) {
      // Ignore
    }
  }
}

final searchHistoryProvider = StateNotifierProvider<SearchHistoryNotifier, List<SearchHistory>>((ref) {
  return SearchHistoryNotifier();
});
