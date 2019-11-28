import json
import logging
import os
from time import sleep

from kafka import KafkaProducer


def main():
    kafka_host_port = os.getenv("KAFKA_HOST_PORT", "localhost:9092")
    seed_name = os.getenv("SEED_NAME", "wpbdry")
    insert_topic = os.getenv("KAFKA_INSERT_TOPIC", "user_names")
    wait = int(os.getenv("SLEEP_SECONDS", "0"))

    logging.info(f"sleep for {wait} seconds")
    sleep(wait)

    producer = KafkaProducer(
        bootstrap_servers=kafka_host_port,
        value_serializer=lambda v: json.dumps(v).encode('utf-8'),
        reconnect_backoff_ms=500,
        reconnect_backoff_max_ms=5000,
    )

    logging.info(f"sleep for {wait} seconds")
    sleep(wait)

    logging.info(f"Send user_name {seed_name} to kafka/{insert_topic}")
    producer.send(insert_topic, seed_name)
    producer.flush()


if __name__ == "__main__":
    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )
    main()
