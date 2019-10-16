FROM python:3.7-slim

WORKDIR /src

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY ./twitterscraper ./twitterscraper

ENTRYPOINT [ "python" ]
CMD [ "-c", "raise Exception('Please set the CMD to either tweets_scraper.py or user_scraper.py')" ]
