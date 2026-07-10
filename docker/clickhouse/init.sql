CREATE DATABASE IF NOT EXISTS raw_telemetry;

CREATE TABLE IF NOT EXISTS raw_telemetry.avl_raw
(
    imei String,
    timestamp DateTime,
    latitude Float64,
    longitude Float64,
    speed UInt16,
    raw_json String,
    priority UInt8,
    altitude UInt16,
    angle UInt16,
    satellites UInt8,
    created_at DateTime DEFAULT now(),
    processed UInt8 DEFAULT 0
)
ENGINE = MergeTree()
ORDER BY (imei, timestamp);