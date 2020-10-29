FROM ubuntu:bionic

LABEL maintainer="Pieter Droogendijk"

WORKDIR /app
ADD target/static ./static
COPY target/contentserver .

EXPOSE 8080

CMD ["./contentserver", "-addr=:8080", "-static=./static"]
