FROM python:3.6-slim-stretch

RUN apt-get -y update
RUN apt-get install -y --fix-missing \
  build-essential \
  cmake \
  gfortran \
  git \
  wget \
  curl \
  graphicsmagick \
  libgraphicsmagick1-dev \
  libatlas-dev \
  libavcodec-dev \
  libavformat-dev \
  libgtk2.0-dev \
  libjpeg-dev \
  liblapack-dev \
  libswscale-dev \
  pkg-config \
  python3-dev \
  python3-numpy \
  software-properties-common \
  zip \
  && apt-get clean && rm -rf /tmp/* /var/tmp/*

RUN cd ~ && \
  mkdir -p dlib && \
  git clone -b 'v19.9' --single-branch https://github.com/davisking/dlib.git dlib/ && \
  cd  dlib/ && \
  python3 setup.py install --yes USE_AVX_INSTRUCTIONS

WORKDIR /src

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY recognizer_pb2_grpc.py .
COPY recognizer_pb2.py .
COPY recognizer.py .
COPY metrics.py .
COPY server.py .

CMD [ "python", "server.py" ]
