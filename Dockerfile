#FROM golang:1.20 AS builder
#
#LABEL stage=gobuilder
#
#ENV GOOS linux
#
##RUN apk update --no-cashe && apk add --no-cashe tzdata
#
#WORKDIR /rpa_clone
#
#ADD go.mod ./
#
#ADD go.sum ./
#
#RUN go mod download
#
#COPY . .
#
#RUN go build -ldflags="-s -w" -o  app/build ./cmd/main.go
#
#FROM golang:1.20
#
#RUN #apk update --no-cashe && apk add --no-cashe ca-certificates
#
#COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York
#
#ENV TZ America/New_York
#
#WORKDIR /app
#
##COPY --from=builder /app/build /app/build
#
#CMD [". app/build"]
FROM golang:1.20 AS builder

ENV GOPATH=/

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

CMD ["./app" ]