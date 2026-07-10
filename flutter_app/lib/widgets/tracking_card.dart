import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:latlong2/latlong.dart';
import '../models/tracking_model.dart';

class TrackingCard extends StatelessWidget {
  final TrackingModel data;

  const TrackingCard({super.key, required this.data});

  @override
  Widget build(BuildContext context) {
    final point = LatLng(data.latitude, data.longitude);
    final isMoving = data.status.toLowerCase() == 'moving';

    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(18),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.06),
            blurRadius: 12,
            offset: const Offset(0, 6),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('IMEI: ${data.imei}', style: const TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),

          SizedBox(
            height: 180,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(14),
              child: FlutterMap(
                options: MapOptions(
                  initialCenter: point,
                  initialZoom: 13,
                ),
                children: [
                  TileLayer(
                    urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                    userAgentPackageName: 'com.example.flutter_app',
                  ),
                  MarkerLayer(
                    markers: [
                      Marker(
                        point: point,
                        width: 45,
                        height: 45,
                        child: const Text('🚗', style: TextStyle(fontSize: 32)),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),

            const SizedBox(height: 12),
            Text('📍 Lat: ${data.latitude}'),
            Text('📍 Long: ${data.longitude}'),
            Text('⚡ Speed: ${data.speed} km/h'),
            Text('🚨 Priority: ${data.priority}'),
            Text('⛰️ Altitude: ${data.altitude} m'),
            Text('🧭 Angle: ${data.angle}°'),
            Text('🛰️ Satellites: ${data.satellites}'),
            Text(isMoving ? '🟢 Status: MOVING' : '🔴 Status: ${data.status.toUpperCase()}'),
            Text('🕒 Time: ${data.timestamp}'),
        ],
      ),
    );
  }
}