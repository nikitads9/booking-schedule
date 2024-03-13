FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata
RUN apk add make

COPY . /github.com/nikitads9/booking-schedule/

WORKDIR /github.com/nikitads9/booking-schedule/

RUN make build-auth
RUN chown -R root ./bin/auth

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/nikitads9/booking-schedule/bin .
COPY --from=builder /github.com/nikitads9/booking-schedule/configs/auth_config.yml .

CMD ["./bookings", "-config", "auth_config.yml"]