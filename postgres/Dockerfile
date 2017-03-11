FROM postgres:9.6-alpine

ENV POSTGRES_USER docker
ENV POSTGRES_PASSWORD docker
ENV POSTGRES_DB docker
EXPOSE 5432

ADD init.sql /docker-entrypoint-initdb.d/
