import logging

from twitterscraper.user_scraper import UserScraper
from twitterscraper.scraper_manager import ScraperManager


class Scraper(ScraperManager):
    @staticmethod
    def scrape(user_name: str):
        logging.info(f"scrape user {user_name}")
        user_scraper = UserScraper(user_name)
        user = user_scraper.scrape()
        return user


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

    logging.info(
        f"# ENV VARS\n{insert_topic}\n{fetch_topic}\n{kafka_host_port}")

    scraper_manager = Scraper(
        insert_topic=insert_topic,
        fetch_topic=fetch_topic,
        kafka_host_port=kafka_host_port,
    )
    scraper_manager.consume_scrape_produce()
