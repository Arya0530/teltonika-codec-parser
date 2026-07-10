import json
import os
import urllib.request

CLICKHOUSE_URL = os.getenv(
    "CLICKHOUSE_URL",
    "http://127.0.0.1:8123/?user=arya&password=123456",
)


def execute_clickhouse(query):
    req = urllib.request.Request(
        CLICKHOUSE_URL,
        data=query.encode("utf-8"),
        method="POST"
    )
    return urllib.request.urlopen(req).read().decode("utf-8")


def insert_clickhouse(data):
    query = """
INSERT INTO raw_telemetry.avl_raw
(imei, timestamp, latitude, longitude, speed, raw_json, priority, altitude, angle, satellites, processed)
FORMAT JSONEachRow
"""

    row = {
        "imei": data["imei"],
        "timestamp": data["timestamp"],
        "latitude": data["latitude"],
        "longitude": data["longitude"],
        "speed": data["speed"],
        "raw_json": json.dumps(data),
        "priority": data["priority"],
        "altitude": data["altitude"],
        "angle": data["angle"],
        "satellites": data["satellites"],
        "processed": 0,
    }

    execute_clickhouse(query + "\n" + json.dumps(row))
    print("[ClickHouse] Raw data inserted from RabbitMQ")


def fetch_clickhouse_raw(limit=20):
    query = f"""
SELECT imei, timestamp, latitude, longitude, speed, priority, altitude, angle, satellites
FROM raw_telemetry.avl_raw
WHERE processed = 0
ORDER BY created_at ASC
LIMIT {limit}
FORMAT JSONEachRow
"""

    response = execute_clickhouse(query)

    rows = []
    for line in response.strip().splitlines():
        if line.strip():
            rows.append(json.loads(line))

    return rows


def mark_as_processed(row):
    query = f"""
ALTER TABLE raw_telemetry.avl_raw
UPDATE processed = 1
WHERE imei = '{row["imei"]}'
  AND timestamp = toDateTime('{row["timestamp"]}')
  AND latitude = {row["latitude"]}
  AND longitude = {row["longitude"]}
  AND speed = {row["speed"]}
"""
    execute_clickhouse(query)
    print("[ClickHouse] Raw data marked as processed")