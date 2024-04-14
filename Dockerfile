FROM golang:1.22-alpine as builder

WORKDIR /go/src/banner

COPY go.mod go.sum /go/src/banner/
RUN go mod download && go mod verify

COPY . .
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
RUN golangci-lint run -v
RUN go build -v -o /usr/local/bin/app ./cmd/server

FROM alpine:3
COPY --from=builder /usr/local/bin/app /bin/main
ENTRYPOINT ["/bin/main"]