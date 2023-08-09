FROM golang:alpine as build 
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY main.go main.go
COPY pkg/ pkg/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
RUN mkdir keys/
RUN go run pkg/auth/pem/main.go /app/keys terra-major-client
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/terra-major-api main.go

FROM alpine:latest as release  
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/bin ./bin
COPY --from=build /app/keys ./keys
ENTRYPOINT ./bin/terra-major-api