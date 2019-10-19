import json
import logging

from kafka import KafkaProducer, KafkaConsumer


class ScraperManager:

    name = "scraper_manager"

    def __init__(
        self,
        fetch_topic: str,
        insert_topic: str,
        kafka_consumer_group: str = "scraper_manager",
        kafka_host_port: str = "localhost:9092",
    ):
        self.consumer = KafkaConsumer(
            fetch_topic,
            bootstrap_servers=kafka_host_port,
            group_id=kafka_consumer_group,
            reconnect_backoff_ms=500,
            reconnect_backoff_max_ms=10000,
        )
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_host_port,
            value_serializer=lambda v: json.dumps(v).encode('utf-8'),
            reconnect_backoff_ms=500,
            reconnect_backoff_max_ms=10000,
        )
        self.insert_topic = insert_topic

    def consume_scrape_produce(self):
        try:
            while True:
                self._consume_scrape_produce()
        except Exception:
            logging.error(
                "Caught error. Going to flush KafkaProducer and then throw error further."
            )
            self.producer.flush()
            raise

    def _consume_scrape_produce(self) -> None:
        """
        Consumes from kafka,
        scrapes via twint,
        and produces/sends scraped msges to kafka
        """
        new_users = self.consume()
        logging.info(f"New users received: {new_users}")

        for user_name in new_users:
            self.scrape_and_produce(user_name)

    def consume(self, blocking: bool = True) -> dict:
        timeout_ms = float("inf") if blocking is True else 0
        partition_dict = self.consumer.poll(
            timeout_ms=timeout_ms,
            max_records=1,
        )

        ret = []
        for consumer_list in partition_dict.values():
            names = [consumer.value.decode("utf-8") for consumer in consumer_list]
            ret.extend(names)
        return ret

    def scrape_and_produce(self, user_name: str) -> None:
        msg = self.scrape(user_name)
        if type(msg) is list:
            for m in msg:
                self.produce(m)
        else:
            self.produce(msg)
        logging.info("done sending")

    def scrape(self, user_name: str):
        """This method will be implemented by the user to scrape either user-profile or tweets"""
        raise NotImplementedError(
            "You need to implement a scrape(user_name: str) method, "
            "which returns an object to be written to kafka."
        )

    def produce(self, msg) -> None:
        topic = self.insert_topic
        logging.info(
            f"{self.name} sends msg (from {msg.username}) to kafka/{topic}"
        )
        msg_dict = getattr(msg, "__dict__", msg)
        self.producer.send(topic, msg_dict)
