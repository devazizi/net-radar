FROM golang:1.22 as builder

WORKDIR /app

ENV GOPROXY=direct
ENV GOSUMDB=off

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/net-radar ./cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/net-radar /app/net-radar

ENV GIN_MODE release
RUN chmod +x /app/proxy

EXPOSE 3000

ENTRYPOINT ["/app/proxy"]