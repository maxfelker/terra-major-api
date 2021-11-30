FROM golang:alpine as dev 
WORKDIR /app
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY data/ data/
COPY go.mod go.mod
RUN go mod download
EXPOSE 8000
CMD go run cmd/instance-api/main.go

FROM golang:alpine as release 
WORKDIR /app
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY data/ data/
COPY go.mod go.mod
RUN go mod download
RUN PORT=80 go build -o bin/instance-api cmd/instance-api/main.go
EXPOSE 80
#CMD ls -al bin/
#CMD ./bin/instance-api
CMD PORT=80 go run cmd/instance-api/main.go