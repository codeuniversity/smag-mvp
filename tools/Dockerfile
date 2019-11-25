FROM ubuntu

RUN apt update && apt install curl wget file xz-utils postgresql-client zsh git kafkacat jq -y

# install github.com/fgeller/kt
RUN wget https://github.com/fgeller/kt/releases/download/v12.1.0/kt-v12.1.0-linux-amd64.txz && \
  cat kt-v12.1.0-linux-amd64.txz | unxz > kt-v12.1.0-linux-amd64 && \
  tar -xvf kt-v12.1.0-linux-amd64 && \
  mv kt /usr/local/bin && \
  rm kt-v12.1.0-linux-amd64.txz && \
  rm kt-v12.1.0-linux-amd64

RUN zsh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
COPY .zshrc /root

WORKDIR /home/tools

ENTRYPOINT [ "zsh"]
