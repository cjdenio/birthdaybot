FROM golang:latest

WORKDIR /usr/src/app

COPY . .

RUN go get .

EXPOSE 3000

CMD [ "go", "run", "dev.go" ]