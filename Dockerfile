FROM golang:1.23-bullseye AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/server /app/cmd/main.go

FROM ubuntu:22.04
RUN mkdir -p /app/cmd
WORKDIR /app
COPY --from=builder /app/server .
ENTRYPOINT ["/app/server"]
