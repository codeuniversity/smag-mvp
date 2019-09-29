import json
import logging
import sys

from kafka import KafkaProducer

if __name__ == "__main__":
    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )

    argv = sys.argv
    if len(argv) != 3:
        logging.error(
            "Wrong number of arguments.\nUsage: python insert_seed.py <kafka host:port> <seed_name>"
        )

    kafka_host_port = argv[1]
    producer = KafkaProducer(
        bootstrap_servers=kafka_host_port,
        value_serializer=lambda v: json.dumps(v).encode('utf-8'),
    )

    seed_name = argv[2]
    producer.send("user_names", seed_name)
    producer.flush()
