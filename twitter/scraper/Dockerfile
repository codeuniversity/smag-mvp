FROM python:3.7-slim

WORKDIR /src

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY ./twitterscraper ./twitterscraper

ENTRYPOINT [ "python" ]
CMD [ "-c", "raise Exception('Please set the CMD to either `-m twitterscraper.posts_scraper.py`, `-m twitterscraper.users_scraper.py`, `-m twitterscraper.follwers_scraper.py` or `-m twitterscraper.follwing_scraper.py`')" ]
