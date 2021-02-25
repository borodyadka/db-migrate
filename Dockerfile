FROM golang:1.16-alpine AS builder

ENV CGO_ENABLED 0

WORKDIR /app
COPY . /app

RUN go mod vendor && go build -mod vendor -a -ldflags "-w -s" -installsuffix cgo -o ./bin/migrate ./cmd/migrate

FROM alpine:3.12

EXPOSE 8082

WORKDIR /migrations
ENV LOG_LEVEL INFO
ENV SOURCE_DIR /migrations

COPY --from=builder /app/bin/migrate /bin/migrate

ENTRYPOINT ["/bin/migrate"]
CMD ["up"]
