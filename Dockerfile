# Builder container 
FROM golang:1.13-buster as builder
RUN mkdir /build
ADD ./src /go/src/parking-cv/src
WORKDIR /go/src/parking-cv

# Install dependencies
RUN go get -u github.com/go-chi/chi
RUN go get go.mongodb.org/mongo-driver/mongo 

# Compile 
RUN go build -o ./build/main ./src


################################################################

# Runtime container
FROM python:3.7-slim-buster

WORKDIR /

RUN apt-get update
RUN yes | apt-get install build-essential
RUN apt-get -y install cmake

WORKDIR /

# Install requirements
COPY ./requirements.txt /requirements.txt
RUN pip install -r requirements.txt

# Copy files from builder container
COPY --from=builder /go/src/parking-cv/build/main /app/main
COPY ./assets /app/assets
WORKDIR /app

# Change user
# RUN adduser -S -D -H -h /app appuser
# USER appuser

# Start app 
EXPOSE 4321
CMD ["./main"]
