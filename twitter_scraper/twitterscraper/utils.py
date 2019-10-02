import twint


def get_conf(user_name: str) -> twint.config.Config:
    c = twint.Config()
    c.Username = user_name
    c.Store_object = True
    c.Hide_output = True
    return c
