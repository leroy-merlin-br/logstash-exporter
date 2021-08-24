FROM golang:1.16.7-alpine3.14 as build

WORKDIR /src

COPY ./ /src/

RUN go build -a -tags netgo -o logstash_exporter .

FROM scratch as final

LABEL maintainer="boitata-sre@leroymerlin.com.br"
LABEL description="Logstash exporter"
LABEL org.opencontainers.image.authors="LMBR Site Reliability Team <boitata-sre@leroymerlin.com.br>"
LABEL org.opencontainers.image.source="https://github.com/leroy-merlin-br/logstash-exporter"
LABEL org.opencontainers.image.licenses="Copyright © 2021 Leroy Merlin Brasil"
LABEL org.opencontainers.image.vendor="Leroy Merlin Brasil"

COPY --from=build /src/logstash_exporter /

EXPOSE 9198

ENTRYPOINT ["/logstash_exporter"]
