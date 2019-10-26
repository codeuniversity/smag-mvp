import logging

import twint

from .scraper_manager import ScraperManager
from .utils import get_conf


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

    # if we only scrape user names (set conf.User_full = False) user names are in follows_list
    # if we scrape profiles of follows (set conf.User_full = True) user objs are in users_list
    ret = []
    ret.extend(twint.output.follows_list)
    ret.extend(twint.output.users_list)
    twint.output.follows_list = []
    twint.output.users_list = []
    return ret


class UserScraper(ScraperManager):
    name = "user_scraper"

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
    user_scraper.consume_scrape_produce()
