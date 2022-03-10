FROM golang:1.16 as builder

ADD . /project
WORKDIR /project/

RUN CGO_ENABLED=0 GOOS=linux go build -o jwt-mock ./build/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /project/jwt-mock .
COPY --from=builder /project/config.yaml ./config.yaml

EXPOSE 80

ENTRYPOINT ["./jwt-mock"]

CMD ["--config", "config.yaml"]
