FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY bin/mealorg_server /usr/bin/mealorg_server
COPY assets/ /assets/

CMD ["mealorg_server"]
