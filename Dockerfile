FROM golang:alpine as build

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY . .

RUN go build -o portscan ./cmd/main.go


FROM alpine:latest

RUN addgroup -g 1000 portscan \
    && adduser -u 1000 -G portscan -s /bin/sh -D portscan

RUN apk add --no-cache ca-certificates nmap

WORKDIR /app
RUN chown portscan:portscan /app

COPY --from=build /app/portscan /usr/local/bin/portscan

EXPOSE 8080

CMD ["portscan"]
