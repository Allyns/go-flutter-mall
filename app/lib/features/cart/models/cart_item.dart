import 'package:go_flutter_mall/features/product/models/product.dart';

/// 购物车项模型
class CartItem {
  final int id;
  final int productId;
  final Product product;
  final int quantity;
  final bool selected;

  CartItem({
    required this.id,
    required this.productId,
    required this.product,
    required this.quantity,
    required this.selected,
  });

  factory CartItem.fromJson(Map<String, dynamic> json) {
    return CartItem(
      id: json['ID'],
      productId: json['product_id'],
      product: Product.fromJson(json['product']),
      quantity: json['quantity'],
      selected: json['selected'],
    );
  }
}
