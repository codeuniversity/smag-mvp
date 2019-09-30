FROM python:3.7-slim

RUN pip install kafka-python

WORKDIR /src
COPY insert_seed.py .

ENTRYPOINT [ "python" ]
CMD [ "insert_seed.py" ]
