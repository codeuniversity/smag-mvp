import json
import logging
import traceback

from kafka import KafkaProducer, KafkaConsumer


class ScraperManager(object):

    name = "scraper_manager"

    def __init__(
        self,
        fetch_topic: str,
        insert_topic: str,
        kafka_consumer_group: str = "scraper_manager",
        kafka_address: str = "localhost:9092",
    ):
        self.consumer = KafkaConsumer(
            fetch_topic,
            bootstrap_servers=kafka_address,
            group_id=kafka_consumer_group,
            reconnect_backoff_ms=500,
            reconnect_backoff_max_ms=10000,
            max_poll_interval_ms=600000,
        )
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_address,
            value_serializer=lambda v: json.dumps(v).encode('utf-8'),
            reconnect_backoff_ms=500,
            reconnect_backoff_max_ms=10000,
            request_timeout_ms=600000,
        )
        self.insert_topic = insert_topic

    def run(self):
        try:
            while True:
                self.consume_scrape_produce()
        except Exception:
            logging.error(
                "Caught error. Going to flush KafkaProducer and then throw error further."
            )
            self.producer.flush()
            raise

    def consume_scrape_produce(self) -> None:
        """
        Consumes from kafka,
        scrapes via custom function,
        and produces/sends scraped msges to kafka
        """

        for msg in self.consumer:
            user_name = msg.value.decode("utf-8")
            try:
                self.scrape_and_produce(user_name)
            except Exception:
                self.consumer.commit()
                traceback.print_exc()
                logging.error(f"Couldn't scrape user {user_name}. Continuing")

    def scrape_and_produce(self, user_name: str) -> None:
        msg = self.scrape(user_name)
        msg_list = msg if type(msg) is list else [msg]
        for m in msg_list:
            self.produce(m)
        logging.info(
            f"Done sending {len(msg_list)} element(s) to kafka/{self.insert_topic}"
        )

    def scrape(self, user_name: str):
        """This method will be implemented by the user to scrape either user-profile or tweets"""
        raise NotImplementedError(
            "You need to implement a scrape(user_name: str) method, "
            "which returns an object to be written to kafka."
        )

    def produce(self, msg) -> None:
        topic = self.insert_topic
        logging.debug(
            f"{self.name} sends msg (from {msg.username}) to kafka/{topic}"
        )
        msg_dict = getattr(msg, "__dict__")
        self.producer.send(topic, msg_dict)
