import logging

import twint

from .scraper_manager import ScraperManager
from .utils import get_conf


class TweetsScraper(ScraperManager):
    name = "tweets_scraper"

    @staticmethod
    def scrape(user_name: str):
        logging.info(f"Scrape tweets of user {user_name}")

        tweets = []

        c = get_conf(user_name)
        c.Store_object_tweets_list = tweets

        twint.run.Search(c)
        return tweets


if __name__ == "__main__":
    import os

    log_level = logging.DEBUG if os.getenv("DEBUG", "false") == "true" else logging.INFO

    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=log_level,
    )

    insert_topic = os.getenv("KAFKA_INSERT_TOPIC", "users_scraped")
    fetch_topic = os.getenv("KAFKA_FETCH_TOPIC", "user_names")
    kafka_consumer_group = os.getenv("KAFKA_CONSUMER_GROUP", "tweets_scraper")
    kafka_address = os.getenv("KAFKA_ADDRESS", "localhost:9092")

    tweets_scraper = TweetsScraper(
        insert_topic=insert_topic,
        fetch_topic=fetch_topic,
        kafka_consumer_group=kafka_consumer_group,
        kafka_address=kafka_address,
    )
    tweets_scraper.run()
