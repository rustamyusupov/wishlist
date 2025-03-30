FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o wishes ./cmd/server/main.go

FROM alpine:latest

ARG AUTH_EMAIL
ARG AUTH_PASSWORD_HASH
ARG DB_URL
ARG SESSION_KEY
ARG SESSION_NAME

ENV AUTH_EMAIL=$AUTH_EMAIL
ENV AUTH_PASSWORD_HASH=$AUTH_PASSWORD_HASH
ENV DB_URL=$DB_URL
ENV SESSION_KEY=$SESSION_KEY
ENV SESSION_NAME=$SESSION_NAME

RUN apk add --no-cache libc6-compat ca-certificates

WORKDIR /app

COPY --from=builder /app/wishes /app/wishes
COPY --from=builder /app/web /app/web

EXPOSE 8080
RUN mkdir -p /app/data

CMD ["/app/wishes"]