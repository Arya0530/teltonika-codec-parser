class TrackingModel {
  final int id;
  final String imei;
  final String timestamp;
  final double latitude;
  final double longitude;
  final int speed;
  final String status;
  final int priority;
  final int altitude;
  final int angle;
  final int satellites;

  TrackingModel({
    required this.id,
    required this.imei,
    required this.timestamp,
    required this.latitude,
    required this.longitude,
    required this.speed,
    required this.status,
    required this.priority,
    required this.altitude,
    required this.angle,
    required this.satellites,
  });

  factory TrackingModel.fromJson(Map<String, dynamic> json) {
    return TrackingModel(
      id: json['id'],
      imei: json['imei'],
      timestamp: json['timestamp'].toString(),
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      speed: json['speed'],
      status: json['status'],
      priority: json['priority'],
      altitude: json['altitude'],
      angle: json['angle'],
      satellites: json['satellites'],
    );
  }
}