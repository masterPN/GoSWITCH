FROM golang:1.22.5-alpine

RUN mkdir /app

COPY batchApp /app

CMD [ "/app/batchApp" ]