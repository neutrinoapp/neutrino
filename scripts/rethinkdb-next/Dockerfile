FROM ubuntu:14.04

WORKDIR /rethink

COPY rethinkdb.deb rethinkdb.deb

RUN apt-get update -y
RUN apt-get install -y libcurl3 libprotobuf8
RUN dpkg -i rethinkdb.deb

CMD rethinkdb --bind all
