FROM ubuntu:16.04 AS opencv

# install cmake 3.11
RUN apt-get update && apt-get install -y --no-install-recommends wget && \
      wget --no-check-certificate https://cmake.org/files/v3.11/cmake-3.11.0-Linux-x86_64.tar.gz && \
      tar -xvf cmake-3.11.0-Linux-x86_64.tar.gz && \
      cd cmake-3.11.0-Linux-x86_64 && cp -r bin /usr/ &&  \
      cp -r share /usr/ && \
      cp -r doc /usr/share/ && cp -r man /usr/share/ && \
      cd .. && \
      rm -r cmake-3.11.0-Linux-x86_64 && \
      rm cmake-3.11.0-Linux-x86_64.tar.gz

RUN apt-get install -y --no-install-recommends \
            git build-essential pkg-config g++ unzip libgtk2.0-dev \
            curl ca-certificates libcurl4-openssl-dev libssl-dev \
            libavcodec-dev libavformat-dev libswscale-dev libtbb2 libtbb-dev \
            libjpeg-dev libpng-dev libtiff-dev libdc1394-22-dev && \
            rm -rf /var/lib/apt/lists/*

ARG OPENCV_VERSION="4.0.1"
ENV OPENCV_VERSION $OPENCV_VERSION

RUN curl -Lo opencv.zip https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip && \
            unzip -q opencv.zip && \
            curl -Lo opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip && \
            unzip -q opencv_contrib.zip && \
            rm opencv.zip opencv_contrib.zip && \
            cd opencv-${OPENCV_VERSION} && \
            mkdir build && cd build && \
            cmake -D CMAKE_BUILD_TYPE=RELEASE \
                  -D CMAKE_INSTALL_PREFIX=/usr/local \
                  -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules \
                  -D WITH_JASPER=OFF \
                  -D BUILD_DOCS=OFF \
                  -D BUILD_EXAMPLES=OFF \
                  -D BUILD_TESTS=OFF \
                  -D BUILD_PERF_TESTS=OFF \
                  -D BUILD_opencv_java=NO \
                  -D BUILD_opencv_python=NO \
                  -D BUILD_opencv_python2=NO \
                  -D BUILD_opencv_python3=NO \
                  -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
            make -j $(nproc --all) && \
            make preinstall && make install && ldconfig && \
            cd / && rm -rf opencv*

# clone seabolt-1.7.0 source code
RUN git clone -b v1.7.4 https://github.com/neo4j-drivers/seabolt.git /seabolt
# invoke cmake build and install artifacts - default location is /usr/local
WORKDIR /seabolt/build
# CMAKE_INSTALL_LIBDIR=lib is a hack where we override default lib64 to lib to workaround a defect
# in our generated pkg-config file
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install
RUN curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v0.3.1/gotestsum_0.3.1_linux_amd64.tar.gz" | tar -xz -C /usr/local/bin gotestsum

#################
#  Go + OpenCV  #
#################
FROM opencv

ARG GOVERSION="1.13"
ENV GOVERSION $GOVERSION

RUN apt-get update && apt-get install -y --no-install-recommends \
            git software-properties-common && \
            curl -Lo go${GOVERSION}.linux-amd64.tar.gz https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz && \
            tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz && \
            rm go${GOVERSION}.linux-amd64.tar.gz && \
            rm -rf /var/lib/apt/lists/*

# for circleci
ENV DOCKERVER="19.03.5"
RUN curl -L -o /tmp/docker-$DOCKERVER.tgz https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKERVER.tgz && \
      tar -xz -C /tmp -f /tmp/docker-$DOCKERVER.tgz && \
      mv /tmp/docker/* /usr/bin

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

