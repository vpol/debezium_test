FROM golang:1.19-alpine3.15 as builder

WORKDIR /app
ARG COMMAND_NAME
ARG COMMAND_PATH
ARG TAGS

COPY . /app/

RUN apk add --update --no-cache git

RUN ./.build/build.sh $COMMAND_PATH $COMMAND_NAME $TAGS

FROM google/cloud-sdk:alpine

ARG COMMAND_NAME
ARG COMMAND_PATH
ARG TAGS

WORKDIR /

RUN apk add --update --no-cache ca-certificates

COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/.build/target/$COMMAND_NAME /$COMMAND_NAME

ENV COMMAND_NAME $COMMAND_NAME
