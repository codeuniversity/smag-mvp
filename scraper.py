from twitterscraper.user_scraper import UserScraper


if __name__ == "__main__":
    user_scraper = UserScraper("urhengula5")
    user_scraper.scrape()
    print(user_scraper.to_json())
