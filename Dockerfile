FROM golang:1.14-stretch AS build
ADD ./ /opt/build/golang
WORKDIR /opt/build/golang
RUN go install ./app
FROM ubuntu:18.04 AS release

MAINTAINER Dmitry Kovalev

ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER newuser WITH SUPERUSER PASSWORD 'password';" &&\
    createdb -O newuser forums &&\
    /etc/init.d/postgresql stop
EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root
EXPOSE 5000

COPY --from=build go/bin/app /usr/bin/
COPY --from=build /opt/build/golang/database /database/
CMD service postgresql start && app
