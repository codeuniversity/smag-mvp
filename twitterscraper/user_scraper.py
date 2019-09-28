import json

import twint


class UserScraper:

    def __init__(self, user_name: str):
        self.user_name = user_name
        self.conf = self.get_conf()

    def scrape(self) -> twint.user.user:
        user = self.scrape_user()
        user.followers_list = self.scrape_followers()
        user.following_list = self.scrape_following()

        self.user = user
        return user

    def scrape_user(self) -> twint.user.user:
        twint.run.Lookup(self.conf)
        user = twint.output.users_list.pop()
        return user

    def scrape_followers(self) -> list:
        twint.run.Followers(self.conf)

        ret = twint.output.follows_list
        twint.output.follows_list = []
        return ret

    def scrape_following(self) -> list:
        twint.run.Following(self.conf)

        ret = twint.output.follows_list
        twint.output.follows_list = []
        return ret

    def get_conf(self) -> twint.config.Config:
        c = twint.Config()
        c.Username = self.user_name
        c.Store_object = True
        return c

    def to_json(self) -> str:
        return json.dumps(self.user.__dict__).encode("utf-8")
