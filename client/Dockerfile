FROM golang:1.18 as builder

WORKDIR /app

COPY ./go.mod ./

RUN go mod download

COPY ./cmd .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/client .

CMD ["./client"]
