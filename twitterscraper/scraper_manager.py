import json
import logging

from kafka import KafkaProducer, KafkaConsumer, TopicPartition

from .user_scraper import UserScraper


class ScraperManager:

    def __init__(
        self,
        kafka_host_port: str = "localhost:9092",
        new_user_topic: str = "user_names",
        scraped_user_topic: str = "users_scraped",
    ):
        self.user_consumer = KafkaConsumer(
            new_user_topic,
            bootstrap_servers=kafka_host_port,
        )
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_host_port,
            value_serializer=lambda v: json.dumps(v).encode('utf-8'),
        )
        self.new_user_topic = new_user_topic
        self.scraped_user_topic = scraped_user_topic

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

    def _consume_scrape_produce(self):
        """
        Consumes new user_names from kafka,
        scrapes these with twint,
        and produces/sends scraped users and unknown users to kafka
        """
        new_users = self.consume()
        logging.info(f"New users received: {new_users}")

        for user_name in new_users:
            self.scrape_and_produce(user_name)

    def consume(self, blocking: bool = True) -> dict:
        timeout_ms = float("inf") if blocking is True else 0
        partition_dict = self.user_consumer.poll(
            timeout_ms=timeout_ms,
            max_records=1,
        )

        ret = []
        for consumer_list in partition_dict.values():
            names = [json.loads(consumer.value) for consumer in consumer_list]
            ret.extend(names)
        return ret

    def scrape_and_produce(self, user_name: str) -> None:
        user = self.scrape(user_name)
        self.produce(user)

    def scrape(self, user_name: str):
        logging.info(f"scrape user {user_name}")
        user_scraper = UserScraper(user_name)
        user = user_scraper.scrape()
        return user

    def produce(self, user) -> None:
        self.send_scraped_user(user)
        self.send_new_users(user)

    def send_scraped_user(self, user) -> None:
        topic = self.scraped_user_topic
        logging.info(
            f"Send scraped user {user.username} to kafka/{topic}"
        )
        user_dict = user.__dict__
        self.producer.send(topic, user_dict)

    def send_new_users(self, user) -> None:
        followers_list = user.followers_list
        following_list = user.following_list

        ls = set(followers_list + following_list)
        topic = self.new_user_topic
        for user_name in ls:
            # TODO: check if username is already in topic
            logging.info(f"Send new user {user_name} to kafka/{topic}")
            self.producer.send(topic, user_name)


if __name__ == "__main__":
    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )

    scraper_manager = ScraperManager()
    scraper_manager.consume_scrape_produce()
