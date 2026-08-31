FROM golang:alpine as builder
WORKDIR /app

COPY src/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /service ./cmd/app/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /service .
COPY --from=builder /app/web ./web

EXPOSE 80
CMD [ "./service" ]