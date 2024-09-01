FROM golang:1.22.5-alpine

RUN mkdir /app

COPY redisApp /app

CMD [ "/app/redisApp" ]