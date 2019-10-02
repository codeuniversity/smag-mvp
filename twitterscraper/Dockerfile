FROM python:3.7-slim

WORKDIR /src

COPY Pipfile Pipfile.lock ./
RUN pip install pipenv \
    && pipenv install

COPY ./twitterscraper ./twitterscraper

ENTRYPOINT [ "pipenv", "run", "python" ]
CMD [ "-c", "raise Eexception('Please set the CMD to either tweets_scraper.py or user_scraper.py')" ]
