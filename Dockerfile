FROM ubuntu:14.04

EXPOSE 3007

ADD moviedb-frontend /moviedb-frontend
ADD public /public
ADD templates /templates

CMD ["/moviedb-frontend"]
