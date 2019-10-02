FROM python:3.7-slim

WORKDIR /src

COPY Pipfile Pipfile.lock ./
RUN pip install pipenv \
    && pipenv install

COPY ./twitter_scraper ./twitter_scraper

ENTRYPOINT [ "pipenv", "run", "python" ]
CMD [ "-c", "raise Eexception('Please set the CMD to either tweets_scraper.py or user_scraper.py')" ]
