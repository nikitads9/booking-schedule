FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata
RUN apk add make

COPY . /github.com/nikitads9/booking-schedule/

WORKDIR /github.com/nikitads9/booking-schedule/

RUN make build-bookings
RUN chown -R root ./bin/bookings

FROM alpine:latest


RUN apk add nginx
ADD ./deploy/nginx/nginx.conf /etc/nginx/nginx.conf


WORKDIR /root/

COPY --from=builder /github.com/nikitads9/booking-schedule/bin .
COPY --from=builder /github.com/nikitads9/booking-schedule/configs/booking_config.yml .

CMD [ "./bookings", "-config", "booking_config.yml", "-keyfile", "/etc/ssl/certs/key.pem", "-certfile", "/etc/ssl/certs/cert.pem" ]