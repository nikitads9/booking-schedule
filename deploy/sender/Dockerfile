FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata
RUN apk add make

COPY . /github.com/nikitads9/booking-schedule/

WORKDIR /github.com/nikitads9/booking-schedule/

RUN make build-sender
RUN chown -R root ./bin/sender

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/nikitads9/booking-schedule/bin .
COPY --from=builder /github.com/nikitads9/booking-schedule/configs/sender_config.yml .

CMD ["./sender", "-config", "sender_config.yml"]