import twint


def get_conf(user_name: str) -> twint.config.Config:
    c = twint.Config()
    c.Username = user_name
    c.Store_object = True
    c.Hide_output = True
    return c


class ShallowTwitterUser(object):

    # full user object contains

    # id: str = ""
    # url: str = ""
    # type: str = ""
    # name: str = ""
    # username: str = ""
    # bio: str = ""
    # avatar: str = ""
    # background_image: str = ""
    # location: str = ""
    # join_date: str = ""
    # join_time: str = ""
    # is_private: int = 0
    # is_verified: int = 0
    # following: int = 0
    # following_list: List[str] = [""]
    # followers: int = 0
    # followers_list: List[str] = [""]
    # tweets: int = 0
    # likes: int = 0
    # media_count: int = 0

    def __init__(self, username):
        self.username = username
