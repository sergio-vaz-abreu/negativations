FROM ubuntu:focal
RUN apt update && apt install -y ca-certificates wget gnupg
RUN wget -qO - https://pkgs-ce.cossacklabs.com/gpg | apt-key add -
RUN apt install -y apt-transport-https
RUN echo "deb https://pkgs-ce.cossacklabs.com/stable/ubuntu focal main" | \
      tee /etc/apt/sources.list.d/cossacklabs.list
RUN apt update && apt install -y libthemis-dev
WORKDIR /app
COPY ./.bin/negativations /app/
ENTRYPOINT ["/app/negativations"]
