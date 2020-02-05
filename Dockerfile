# Get dependencies and compile app in temporary container (359 MB)
FROM golang:alpine as builder
RUN mkdir /build
ADD ./src /build/
WORKDIR /build
RUN apk update && \
    apk upgrade && \
    apk add git
RUN go get -u github.com/go-chi/chi
RUN go build -o main .

# Copy binary into prod container (16 MB)
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
