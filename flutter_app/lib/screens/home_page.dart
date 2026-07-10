import 'dart:async';
import 'package:flutter/material.dart';
import '../models/tracking_model.dart';
import '../services/api_service.dart';
import '../widgets/tracking_card.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  late Future<List<TrackingModel>> _future;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _load();
    _timer = Timer.periodic(const Duration(seconds: 10), (_) => _load());
  }

  void _load() {
    setState(() {
      _future = ApiService.fetchTrackingData();
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F8FC),
      appBar: AppBar(
        backgroundColor: const Color(0xFF0A4FAF),
        foregroundColor: Colors.white,
        title: const Text('Widya Tracking Demo'),
      ),
      body: FutureBuilder<List<TrackingModel>>(
        future: _future,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting && !snapshot.hasData) {
            return const Center(child: CircularProgressIndicator());
          }

          if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          }

          final data = snapshot.data ?? [];

          return RefreshIndicator(
            onRefresh: () async => _load(),
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                const Text(
                  'Vehicle Telemetry',
                  style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 6),
                const Text('Auto refresh setiap 10 detik'),
                const SizedBox(height: 16),
                ...data.map((item) => TrackingCard(data: item)),
              ],
            ),
          );
        },
      ),
    );
  }
}