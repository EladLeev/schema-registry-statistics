FROM golang:alpine AS builder
ADD . /schema-registry-statistics
WORKDIR /schema-registry-statistics
RUN apk add git
RUN go build -o sr-stat . && cp sr-stat /sr-stat

FROM golang:alpine
COPY --from=builder /sr-stat /usr/bin/sr-stat
LABEL maintainer="eladleev@gmail.com"

CMD ["/usr/bin/sr-stat"]
