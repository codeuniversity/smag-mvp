import logging

import twint

from .scraper_manager import ScraperManager
from .utils import get_conf


class TweetsScraper(ScraperManager):
    name = "tweets_scraper"

    @staticmethod
    def scrape(user_name: str):
        logging.info(f"scrape tweets of user {user_name}")

        tweets = []

        c = get_conf(user_name)
        c.Store_object_tweets_list = tweets

        twint.run.Search(c)
        return tweets


if __name__ == "__main__":
    import os

    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )

    insert_topic = os.getenv("KAFKA_INSERT_TOPIC", "users_scraped")
    fetch_topic = os.getenv("KAFKA_FETCH_TOPIC", "user_names")
    kafka_host_port = os.getenv("KAFKA_HOST_PORT", "localhost:9092")

    tweets_scraper = TweetsScraper(
        insert_topic=insert_topic,
        fetch_topic=fetch_topic,
        kafka_host_port=kafka_host_port,
    )
    tweets_scraper.consume_scrape_produce()