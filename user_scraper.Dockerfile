FROM python:3.7

WORKDIR /src

COPY Pipfile Pipfile.lock ./
RUN pip install pipenv \
    && pipenv install

COPY ./twitterscraper ./twitterscraper
COPY user_scraper.py .

ENTRYPOINT [ "pipenv", "run", "python" ]
CMD [ "user_scraper.py" ]
