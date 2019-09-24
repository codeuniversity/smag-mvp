import twint


if __name__ == "__main__":

    # Configure
    c = twint.Config()
    c.Username = "noneprivacy"
    c.Search = "#osint"
    c.Format = "Tweet id: {id} | Tweet: {tweet}"

    # Run
    twint.run.Search(c)
