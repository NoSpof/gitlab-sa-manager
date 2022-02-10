FROM golang:1.17.6-alpine as build-env
COPY src /go/src
WORKDIR /go/src
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static
COPY --from=build-env /go/bin/app /
CMD ["/app"]