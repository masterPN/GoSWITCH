FROM mcr.microsoft.com/mssql/server:2022-latest

USER root

COPY setup.sql .
COPY entrypoint.sh .
COPY import-data.sh .

# RUN dos2unix *
RUN chmod +x ./import-data.sh

USER mssql

ENTRYPOINT /bin/bash ./entrypoint.sh