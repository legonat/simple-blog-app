FROM golang:1.17-buster

WORKDIR /app

RUN go version
ENV GOPATH=/

COPY ./ /app

# build go app
RUN go mod download
RUN go build -o blog-app ./cmd/app/main.go

CMD ["./blog-app -f initVmDb"]
CMD ["./blog-app -f startServer"]