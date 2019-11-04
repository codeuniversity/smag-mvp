import logging

import twint

from .scraper_manager import ScraperManager
from .utils import get_conf


def scrape(user_name: str) -> twint.user.user:
    conf = get_conf(user_name)
    twint.run.Lookup(conf)
    user = twint.output.users_list.pop()
    return user


class UserScraper(ScraperManager):
    name = "user_scraper"

    @staticmethod
    def scrape(user_name: str):
        logging.info(f"Scrape user {user_name}")
        user = scrape(user_name)
        return user


if __name__ == "__main__":
    import os

    log_level = logging.DEBUG if os.getenv("DEBUG", "false") == "true" else logging.INFO

    logging.basicConfig(
        format="%(asctime)s.%(msecs)03d - %(module)s - %(levelname)s - %(message)s",
        datefmt="%H:%M:%S",
        level=log_level,
    )


    fetch_topic = os.getenv("KAFKA_FETCH_TOPIC", "user_names")
    insert_topic = os.getenv("KAFKA_INSERT_TOPIC", "users_scraped")
    kafka_consumer_group = os.getenv("KAFKA_CONSUMER_GROUP", "user_scraper")
    kafka_address = os.getenv("KAFKA_ADDRESS", "localhost:9092")

    user_scraper = UserScraper(
        insert_topic=insert_topic,
        fetch_topic=fetch_topic,
        kafka_consumer_group=kafka_consumer_group,
        kafka_address=kafka_address,
    )
    user_scraper.run()
