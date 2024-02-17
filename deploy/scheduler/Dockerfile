FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata
RUN apk add make

COPY . /github.com/nikitads9/event-schedule/

WORKDIR /github.com/nikitads9/event-schedule/

RUN go mod download
RUN go build -o ./bin/event-schedule cmd/server/schedule.go
RUN chown -R root ./bin/scheduler

FROM alpine:latest

WORKDIR /bin/

COPY --from=builder /github.com/nikitads9/event-schedule/bin .
COPY --from=builder /github.com/nikitads9/event-schedule/config.yml .

CMD ["./scheduler", "-config", "config.yml"]