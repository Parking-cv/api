FROM golang:1.13-buster as builder
RUN mkdir /build
ADD ./ /
WORKDIR /
RUN go get -u github.com/go-chi/chi
RUN go get go.mongodb.org/mongo-driver/mongo 
RUN go build -o main .

FROM python:3.7-slim-buster
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
