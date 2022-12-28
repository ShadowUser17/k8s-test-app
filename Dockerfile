FROM golang:1.19-bullseye AS build

WORKDIR /root

COPY ./cmd ./
COPY ./go.mod ./

RUN go mod tidy
RUN go build -ldflags="-s -w" -o ./app ./cmd/main.go


FROM alpine:3.17.0

WORKDIR /root
COPY --from=build ./app ./

EXPOSE 8080/tcp

ENTRYPOINT ["./app"]
CMD ["-listen", "0.0.0.0:8080"]
