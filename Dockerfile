FROM golang:1.8 AS build
COPY . /go/src/github.com/watercraft/go-md2man
WORKDIR /go/src/github.com/watercraft/go-md2man
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /go/src/github.com/watercraft/go-md2man/go-md2man /go-md2man
ENTRYPOINT ["/go-md2man"]
