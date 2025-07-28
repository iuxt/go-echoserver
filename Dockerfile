FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /echo-server main.go

FROM scratch AS runner
COPY --from=builder /echo-server /echo-server
EXPOSE 8080
ENTRYPOINT ["/echo-server"]