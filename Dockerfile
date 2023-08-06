FROM ubuntu:latest

WORKDIR /app

COPY ./bin ./bin

RUN [ "chmod", "+x", "/app/bin/ib_db_services_lnx.exe"]

ENTRYPOINT ["./bin/ib_db_services_lnx.exe"]

EXPOSE 9191
