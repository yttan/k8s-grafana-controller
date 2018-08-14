FROM ubuntu:16.04
RUN apt update -y
ENV CONFIG_PATH="/adapter"
ENV GRAFANA_IP="10.111.251.8"
ENV PROMETHEUS_IP="10.101.3.181"
ENV ADMIN_NAME="admin"
ENV ADMIN_PASSWORD="admin"


RUN mkdir -p /adapter
ADD main /adapter
ADD config /adapter
CMD exec /adapter/main

