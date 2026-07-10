import os

import psycopg2
from psycopg2.extras import RealDictCursor


def get_connection():
    return psycopg2.connect(
        host=os.getenv("POSTGRES_HOST", "localhost"),
        database=os.getenv("POSTGRES_DB", "iot_teltonika"),
        user=os.getenv("POSTGRES_USER", "postgres"),
        password=os.getenv("POSTGRES_PASSWORD", "admin"),
        port=int(os.getenv("POSTGRES_PORT", "5434")),
    )


def get_tracking_data():
    conn = get_connection()
    cur = conn.cursor(cursor_factory=RealDictCursor)

    try:
        cur.execute("""
            SELECT
                p.id,
                d.imei,
                d.device_name,
                p.timestamp,
                p.latitude,
                p.longitude,
                p.speed,
                p.status,
                p.priority,
                p.altitude,
                p.angle,
                p.satellites
            FROM processed_avl p
            JOIN devices d
              ON p.device_id = d.id
            ORDER BY p.id DESC
            LIMIT 20
        """)

        return cur.fetchall()
    finally:
        cur.close()
        conn.close()