FROM ubuntu:16.04
RUN apt update -y
ENV CONFIG_PATH="/controller"
ENV GRAFANA_IP="10.111.251.8"
ENV PROMETHEUS_IP="10.101.3.181"
ENV ADMIN_NAME="admin"
ENV ADMIN_PASSWORD="admin"


RUN mkdir -p /controller
ADD main /controller
ADD config /controller
CMD exec /controller/main
