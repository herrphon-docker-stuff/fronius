# TODO: move to minimal container l8r... and fixme... this is just typed from my head

#FROM golang as builder
#ADD . .
#RUN go build

#FROM scratch
#COPY --from=builder /app/app /app
#CMD ["/app"]


FROM golang:1.16

RUN go get -u github.com/tgulacsi/fronius
# RUN go install github.com/tgulacsi/fronius@latest

ENV INFLUX_USER=influxusername
ENV INFLUX_PASSW=influxuserpassword 

CMD ["fronius", "serve"]


