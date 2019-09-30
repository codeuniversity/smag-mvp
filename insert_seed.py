if __name__ == "__main__":
    import json
    import logging
    import os

    from kafka import KafkaProducer

    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )

    kafka_host_port = os.getenv("KAFKA_HOST_PORT", "localhost:9092")
    producer = KafkaProducer(
        bootstrap_servers=kafka_host_port,
        value_serializer=lambda v: json.dumps(v).encode('utf-8'),
    )

    seed_name = os.getenv("SEED_NAME", "urhengula5")
    producer.send("user_names", seed_name)
    producer.flush()
