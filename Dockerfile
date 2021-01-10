FROM golang:1.15 AS build

ADD ./ /opt/build/golang
WORKDIR /opt/build/golang
RUN go install ./app

FROM ubuntu:20.04 AS release

MAINTAINER Dmitry Kovalev

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
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
