# syntax=docker/dockerfile:1

FROM golang:1.22.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
# COPY config /config

RUN go mod tidy

COPY . .

#Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /server
# RUN touch .env
# EXPOSE 8080

# CMD ["/server"]

FROM alpine

WORKDIR /
COPY --from=build-stage /server /server
RUN touch .env
EXPOSE 8080

CMD ["/server"]



