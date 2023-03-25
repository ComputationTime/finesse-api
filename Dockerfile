FROM golang:1.20

WORKDIR /usr/src/api

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build 

CMD ["./finesse-api"]