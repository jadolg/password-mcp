FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w" \
      -trimpath \
      -o password-mcp .

FROM scratch
COPY --from=builder /build/password-mcp /password-mcp
EXPOSE 8080

ENTRYPOINT ["/password-mcp", "-http", "-addr", ":8080"]
