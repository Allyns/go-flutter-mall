import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:go_flutter_mall/features/product/providers/product_provider.dart';
import 'package:go_flutter_mall/features/product/providers/search_history_provider.dart';

class ProductListScreen extends ConsumerStatefulWidget {
  const ProductListScreen({super.key});

  @override
  ConsumerState<ProductListScreen> createState() => _ProductListScreenState();
}

class _ProductListScreenState extends ConsumerState<ProductListScreen> {
  final FocusNode _searchFocusNode = FocusNode();
  final TextEditingController _searchController = TextEditingController();
  bool _isSearchActive = false;

  @override
  void initState() {
    super.initState();
    _searchFocusNode.addListener(() {
      setState(() {
        _isSearchActive = _searchFocusNode.hasFocus;
      });
    });
  }

  @override
  void dispose() {
    _searchFocusNode.dispose();
    _searchController.dispose();
    super.dispose();
  }

  void _onSearchSubmitted(String value) {
    if (value.trim().isNotEmpty) {
      ref.read(searchHistoryProvider.notifier).addHistory(value);
      ref.read(productSearchQueryProvider.notifier).state = value;
      ref.invalidate(productListProvider);
    }
    _searchFocusNode.unfocus(); // 收起键盘
    // 保持 _isSearchActive 为 false，展示搜索结果
  }

  void _onHistoryTap(String keyword) {
    _searchController.text = keyword;
    _onSearchSubmitted(keyword);
  }

