FROM golang:1.13-buster as builder
RUN mkdir /build
ADD ./src /build/
WORKDIR /build
RUN go get -u github.com/go-chi/chi
RUN go build -o main .

FROM python:3.7-slim-buster
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
