FROM docker.io/alpine:3.5

MAINTAINER bobliu0909@gmail.com

RUN mkdir -p /opt/app/gtlservice/etc

COPY gtlgateway /opt/app/gtlservice

COPY etc/config.yaml /opt/app/gtlservice/etc/config.yaml

RUN chmod +x /opt/app/gtlservice/gtlgateway

WORKDIR /opt/app/gtlservice

CMD ["./gtlgateway"]

EXPOSE 30000