FROM golang as builder

ADD . .
RUN go build

FROM scratch

COPY --from=builder /app/app /app

CMD ["/app"]

// TODO: fixme... this is just typed from my head

FROM golang

RUN go get -u github.com/tgulacsi/fronius


