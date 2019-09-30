import logging

import twint

from twitterscraper.scraper_manager import ScraperManager
from twitterscraper.utils import get_conf


def scrape(user_name: str) -> twint.user.user:
    conf = get_conf(user_name)
    user = scrape_user(conf)
    user.followers_list = scrape_follows_list(twint.run.Followers, conf)
    user.following_list = scrape_follows_list(twint.run.Following, conf)

    return user


def scrape_user(conf: twint.Config) -> twint.user.user:
    twint.run.Lookup(conf)
    user = twint.output.users_list.pop()
    return user


def scrape_follows_list(func, conf: twint.Config) -> list:
    func(conf)

    ret = twint.output.follows_list
    twint.output.follows_list = []
    return ret


class Scraper(ScraperManager):
    @staticmethod
    def scrape(user_name: str):
        logging.info(f"scrape user {user_name}")
        user = scrape(user_name)
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

    scraper_manager = Scraper(
        insert_topic=insert_topic,
        fetch_topic=fetch_topic,
        kafka_host_port=kafka_host_port,
    )
    scraper_manager.consume_scrape_produce()
