if __name__ == "__main__":
    import json
    import logging
    import os
    from time import sleep

    from kafka import KafkaProducer
    from kafka.errors import NoBrokersAvailable

    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )

    kafka_host_port = os.getenv("KAFKA_HOST_PORT", "localhost:9092")

    wait = int(os.getenv("SLEEP_SECONDS", "5"))
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

    seed_name = os.getenv("SEED_NAME", "urhengula5")
    producer.send("user_names", seed_name)
    producer.flush()
