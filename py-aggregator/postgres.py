import os
import psycopg2


def get_connection():
    return psycopg2.connect(
        host=os.getenv("POSTGRES_HOST", "localhost"),
        database=os.getenv(
            "POSTGRES_DB",
            "iot_teltonika",
        ),
        user=os.getenv("POSTGRES_USER", "postgres"),
        password=os.getenv(
            "POSTGRES_PASSWORD",
            "admin",
        ),
        port=int(os.getenv("POSTGRES_PORT", "5434")),
    )


def get_device_id_by_imei(cur, imei):
    cur.execute("""
        SELECT id
        FROM devices
        WHERE imei = %s
          AND is_active = TRUE
        LIMIT 1
    """, (imei,))

    row = cur.fetchone()

    if row is None:
        return None

    return row[0]


def insert_postgres(data):
    conn = get_connection()
    cur = conn.cursor()

    try:
        imei = data["imei"]

        device_id = get_device_id_by_imei(cur, imei)

        if device_id is None:
            print(f"[PostgreSQL] Device IMEI tidak valid / tidak aktif: {imei}")
            conn.rollback()
            return

        cur.execute("""
            INSERT INTO processed_avl
            (device_id, timestamp, latitude, longitude, speed, status,
             priority, altitude, angle, satellites)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            device_id,
            data["timestamp"],
            data["latitude"],
            data["longitude"],
            data["speed"],
            data["status"],
            data["priority"],
            data["altitude"],
            data["angle"],
            data["satellites"]
        ))

        conn.commit()
        print("[PostgreSQL] Clean data inserted with device relation")

    except Exception as e:
        conn.rollback()
        print("[PostgreSQL] Insert error:", e)

    finally:
        cur.close()
        conn.close()