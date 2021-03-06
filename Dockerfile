FROM golang:1.17-alpine

RUN mkdir -p auth-etf2l/bin/
WORKDIR auth-etf2l/

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY pkg/ .

RUN go build -o ./bin/app

FROM scratch

COPY src/ . 
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/auth-etf2l/bin/app .

ENTRYPOINT ["./app"]
