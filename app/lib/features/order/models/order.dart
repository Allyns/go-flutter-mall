/// 订单模型
class Order {
  final int id;
  final String orderNo;
  final double totalAmount;
  final int status; // 0: Pending, 1: Paid, etc.
  final String createdAt;
  final List<OrderItem> items;

  Order({
    required this.id,
    required this.orderNo,
    required this.totalAmount,
    required this.status,
    required this.createdAt,
    required this.items,
  });

  factory Order.fromJson(Map<String, dynamic> json) {
    return Order(
      id: json['id'],
      orderNo: json['order_no'],
      totalAmount: (json['total_amount'] as num).toDouble(),
      status: json['status'],
      createdAt: json['created_at'],
      items: (json['items'] as List)
          .map((item) => OrderItem.fromJson(item))
          .toList(),
    );
  }

  // 获取状态描述
  String get statusText {
    switch (status) {
      case 0:
        return '待付款';
      case 1:
        return '待发货';
      case 2:
        return '待收货';
      case 3:
        return '待评价';
      case 4:
        return '已完成';
      case 5:
        return '售后中';
      case -1:
        return '已取消';
      default:
        return '未知状态';
    }
  }
}

/// 订单项模型
class OrderItem {
  final int id;
  final String productName;
  final double price;
  final int quantity;
  final String productImage;

  OrderItem({
    required this.id,
    required this.productName,
    required this.price,
    required this.quantity,
    required this.productImage,
  });

  factory OrderItem.fromJson(Map<String, dynamic> json) {
    return OrderItem(
      id: json['id'],
      productName: json['product_name'],
      price: (json['price'] as num).toDouble(),
      quantity: json['quantity'],
      productImage: json['product_image'],
    );
  }
}
