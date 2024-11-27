FROM golang:1.22.5-alpine

RUN adduser -D -g '' appuser

WORKDIR /app

COPY redisApp /app

RUN chown -R appuser:appuser /app

USER appuser

CMD [ "/app/redisApp" ]