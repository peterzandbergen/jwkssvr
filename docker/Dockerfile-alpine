# cir-cn.chp.belastingdienst.nl/zandp06/jwkssvr
FROM cir-cn.chp.belastingdienst.nl/cpet/golang:1.19 as build

WORKDIR /app

COPY . .
ADD ./certs /etc/ssl/certs

ENV http_proxy=http://proxyprokdmzbulk.belastingdienst.nl:8080
ENV https_proxy=http://proxyprokdmzbulk.belastingdienst.nl:8080
ENV no_proxy=localhost,127.0.0.1,*.belastingdienst.nl

RUN go mod download

EXPOSE 8080

RUN CGO_ENABLED=0 go build -o main ./cmd/jwkssvr

# docker pull cir-cn.chp.belastingdienst.nl/zandp06/docker.io/alpine:3

FROM cir-cn.chp.belastingdienst.nl/zandp06/docker.io/alpine:3

ADD ./certs /etc/ssl/certs

COPY --from=build /app/main /app/

ENTRYPOINT [ "/app/main" ]