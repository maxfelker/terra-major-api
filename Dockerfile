FROM golang:alpine as build 
WORKDIR /app
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY data/ data/
COPY go.mod go.mod
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/instance-api cmd/instance-api/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/bin ./bin
COPY --from=build /app/data ./data/
ENTRYPOINT ./bin/instance-api