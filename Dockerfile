FROM golang:1.22.5-alpine

WORKDIR /usr/src/app

RUN apk add tzdata
ENV TZ Asia/Bangkok

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

RUN go build -o ./output/zuck_my_cloth_backend
EXPOSE 3000
CMD [ "./output/zuck_my_cloth_backend"]
