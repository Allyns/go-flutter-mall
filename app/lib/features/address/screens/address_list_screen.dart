import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';

// 地址模型
class Address {
  final int id;
  final String receiverName;
  final String phone;
  final String province;
  final String city;
  final String district;
  final String detailAddress;
  final bool isDefault;

  Address({
    required this.id,
    required this.receiverName,
    required this.phone,
    required this.province,
    required this.city,
    required this.district,
    required this.detailAddress,
    required this.isDefault,
  });

  factory Address.fromJson(Map<String, dynamic> json) {
    return Address(
      id: json['ID'] ?? 0,
      receiverName: json['receiver_name'] ?? '',
      phone: json['phone'] ?? '',
      province: json['province'] ?? '',
      city: json['city'] ?? '',
      district: json['district'] ?? '',
      detailAddress: json['detail_address'] ?? '',
      isDefault: json['is_default'] ?? false,
    );
  }

  String get fullAddress => "$province $city $district $detailAddress";
}

// 地址列表 Provider
final addressListProvider = FutureProvider.autoDispose<List<Address>>((ref) async {
  final response = await HttpClient().dio.get('/addresses');
  final List<dynamic> data = response.data;
  return data.map((json) => Address.fromJson(json)).toList();
});

class AddressListScreen extends ConsumerWidget {
  const AddressListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final addressAsync = ref.watch(addressListProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('收货地址管理')),
      body: addressAsync.when(
        data: (addresses) => addresses.isEmpty
            ? Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.location_off, size: 64, color: Colors.grey),
                    const SizedBox(height: 16),
                    const Text('暂无收货地址', style: TextStyle(color: Colors.grey)),
                    const SizedBox(height: 24),
                    ElevatedButton(
                      onPressed: () {
                        context.push('/addresses/add');
                      },
                      child: const Text('添加新地址'),
                    ),
                  ],
                ),
              )
            : ListView.separated(
                itemCount: addresses.length,
                separatorBuilder: (_, __) => const Divider(height: 1),
                itemBuilder: (context, index) {
                  final address = addresses[index];
                  return ListTile(
                    title: Row(
                      children: [
                        Text(address.receiverName, style: const TextStyle(fontWeight: FontWeight.bold)),
                        const SizedBox(width: 10),
                        Text(address.phone),
                        if (address.isDefault) ...[
                          const SizedBox(width: 10),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
                            decoration: BoxDecoration(
                              color: Colors.red[50],
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: const Text('默认', style: TextStyle(color: Colors.red, fontSize: 10)),
                          ),
                        ],
                      ],
                    ),
                    subtitle: Text(address.fullAddress),
                    trailing: IconButton(
                      icon: const Icon(Icons.edit),
                      onPressed: () {
                        context.push('/addresses/edit', extra: address);
                      },
                    ),
                    onTap: () {
                       // 如果是从结算页跳过来的，可能需要返回选中的地址
                       // 这里暂时简单处理，点击也进编辑，或者做成选择模式
                       // 为了支持选择，我们可以简单的 pop 并带回数据
                       // 但由于 GoRouter 的 push/pop 机制，最好是通过 provider 传递选中状态
                       // 或者这里仅作为管理页，选择逻辑在 Checkout 页内部处理
                       context.push('/addresses/edit', extra: address);
                    },
                  );
                },
              ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => Center(child: Text('加载失败: $err')),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          context.push('/addresses/add');
        },
        child: const Icon(Icons.add),
      ),
    );
  }
}
