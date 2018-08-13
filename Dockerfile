FROM ubuntu:16.04
RUN apt update -y
ENV CONFIG_PATH="/adapter"
ENV GRAFANA_IP="10.111.251.8"
ENV PROMETHEUS_IP="10.101.3.181"

RUN mkdir -p /adapter/.kube
ADD grafana_adapter /adapter
ADD config /adapter/.kube
CMD exec /adapter/grafana_adapter

