FROM golang:1.22.5-alpine

WORKDIR /usr/src/app

ENV TZ="Asia/Bangkok"

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

RUN go build -o ./build/server
EXPOSE 3000
CMD [ "./build/server"]
