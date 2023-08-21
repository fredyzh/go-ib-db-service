FROM ubuntu:latest

WORKDIR /app

COPY ./bin/ib_db_services_lnx.exe ./bin/ib_db_services_lnx.exe

RUN [ "chmod", "+x", "/app/bin/ib_db_services_lnx.exe"]

ENTRYPOINT ["./bin/ib_db_services_lnx.exe"]

EXPOSE 9191
