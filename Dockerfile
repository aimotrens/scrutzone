FROM golang:1.25-alpine@sha256:72a77777b4d52d80820a9b5828c06cca1cd3c24721b6ca25b0269b5f6fb00daa AS builder
WORKDIR /src

ARG SCRUTZONE_VERSION

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -ldflags "-X \"main.scrutzoneVersion=${SCRUTZONE_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o /src/scrutzone


FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

WORKDIR /scrutzone
COPY --from=builder /src/scrutzone .

CMD ["./scrutzone"]
