FROM golang:1.22.5-alpine

RUN mkdir /app

COPY mssqlApp /app

CMD [ "/app/mssqlApp" ]