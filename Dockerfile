FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o /star-tour .

FROM alpine:3.18
COPY --from=build /star-tour /star-tour
COPY templates /templates
COPY static /static
EXPOSE 8080
ENTRYPOINT ["/star-tour"]
