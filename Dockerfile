FROM golang:1.23-alpine AS builder
ARG VERSION=dev
ARG COMMIT=none
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o /stripe-cm ./cmd/stripe-cm/

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /stripe-cm /usr/local/bin/stripe-cm
ENTRYPOINT ["stripe-cm"]