  Widget _buildCategoryItem(
    WidgetRef ref,
    String title,
    int id,
    bool isSelected,
  ) {
    return GestureDetector(
      onTap: () {
        ref.read(productCategoryProvider.notifier).state = id;
        ref.invalidate(productListProvider);
      },
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            title,
            style: TextStyle(
              color: isSelected ? Colors.white : Colors.white70,
              fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 4),
          if (isSelected) Container(width: 20, height: 2, color: Colors.white),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final productsAsyncValue = ref.watch(productListProvider);
    final selectedCategoryId = ref.watch(productCategoryProvider);
    final searchHistory = ref.watch(searchHistoryProvider);

    return Scaffold(
      body: Stack(
        children: [
          CustomScrollView(
            slivers: [
              // 仿淘宝顶部搜索栏
              SliverAppBar(
                pinned: true,
                floating: true,
                backgroundColor: Colors.black,
                expandedHeight: 100, // 增加高度以容纳 FlexibleSpace
                title: Container(
                  height: 36,
                  decoration: BoxDecoration(
                    color: const Color(0xFF2C2C2C),
                    borderRadius: BorderRadius.circular(18),
                  ),
                  child: TextField(
                    controller: _searchController,
                    focusNode: _searchFocusNode,
                    textInputAction: TextInputAction.search,
                    onSubmitted: _onSearchSubmitted,
                    onChanged: (value) {
                      // 可选：实时搜索或仅在提交时搜索
                      // ref.read(productSearchQueryProvider.notifier).state = value;
                    },
                    decoration: InputDecoration(
                      hintText: '搜索商品...',
                      hintStyle: TextStyle(
                        color: Colors.grey[400],
                        fontSize: 14,
                      ),
                      prefixIcon: const Icon(
                        Icons.search,
                        color: Colors.grey,
                        size: 20,
                      ),
                      suffixIcon: _searchController.text.isNotEmpty
                          ? GestureDetector(
                              onTap: () {
                                _searchController.clear();
                                // ref.read(productSearchQueryProvider.notifier).state = '';
                                // ref.invalidate(productListProvider);
                              },
                              child: const Icon(
                                Icons.clear,
                                color: Colors.grey,
                                size: 20,
                              ),
                            )
                          : null,
                      border: InputBorder.none,
                      enabledBorder: InputBorder.none,
                      focusedBorder: InputBorder.none,
                      contentPadding: const EdgeInsets.symmetric(vertical: 8),
                      fillColor: Colors.transparent,
                    ),
                    style: const TextStyle(fontSize: 14, color: Colors.white),
                  ),
                ),
                flexibleSpace: FlexibleSpaceBar(
                  background: Container(
                    decoration: const BoxDecoration(
                      gradient: LinearGradient(
                        colors: [Colors.black, Color(0xFF1B5E20)],
                        begin: Alignment.topLeft,
                        end: Alignment.bottomRight,
                      ),
                    ),
                    child: Align(
                      alignment: Alignment.bottomCenter,
                      child: Padding(
                        padding: const EdgeInsets.only(bottom: 10),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            _buildCategoryItem(
                              ref,
                              "全部",
                              0,
                              selectedCategoryId == 0,
                            ),
                            _buildCategoryItem(
                              ref,
                              "食品",
                              3,
                              selectedCategoryId == 3,
                            ), // 对应数据库 ID
                            _buildCategoryItem(
                              ref,
                              "生鲜",
                              4,
                              selectedCategoryId == 4,
                            ),
                            _buildCategoryItem(
                              ref,
                              "数码",
                              1,
                              selectedCategoryId == 1,
                            ),
                            _buildCategoryItem(
                              ref,
                              "服饰",
                              2,
                              selectedCategoryId == 2,
                            ),
                            _buildCategoryItem(
                              ref,
                              "家电",
                              5,
                              selectedCategoryId == 5,
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
              ),

              // 商品列表
              productsAsyncValue.when(
                data: (products) => SliverPadding(
                  padding: const EdgeInsets.all(10),
                  sliver: SliverGrid(
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2,
                          childAspectRatio: 0.7, // 调整比例让卡片更高，适合展示图片和信息
                          crossAxisSpacing: 10,
                          mainAxisSpacing: 10,
                        ),
                    delegate: SliverChildBuilderDelegate((context, index) {
                      final product = products[index];
                      return GestureDetector(
                        onTap: () => context.push('/product/${product.id}'),
                        child: Container(
                          decoration: BoxDecoration(
                            color: const Color(0xFF1E1E1E),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          clipBehavior: Clip.antiAlias,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              // 商品图片
                              Expanded(
                                child: Stack(
                                  children: [
                                    Image.network(
                                      product.coverImage,
                                      fit: BoxFit.cover,
                                      width: double.infinity,
                                      height: double.infinity,
                                      errorBuilder: (_, __, ___) => Container(
                                        color: Colors.grey[800],
                                        child: const Center(
                                          child: Icon(
                                            Icons.image_not_supported,
                                            color: Colors.grey,
                                          ),
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              // 商品信息
                              Padding(
                                padding: const EdgeInsets.all(10.0),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      product.name,
                                      maxLines: 2,
                                      overflow: TextOverflow.ellipsis,
                                      style: const TextStyle(
                                        fontSize: 14,
                                        color: Colors.white,
                                        fontWeight: FontWeight.w500,
                                      ),
                                    ),
                                    const SizedBox(height: 6),
                                    Row(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.baseline,
                                      textBaseline: TextBaseline.alphabetic,
                                      children: [
                                        const Text(
                                          '¥',
                                          style: TextStyle(
                                            color: Colors.green,
                                            fontSize: 12,
                                            fontWeight: FontWeight.bold,
                                          ),
                                        ),
                                        Text(
                                          '${product.price}',
                                          style: const TextStyle(
                                            color: Colors.green,
                                            fontSize: 18,
                                            fontWeight: FontWeight.bold,
                                          ),
                                        ),
                                        const SizedBox(width: 6),
                                        Text(
                                          '${product.stock} sold',
                                          style: TextStyle(
                                            color: Colors.grey[500],
                                            fontSize: 10,
                                          ),
                                        ),
                                      ],
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    }, childCount: products.length),
                  ),
                ),
                loading: () => const SliverFillRemaining(
                  child: Center(child: CircularProgressIndicator()),
                ),
                error: (err, stack) => SliverFillRemaining(
                  child: Center(child: Text('Error: $err')),
                ),
              ),
            ],
          ),
          // Search History Overlay
          if (_isSearchActive)
            Positioned(
              top: 100, // AppBar Height
              left: 0,
              right: 0,
              bottom: 0,
              child: Container(
                color: const Color(0xFF121212),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (searchHistory.isNotEmpty) ...[
                      Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text(
                              '历史搜索',
                              style: TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.bold,
                                color: Colors.white,
                              ),
                            ),
                            GestureDetector(
                              onTap: () {
                                ref
                                    .read(searchHistoryProvider.notifier)
                                    .clearHistory();
                              },
                              child: const Icon(
                                Icons.delete_outline,
                                size: 18,
                                color: Colors.grey,
                              ),
                            ),
                          ],
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16.0),
                        child: Wrap(
                          spacing: 8.0,
                          runSpacing: 8.0,
                          children: searchHistory.map((history) {
                            return GestureDetector(
                              onTap: () => _onHistoryTap(history.keyword),
                              child: Container(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 12,
                                  vertical: 6,
                                ),
                                decoration: BoxDecoration(
                                  color: Colors.grey[800],
                                  borderRadius: BorderRadius.circular(16),
                                ),
                                child: Text(
                                  history.keyword,
                                  style: const TextStyle(
                                    fontSize: 12,
                                    color: Colors.white70,
                                  ),
                                ),
                              ),
                            );
                          }).toList(),
                        ),
                      ),
                    ] else
                      const Padding(
                        padding: EdgeInsets.all(32.0),
                        child: Center(
                          child: Text(
                            '暂无搜索历史',
                            style: TextStyle(color: Colors.grey),
                          ),
                        ),
                      ),
                  ],
                ),
              ),
            ),
        ],
      ),
    );
  }
}
