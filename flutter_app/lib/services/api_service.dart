import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/tracking_model.dart';

class ApiService {
  static const String baseUrl = 'http://localhost:8000';

  static Future<List<TrackingModel>> fetchTrackingData() async {
    final response = await http.get(Uri.parse('$baseUrl/tracking'));

    if (response.statusCode == 200) {
      final List data = jsonDecode(response.body);
      return data.map((item) => TrackingModel.fromJson(item)).toList();
    } else {
      throw Exception('Gagal mengambil data tracking');
    }
  }
}