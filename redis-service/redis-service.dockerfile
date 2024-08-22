FROM golang:1.22.5

RUN mkdir /app

COPY redisApp /app

CMD [ "/app/redisApp" ]