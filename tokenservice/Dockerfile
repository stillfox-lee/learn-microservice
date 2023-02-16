FROM golang:1.19 AS build

WORKDIR /app/src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o token .

FROM golang:1.19 AS dev
COPY --from=build /app/src/token /app/bin/token
CMD ["/app/bin/token", "--redis-host", "redis-db", "--redis-port", "6379", "--redis-password", "test124", "--redis-db", "0", "--server-port", "80", "--service-name", "token_service", "--consul-port", "8500", "--consul-host", "consul"]