FROM golang:1.23-alpine@sha256:04ec5618ca64098b8325e064aa1de2d3efbbd022a3ac5554d49d5ece99d41ad5 AS builder
WORKDIR /src

ARG SCRUTZONE_VERSION

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -ldflags "-X \"main.scrutzoneVersion=${SCRUTZONE_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o /src/scrutzone


FROM alpine:3.21@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

WORKDIR /scrutzone
COPY --from=builder /src/scrutzone .

CMD ["./scrutzone"]
