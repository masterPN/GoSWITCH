FROM golang:1.22.5-alpine

RUN mkdir /app

COPY eslApp /app

CMD [ "/app/eslApp" ]