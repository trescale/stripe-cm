FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY stripe-cm /usr/local/bin/stripe-cm
ENTRYPOINT ["stripe-cm"]
