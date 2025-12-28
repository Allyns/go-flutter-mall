/// 用户模型
/// 对应后端返回的用户数据结构
class User {
  final int id;
  final String username;
  final String email;
  final String? avatar;

  User({
    required this.id,
    required this.username,
    required this.email,
    this.avatar,
  });

  /// 从 JSON 构造 User 对象
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      username: json['username'],
      email: json['email'],
      avatar: json['avatar'],
    );
  }

  /// 将 User 对象转换为 JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'username': username,
      'email': email,
      'avatar': avatar,
    };
  }
}
