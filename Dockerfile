FROM golang:1.22-alpine AS build 

WORKDIR /epl-fantasy

COPY src/config/  src/config/
COPY src/db/  src/db/
COPY src/handlers/  src/handlers/
COPY src/service/  src/service/
COPY main.go  main.go
COPY URL.env  URL.env
RUN go mod init epl-fantasy
RUN go get .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -o main .

FROM alpine:latest

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /usr/bin

COPY --from=build /epl-fantasy/main .
COPY --from=build /epl-fantasy/URL.env .

EXPOSE 8080

CMD ["./main"]



