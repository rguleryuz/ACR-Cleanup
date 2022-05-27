FROM golang:1.18.1-bullseye
RUN mkdir app
WORKDIR /app
COPY . .
RUN go mod init app
RUN go get -d ./...
RUN go build -o app
ENTRYPOINT [ "./app" ]