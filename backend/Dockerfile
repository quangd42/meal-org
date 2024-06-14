FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY bin/meal_planner /usr/bin/meal_planner

CMD ["meal_planner"]
