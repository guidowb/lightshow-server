FROM golang:alpine as build

WORKDIR /app
ADD . /app
RUN go build -o lightshow-server main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=build /app/lightshow-server /app/
WORKDIR /app
CMD ["./lightshow-server"]