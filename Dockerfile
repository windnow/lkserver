FROM golang:1.23-alpine as builder

WORKDIR /build

RUN apk add --no-cache git build-base

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN ls -al

RUN go build -o main ./cmd/lkserver

FROM alpine:latest
WORKDIR /srv/
RUN apk update && apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /build/main .
RUN ls -al
COPY --from=builder /build/data ./data
RUN mkdir res
RUN touch res/config.toml
RUN pwd
RUN ls ./data -al
EXPOSE 8833
ENTRYPOINT ["./main"]