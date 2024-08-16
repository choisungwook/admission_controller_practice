# syntax=docker/dockerfile:1.2

# build go project
FROM golang:1.21.6 AS build-stage

WORKDIR /app

COPY src/* /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /server /server

EXPOSE 443

ENTRYPOINT ["/server"]
