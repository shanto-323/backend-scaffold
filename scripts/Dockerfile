FROM golang:1.25.4-alpine3.22 AS builder
WORKDIR /project

COPY . .

RUN go build -o /cmd/bin/app -trimpath -mod=vendor \
    -gcflags "-trimpath $PWD" -asmflags "-trimpath $PWD" \
    -ldflags "-s -w -X main.version=$VERSION -X main.commitID=$COMMIT" \
    ./cmd

FROM alpine:3.22
WORKDIR /usr/bin
COPY --from=builder /cmd/bin/app .
CMD [ "app" ]
