FROM golang:1.25-alpine@sha256:3587db7cc96576822c606d119729370dbf581931c5f43ac6d3fa03ab4ed85a10 AS builder
WORKDIR /src

ARG SCRUTZONE_VERSION

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -ldflags "-X \"main.scrutzoneVersion=${SCRUTZONE_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o /src/scrutzone


FROM alpine:3.23@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375

WORKDIR /scrutzone
COPY --from=builder /src/scrutzone .

CMD ["./scrutzone"]
