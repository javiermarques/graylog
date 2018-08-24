FROM golang:alpine
# Install SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo main.go

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/app/main /
ENTRYPOINT ["/main"]
