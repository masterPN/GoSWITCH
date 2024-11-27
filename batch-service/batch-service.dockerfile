FROM golang:1.22.5-alpine

RUN adduser -D -g '' appuser

WORKDIR /app

COPY batchApp /app

RUN chown -R appuser:appuser /app

USER appuser

CMD [ "/app/batchApp" ]