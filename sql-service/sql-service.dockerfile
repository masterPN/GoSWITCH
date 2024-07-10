FROM golang:1.22.5

RUN mkdir /app

COPY sqlApp /app

CMD [ "/app/sqlApp" ]