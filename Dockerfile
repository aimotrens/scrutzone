FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS builder
WORKDIR /src

ARG SCRUTZONE_VERSION

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -ldflags "-X \"main.scrutzoneVersion=${SCRUTZONE_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o /src/scrutzone


FROM alpine:3.24@sha256:e7c4abb69531cb09e2a2bbb56fad3367ab694865c49df898c1c683185cc4376c

WORKDIR /scrutzone
COPY --from=builder /src/scrutzone .

CMD ["./scrutzone"]
