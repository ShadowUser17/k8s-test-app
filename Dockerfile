FROM golang:1.19-alpine AS build

WORKDIR /root

COPY ./cmd/main.go ./
COPY ./go.mod ./

RUN go mod tidy
RUN go build -ldflags="-s -w" -o ./app ./main.go


FROM alpine:3.17.0

WORKDIR /root
COPY --from=build /root/app ./

EXPOSE 8080/tcp

ENTRYPOINT ["/root/app"]
CMD ["-listen", "0.0.0.0:8080"]
