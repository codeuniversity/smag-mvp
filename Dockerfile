FROM python:3.7

WORKDIR /src

COPY Pipfile Pipfile.lock ./
RUN pip install pipenv \
    && pipenv install

COPY ./twitterscraper ./twitterscraper
COPY user_scraper.py tweets_scraper.py ./

ENTRYPOINT [ "pipenv", "run", "python" ]
CMD [ "-c", "raise Eexception('Please set the CMD to either tweets_scraper.py or user_scraper.py')" ]
