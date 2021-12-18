FROM golang:1.17-stretch as builder

WORKDIR /go/src/vtuber-cho

COPY go.mod go.sum ./
RUN go mod download


COPY ./ /go/src/vtuber-cho

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/vtuber-cho/main

# FROM scratch as prod

# WORKDIR /go/src/vtuber-cho
# COPY --from=builder /go/src/vtuber-cho/main /go/src/

EXPOSE 8000
CMD ["./main"]
