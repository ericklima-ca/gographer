# syntax=docker/dockerfile:1

FROM  golang:1.19-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -o /gographer

FROM gcr.io/distroless/base-debian11
LABEL maintainer="Erick Amorim <ericklima.ca@yahoo.com>"
COPY --from=builder /gographer /gographer
EXPOSE 8000
ENV GIN_MOD=release
ENTRYPOINT ["/gographer"]
