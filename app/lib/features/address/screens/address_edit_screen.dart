import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:go_flutter_mall/core/http/http_client.dart';
import 'package:go_flutter_mall/features/address/screens/address_list_screen.dart'; // 引用 Address 模型和 Provider

class AddressEditScreen extends HookConsumerWidget {
  final Address? address; // 如果为 null 则是添加，否则是编辑

  const AddressEditScreen({super.key, this.address});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final nameController = useTextEditingController(text: address?.receiverName ?? '');
    final phoneController = useTextEditingController(text: address?.phone ?? '');
    final provinceController = useTextEditingController(text: address?.province ?? '');
    final cityController = useTextEditingController(text: address?.city ?? '');
    final districtController = useTextEditingController(text: address?.district ?? '');
    final detailController = useTextEditingController(text: address?.detailAddress ?? '');
    final isDefault = useState(address?.isDefault ?? false);
    final isLoading = useState(false);

    return Scaffold(
      appBar: AppBar(title: Text(address == null ? '添加收货地址' : '编辑收货地址')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            TextField(
              controller: nameController,
              decoration: const InputDecoration(labelText: '收货人姓名'),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: phoneController,
              decoration: const InputDecoration(labelText: '手机号码'),
              keyboardType: TextInputType.phone,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: provinceController,
                    decoration: const InputDecoration(labelText: '省份'),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: TextField(
                    controller: cityController,
                    decoration: const InputDecoration(labelText: '城市'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            TextField(
              controller: districtController,
              decoration: const InputDecoration(labelText: '区/县'),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: detailController,
              decoration: const InputDecoration(labelText: '详细地址'),
            ),
            const SizedBox(height: 16),
            SwitchListTile(
              title: const Text('设为默认地址'),
              value: isDefault.value,
              onChanged: (val) => isDefault.value = val,
            ),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton(
                onPressed: isLoading.value
                    ? null
                    : () async {
                        if (nameController.text.isEmpty ||
                            phoneController.text.isEmpty ||
                            provinceController.text.isEmpty ||
                            cityController.text.isEmpty ||
                            districtController.text.isEmpty ||
                            detailController.text.isEmpty) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('请填写完整信息')),
                          );
                          return;
                        }

                        isLoading.value = true;
                        try {
                          final data = {
                            "receiver_name": nameController.text,
                            "phone": phoneController.text,
                            "province": provinceController.text,
                            "city": cityController.text,
                            "district": districtController.text,
                            "detail_address": detailController.text,
                            // 后端接口若支持 is_default 字段
                            // "is_default": isDefault.value, 
                          };

                          if (address == null) {
                            // 创建
                            await HttpClient().dio.post('/addresses', data: data);
                          } else {
                            // 更新
                            await HttpClient().dio.put('/addresses/${address!.id}', data: data);
                          }
                          
                          if (context.mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text(address == null ? '添加成功' : '修改成功')),
                            );
                            ref.invalidate(addressListProvider); // 刷新列表
                            context.pop();
                          }
                        } catch (e) {
                          if (context.mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text('操作失败: $e')),
                            );
                          }
                        } finally {
                          isLoading.value = false;
                        }
                      },
                child: isLoading.value
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('保存'),
              ),
            ),
            if (address != null) ...[
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: OutlinedButton(
                  style: OutlinedButton.styleFrom(foregroundColor: Colors.red),
                  onPressed: isLoading.value
                      ? null
                      : () async {
                          final confirm = await showDialog<bool>(
                            context: context,
                            builder: (context) => AlertDialog(
                              title: const Text('确认删除'),
                              content: const Text('确定要删除该地址吗？'),
                              actions: [
                                TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('取消')),
                                TextButton(onPressed: () => Navigator.pop(context, true), child: const Text('删除', style: TextStyle(color: Colors.red))),
                              ],
                            ),
                          );

                          if (confirm == true) {
                             isLoading.value = true;
                             try {
                               await HttpClient().dio.delete('/addresses/${address!.id}');
                               if (context.mounted) {
                                 ScaffoldMessenger.of(context).showSnackBar(
                                   const SnackBar(content: Text('删除成功')),
                                 );
                                 ref.invalidate(addressListProvider);
                                 context.pop();
                               }
                             } catch (e) {
                               if (context.mounted) {
                                 ScaffoldMessenger.of(context).showSnackBar(
                                   SnackBar(content: Text('删除失败: $e')),
                                 );
                               }
                             } finally {
                               isLoading.value = false;
                             }
                          }
                        },
                  child: const Text('删除地址'),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
