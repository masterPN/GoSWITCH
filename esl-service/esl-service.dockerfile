FROM golang:1.22.5

RUN mkdir /app

COPY eslApp /app

CMD [ "/app/eslApp" ]