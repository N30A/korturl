FROM golang:1.27.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/korturl ./cmd/korturl

FROM gcr.io/distroless/static-debian13

WORKDIR /app

COPY --from=builder /app/korturl /app/korturl

USER 1000:1000

EXPOSE 8080

CMD ["/app/korturl"]
