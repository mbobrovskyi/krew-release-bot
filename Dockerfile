FROM golang:1.17.3-alpine3.14 as builder

WORKDIR /go/src/github.com/rajatjindal/krew-release-bot
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go test -mod vendor ./... -cover
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor --ldflags "-s -w" -o krew-release-bot cmd/action/*

FROM alpine:3.14.3

RUN mkdir -p /home/app

# Add non root user
RUN addgroup -S app && adduser app -S -G app
RUN chown app /home/app

WORKDIR /home/app

USER app

COPY --from=builder /go/src/github.com/rajatjindal/krew-release-bot/krew-release-bot /usr/local/bin/

CMD ["krew-release-bot", "action"]
