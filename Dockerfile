FROM golang AS build
WORKDIR /go/src/app/
COPY . .
RUN go get -d
RUN go build -i -o transmission-exporter

FROM alpine
RUN apk update && apk add --no-cache \
        libc6-compat ca-certificates
COPY --from=build /go/src/app/transmission-exporter /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/transmission-exporter"]
