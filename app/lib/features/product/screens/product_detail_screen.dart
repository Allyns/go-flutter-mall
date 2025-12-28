import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:go_flutter_mall/features/product/providers/product_provider.dart';
import 'package:go_flutter_mall/features/cart/providers/cart_provider.dart';
import 'package:go_flutter_mall/features/product/models/product.dart';

class ProductDetailScreen extends ConsumerStatefulWidget {
  final int productId;

  const ProductDetailScreen({super.key, required this.productId});

  @override
  ConsumerState<ProductDetailScreen> createState() =>
      _ProductDetailScreenState();
}

class _ProductDetailScreenState extends ConsumerState<ProductDetailScreen> {
  ProductSKU? _selectedSku;
  int _quantity = 1;

  @override
  Widget build(BuildContext context) {
    final productAsyncValue = ref.watch(
      productDetailProvider(widget.productId),
    );
    final reviewsAsyncValue = ref.watch(
      productReviewsProvider(widget.productId),
    );

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        forceMaterialTransparency: true,
        backgroundColor: Colors.transparent,
        leading: Container(
          margin: const EdgeInsets.all(8),
          decoration: const BoxDecoration(
            color: Colors.black26,
            shape: BoxShape.circle,
          ),
          child: IconButton(
            icon: const Icon(Icons.arrow_back, color: Colors.white),
            onPressed: () => Navigator.of(context).pop(),
          ),
        ),
        actions: [
          Container(
            margin: const EdgeInsets.all(8),
            decoration: const BoxDecoration(
              color: Colors.black26,
              shape: BoxShape.circle,
            ),
            child: IconButton(
              icon: const Icon(Icons.share, color: Colors.white),
              onPressed: () {},
            ),
          ),
        ],
      ),
      body: productAsyncValue.when(
        data: (product) {
          // 如果没有选中 SKU 且有 SKU 列表，默认选中第一个 (或者不选中强制用户选)
          // 这里逻辑：如果有 SKU，必须选一个才能加购。
          // 初始状态下 _selectedSku 为 null

          return Column(
            children: [
              Expanded(
                child: SingleChildScrollView(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // 商品大图
                      Image.network(
                        product.coverImage,
                        height: 375,
                        width: double.infinity,
                        fit: BoxFit.cover,
                      ),

                      Container(
                        color: const Color(0xFF1E1E1E),
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // 价格栏
                            Row(
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                const Text(
                                  '¥',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 16,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                Text(
                                  '${_selectedSku?.price ?? product.price}',
                                  style: const TextStyle(
                                    color: Colors.white,
                                    fontSize: 28,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                const Spacer(),
                                Text(
                                  '库存: ${_selectedSku?.stock ?? product.stock}',
                                  style: const TextStyle(
                                    color: Colors.grey,
                                    fontSize: 12,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 12),
                            Text(
                              product.name,
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.w600,
                                height: 1.3,
                              ),
                            ),
                            const SizedBox(height: 12),
                            Text(
                              product.description,
                              style: TextStyle(
                                color: Colors.grey[600],
                                fontSize: 14,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 10),

                      // 规格选择
                      if (product.skus.isNotEmpty)
                        Container(
                          color: const Color(0xFF1E1E1E),
                          padding: const EdgeInsets.all(16),
                          width: double.infinity,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Text(
                                "选择规格",
                                style: TextStyle(
                                  color: Colors.white,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 16,
                                ),
                              ),
                              const SizedBox(height: 12),
                              Wrap(
                                spacing: 8,
                                runSpacing: 8,
                                children: product.skus.map((sku) {
                                  final isSelected = _selectedSku?.id == sku.id;
                                  return ChoiceChip(
                                    label: Text(sku.name),
                                    selected: isSelected,
                                    onSelected: (selected) {
                                      setState(() {
                                        _selectedSku = selected ? sku : null;
                                      });
                                    },
                                    selectedColor: const Color(0xFFFF9000),
                                    backgroundColor: const Color(0xFF2C2C2C),
                                    labelStyle: TextStyle(
                                      color: isSelected
                                          ? Colors.black
                                          : Colors.white,
                                    ),
                                  );
                                }).toList(),
                              ),
                            ],
                          ),
                        ),

                      const SizedBox(height: 10),

                      // 数量选择
                      Container(
                        color: const Color(0xFF1E1E1E),
                        padding: const EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 12,
                        ),
                        child: Row(
                          children: [
                            const Text(
                              "数量",
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                            const Spacer(),
                            Container(
                              decoration: BoxDecoration(
                                border: Border.all(
                                  color: Colors.grey,
                                  width: 0.5,
                                ),
                                borderRadius: BorderRadius.circular(4),
                              ),
                              child: Row(
                                children: [
                                  IconButton(
                                    icon: const Icon(Icons.remove, size: 16),
                                    onPressed: _quantity > 1
                                        ? () => setState(() => _quantity--)
                                        : null,
                                  ),
                                  Text(
                                    '$_quantity',
                                    style: const TextStyle(fontSize: 16),
                                  ),
                                  IconButton(
                                    icon: const Icon(Icons.add, size: 16),
                                    onPressed: () =>
                                        setState(() => _quantity++),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),

                      const SizedBox(height: 10),

                      // 评价列表
                      Container(
                        color: const Color(0xFF1E1E1E),
                        width: double.infinity,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Padding(
                              padding: const EdgeInsets.all(16.0),
                              child: Row(
                                children: [
                                  const Text(
                                    "商品评价",
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: Colors.white,
                                    ),
                                  ),
                                  const Spacer(),
                                  reviewsAsyncValue.maybeWhen(
                                    data: (reviews) => Text(
                                      "(${reviews.length})",
                                      style: const TextStyle(
                                        color: Colors.grey,
                                        fontSize: 14,
                                      ),
                                    ),
                                    orElse: () => const SizedBox(),
                                  ),
                                ],
                              ),
                            ),
                            reviewsAsyncValue.when(
                              data: (reviews) {
                                if (reviews.isEmpty) {
                                  return const Padding(
                                    padding: EdgeInsets.fromLTRB(16, 0, 16, 16),
                                    child: Text(
                                      "暂无评价",
                                      style: TextStyle(color: Colors.grey),
                                    ),
                                  );
                                }
                                return ListView.separated(
                                  padding: EdgeInsets.zero,
                                  shrinkWrap: true,
                                  physics: const NeverScrollableScrollPhysics(),
                                  itemCount: reviews.length > 3
                                      ? 3
                                      : reviews.length, // 只显示前3条
                                  separatorBuilder: (ctx, index) =>
                                      const Divider(
                                        height: 1,
                                        color: Colors.black12,
                                      ),
                                  itemBuilder: (ctx, index) {
                                    final review = reviews[index];
                                    return ListTile(
                                      leading: CircleAvatar(
                                        backgroundImage:
                                            review.userAvatar.isNotEmpty
                                            ? NetworkImage(review.userAvatar)
                                            : null,
                                        child: review.userAvatar.isEmpty
                                            ? const Icon(Icons.person)
                                            : null,
                                      ),
                                      title: Row(
                                        children: [
                                          Text(
                                            review.userName,
                                            style: const TextStyle(
                                              fontSize: 14,
                                            ),
                                          ),
                                          const Spacer(),
                                          Text(
                                            review.createdAt.split('T')[0],
                                            style: const TextStyle(
                                              fontSize: 12,
                                              color: Colors.grey,
                                            ),
                                          ),
                                        ],
                                      ),
                                      subtitle: Column(
                                        crossAxisAlignment:
                                            CrossAxisAlignment.start,
                                        children: [
                                          Row(
                                            children: List.generate(
                                              5,
                                              (i) => Icon(
                                                Icons.star,
                                                size: 14,
                                                color: i < review.rating
                                                    ? Colors.amber
                                                    : Colors.grey,
                                              ),
                                            ),
                                          ),
                                          const SizedBox(height: 4),
                                          Text(review.content),
                                        ],
                                      ),
                                    );
                                  },
                                );
                              },
                              loading: () => const Padding(
                                padding: EdgeInsets.all(16.0),
                                child: CircularProgressIndicator(),
                              ),
                              error: (_, __) => const SizedBox(),
                            ),
                            if (reviewsAsyncValue.valueOrNull?.isNotEmpty ??
                                false)
                              Center(
                                child: TextButton(
                                  onPressed: () {
                                    // TODO: View all reviews
                                  },
                                  child: const Text("查看全部评价"),
                                ),
                              ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 10),

                      // 商品详情图
                      if (product.images.isNotEmpty)
                        Container(
                          color: const Color(0xFF1E1E1E),
                          width: double.infinity,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Padding(
                                padding: EdgeInsets.all(16.0),
                                child: Text(
                                  "商品详情",
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.white,
                                  ),
                                ),
                              ),
                              ...product.images.map(
                                (img) => Image.network(
                                  img,
                                  width: double.infinity,
                                  fit: BoxFit.cover,
                                ),
                              ),
                            ],
                          ),
                        ),
                    ],
                  ),
                ),
              ),
              // 底部操作栏
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 8,
                ),
                decoration: BoxDecoration(
                  color: const Color(0xFF1E1E1E),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withValues(alpha: 0.05),
                      blurRadius: 10,
                      offset: const Offset(0, -2),
                    ),
                  ],
                ),
                child: Row(
                  children: [
                    InkWell(
                      onTap: () {
                        ScaffoldMessenger.of(
                          context,
                        ).showSnackBar(const SnackBar(content: Text('进店逛逛')));
                      },
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: const [
                          Icon(Icons.store, color: Colors.grey),
                          Text(
                            "店铺",
                            style: TextStyle(fontSize: 10, color: Colors.grey),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 20),
                    InkWell(
                      onTap: () {
                        ScaffoldMessenger.of(
                          context,
                        ).showSnackBar(const SnackBar(content: Text('联系客服')));
                      },
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: const [
                          Icon(Icons.headset_mic, color: Colors.grey),
                          Text(
                            "客服",
                            style: TextStyle(fontSize: 10, color: Colors.grey),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 20),
                    Expanded(
                      child: SizedBox(
                        height: 40,
                        child: ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFFFF9000),
                            shape: const RoundedRectangleBorder(
                              borderRadius: BorderRadius.horizontal(
                                left: Radius.circular(20),
                              ),
                            ),
                          ),
                          onPressed: () {
                            if (product.skus.isNotEmpty &&
                                _selectedSku == null) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(content: Text('请选择规格')),
                              );
                              return;
                            }
                            ref
                                .read(cartProvider.notifier)
                                .addToCart(
                                  product.id,
                                  _quantity,
                                  skuId: _selectedSku?.id,
                                );
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('已加入购物车')),
                            );
                          },
                          child: const Text('加入购物车'),
                        ),
                      ),
                    ),
                    Expanded(
                      child: SizedBox(
                        height: 40,
                        child: ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFFFF5000),
                            shape: const RoundedRectangleBorder(
                              borderRadius: BorderRadius.horizontal(
                                right: Radius.circular(20),
                              ),
                            ),
                          ),
                          onPressed: () {
                            if (product.skus.isNotEmpty &&
                                _selectedSku == null) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                  content: Text('Please select a type'),
                                ),
                              );
                              return;
                            }
                            ref
                                .read(cartProvider.notifier)
                                .addToCart(
                                  product.id,
                                  _quantity,
                                  skuId: _selectedSku?.id,
                                );
                            context.go('/');
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('已加入购物车，请前往购物车结算')),
                            );
                          },
                          child: const Text('立即购买'),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => Center(child: Text('加载失败: $err')),
      ),
    );
  }
}
