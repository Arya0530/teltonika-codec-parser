import json
import os

import pika

from clickhouse import insert_clickhouse

connection = pika.BlockingConnection(
    pika.ConnectionParameters(
        host=os.getenv("RABBITMQ_HOST", "127.0.0.1"),
        port=int(os.getenv("RABBITMQ_PORT", "5672")),
        virtual_host=os.getenv("RABBITMQ_VHOST", "/"),
        credentials=pika.PlainCredentials(
            os.getenv("RABBITMQ_USER", "guest"),
            os.getenv("RABBITMQ_PASSWORD", "guest"),
        ),
    )
)

channel = connection.channel()
channel.queue_declare(queue="avl_raw_queue", durable=True)


def callback(ch, method, properties, body):
    data = json.loads(body)

    print("==============================")
    print("MESSAGE DARI RABBITMQ")
    print(json.dumps(data, indent=4))

    try:
        insert_clickhouse(data)
        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as error:
        print("[Consumer] Gagal insert ClickHouse:", error)

        # Message dikembalikan ke queue.
        ch.basic_nack(
            delivery_tag=method.delivery_tag,
            requeue=True,
        )


channel.basic_qos(prefetch_count=1)

channel.basic_consume(
    queue="avl_raw_queue",
    on_message_callback=callback,
)

print("Waiting message for ClickHouse writer...")
channel.start_consuming()