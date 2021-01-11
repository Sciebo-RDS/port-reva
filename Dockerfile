FROM golang:1.13

WORKDIR /go/src/github/Sciebo-RDS/port-reva
COPY . .
RUN go build -o port-reva ./cmd && cp /go/src/github/Sciebo-RDS/port-reva/port-reva /go/bin/port-reva
ENTRYPOINT ["/go/bin/port-reva"]
