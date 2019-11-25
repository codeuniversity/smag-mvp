FROM ubuntu as plugin-builder
RUN apt update && apt install -y git maven openjdk-8-jdk
WORKDIR /src
RUN git clone -b es-7.1 --single-branch https://github.com/lior-k/fast-elasticsearch-vector-scoring.git
WORKDIR /src/fast-elasticsearch-vector-scoring
RUN mvn package

FROM elasticsearch:7.1.0
COPY --from=plugin-builder /src/fast-elasticsearch-vector-scoring/target/releases/elasticsearch-binary-vector-scoring-7.1.0.zip /plugins/
RUN bin/elasticsearch-plugin install file:///plugins/elasticsearch-binary-vector-scoring-7.1.0.zip
