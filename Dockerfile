FROM ubuntu:14.04

MAINTAINER JamesClonk

EXPOSE 3007

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY moviedb-frontend /moviedb-frontend
COPY public /public
COPY templates /templates

ENV JCIO_ENV production
ENV PORT 3007
ENV JCIO_MOVIEDB_BACKEND http://moviedb-backend.jamesclonk.io

CMD ["/moviedb-frontend"]
