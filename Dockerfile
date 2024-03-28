FROM golang:1.22-alpine as build

WORKDIR /go/src/github.com/alexellis/firecracker-init-lab/init

COPY init .

RUN go build --tags netgo --ldflags '-s -w -extldflags "-lm -lstdc++ -static"' -o init main.go

FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine3.19

RUN apk add --no-cache curl ca-certificates htop && \
    echo "nameserver 1.1.1.1" > /etc/resolv.conf

COPY --from=build /go/src/github.com/alexellis/firecracker-init-lab/init/init /init