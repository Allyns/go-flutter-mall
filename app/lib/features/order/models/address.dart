class Address {
  final int id;
  final String name;
  final String phone;
  final String province;
  final String city;
  final String district;
  final String detailAddress;
  final bool isDefault;

  Address({
    required this.id,
    required this.name,
    required this.phone,
    required this.province,
    required this.city,
    required this.district,
    required this.detailAddress,
    required this.isDefault,
  });

  factory Address.fromJson(Map<String, dynamic> json) {
    return Address(
      id: json['id'],
      name: json['name'],
      phone: json['phone'],
      province: json['province'],
      city: json['city'],
      district: json['district'],
      detailAddress: json['detail_address'],
      isDefault: json['is_default'] ?? false,
    );
  }

  String get fullAddress => '$province $city $district $detailAddress';
}
