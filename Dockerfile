FROM golang:1.23-alpine@sha256:ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06 AS builder
WORKDIR /src

ARG SCRUTZONE_VERSION

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -ldflags "-X \"main.scrutzoneVersion=${SCRUTZONE_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o /src/scrutzone


FROM alpine:3.20@sha256:e72ad0747b9dc266fca31fb004580d316b6ae5b0fdbbb65f17bbe371a5b24cff

WORKDIR /scrutzone
COPY --from=builder /src/scrutzone .

CMD ["./scrutzone"]
