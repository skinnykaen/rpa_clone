FROM golang:1.18 AS builder

LABEL stage=gobuilder

ENV GOOS linux

RUN apk update --no-cashe && apk add --no-cashe tzdata

WORKDIR /rpa_clone

ADD go.mod ./

ADD go.sum ./

RUN go mod download

COPY build .

RUN go build -ldflags="-s -w" -o /app/build . cmd/main.go

FROM golang:1.18

RUN apk update --no-cashe && apk add --no-cashe ca-certificates

COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York

ENV TZ America/New_York

WORKDIR /app

COPY --from=builder /app/build /app/build

CMD [". /build"]