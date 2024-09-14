FROM postgres:16-alpine

ENV POSTGRES_USER=gorm
ENV POSTGRES_PASSWORD=admin
ENV POSTGRES_DB=database

COPY init-db.sql /docker-entrypoint-initdb.d/

EXPOSE 5432