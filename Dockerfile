FROM golang:1.23-alpine AS build

WORKDIR /src/
COPY go.* .
RUN go mod download

COPY . /src/
RUN CGO_ENABLED=0 go build -o out/bin/swim-spot-checker ./

FROM alpine:latest

RUN mkdir -p app
WORKDIR /app
COPY --from=build /src/out/bin/swim-spot-checker /app/

RUN apk update && apk add --no-cache dcron tini

RUN echo "* * * * * /app/swim-spot-checker >> /var/log/cron.log 2>&1" > /etc/crontabs/root

VOLUME /var/log

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["crond", "-f"]

