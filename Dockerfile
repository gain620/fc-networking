FROM golang:1.22-alpine as build

WORKDIR /go/src/github.com/alexellis/firecracker-init-lab/init

COPY init .

RUN go build --tags netgo --ldflags '-s -w -extldflags "-lm -lstdc++ -static"' -o init main.go

FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine3.19

RUN apk add --no-cache curl ca-certificates htop

COPY --from=build /go/src/github.com/alexellis/firecracker-init-lab/init/init /init/init-go
COPY dotnet /init/dotnet-hello