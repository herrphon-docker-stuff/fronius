# TODO: move to minimal container l8r... and fixme... this is just typed from my head

FROM golang as builder
WORKDIR /app
ADD . .
RUN go build
RUN ls -la

# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
#RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

FROM ubuntu
# FROM scratch
# FROM alpine
COPY --from=builder /app/fronius /fronius

# the tls certificates:
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV INFLUX_SECRET_TOKEN=foobar

CMD ["/fronius", "serve"]


#FROM golang:1.20

#RUN go get -u github.com/tgulacsi/fronius
## RUN go install github.com/tgulacsi/fronius@latest

#ENV INFLUX_USER=influxusername
#ENV INFLUX_PASSW=influxuserpassword

# CMD ["fronius", "serve"]


