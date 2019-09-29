import twint


class UserScraper:

    def __init__(self, user_name: str):
        self.user_name = user_name
        self.conf = self.get_conf()

    def scrape(self) -> twint.user.user:
        user = self.scrape_user()
        user.followers_list = self.scrape_follows_list(twint.run.Followers)
        user.following_list = self.scrape_follows_list(twint.run.Following)

        return user

    def scrape_user(self) -> twint.user.user:
        twint.run.Lookup(self.conf)
        user = twint.output.users_list.pop()
        return user

    def scrape_follows_list(self, func) -> list:
        func(self.conf)

        ret = twint.output.follows_list
        twint.output.follows_list = []
        return ret

    def get_conf(self) -> twint.config.Config:
        c = twint.Config()
        c.Username = self.user_name
        c.Store_object = True
        return c
